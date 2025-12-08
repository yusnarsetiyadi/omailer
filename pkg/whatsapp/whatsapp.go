package whatsapp

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mdp/qrterminal/v3"
	"github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	_ "modernc.org/sqlite"
)

var (
	clientMu sync.RWMutex
	client   *whatsmeow.Client
	isReady  atomic.Bool
)

func Init() error {
	if err := os.MkdirAll("data", 0o700); err != nil {
		return err
	}

	dbFile := "data/session.db"

	dsn := fmt.Sprintf(
		"file:%s?_pragma=foreign_keys=ON&_pragma=journal_mode=WAL&_pragma=busy_timeout=5000",
		dbFile,
	)
	container, err := sqlstore.New(context.Background(), "sqlite", dsn, nil)
	if err != nil {
		return fmt.Errorf("failed create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return err
	}

	c := whatsmeow.NewClient(deviceStore, nil)

	clientMu.Lock()
	client = c
	clientMu.Unlock()

	c.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {

		case *events.Connected:
			isReady.Store(true)
			logrus.Info("WhatsApp client connected & ready")

		case *events.AppStateSyncComplete:
			isReady.Store(true)
			logrus.Info("AppState sync complete")

		case *events.PairSuccess:
			logrus.Infof("PairSuccess! Logged in as: %v", v.ID)

		case *events.ConnectFailure:
			isReady.Store(false)
			logrus.Warnf("Connection failed: %v", v.Reason)

		default:
			logrus.Debugf("Unhandled event: %T", evt)
		}
	})

	if c.Store.ID == nil {
		logrus.Warn("WA not paired yet, showing QR code...")

		qrChan, _ := c.GetQRChannel(context.Background())

		err = c.Connect()
		if err != nil {
			return fmt.Errorf("failed connect: %w", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				logrus.Info("Scan QR above using WhatsApp (Linked Devices)")
			} else {
				logrus.Infof("QR Event: %s", evt.Event)
			}
		}

		logrus.Info("Device paired!")
	} else {
		logrus.Info("Restoring existing WhatsApp session...")
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() {
	clientMu.RLock()
	c := client
	clientMu.RUnlock()

	if c == nil {
		return
	}

	logrus.Info("Logging out and deleting WhatsApp session...")

	err := c.Logout(context.Background())
	if err != nil {
		logrus.Warnf("Failed to logout device: %v", err)
	}

	c.Disconnect()

	err = os.Remove("data/session.db")
	if err != nil {
		logrus.Warnf("Failed to remove session.db: %v", err)
	} else {
		logrus.Info("Local WhatsApp session database removed.")
	}
}

func ensureConnected(c *whatsmeow.Client) error {
	if c == nil {
		return fmt.Errorf("client nil")
	}
	if !c.IsConnected() {
		if err := c.Connect(); err != nil {
			return fmt.Errorf("failed reconnect: %w", err)
		}
	}
	if !isReady.Load() {
		return fmt.Errorf("WhatsApp still syncing, not ready")
	}

	return nil
}

func WaitUntilReady(ctx context.Context) error {
	timeout := time.After(20 * time.Second)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timeout:
			return fmt.Errorf("timeout waiting WhatsApp ready")

		case <-ticker.C:
			clientMu.RLock()
			c := client
			clientMu.RUnlock()

			if isReady.Load() && c != nil && c.IsConnected() {
				return nil
			}
		}
	}
}
