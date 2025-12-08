package main

import (
	"context"
	"net/http"
	httpomailer "omailer/internal/http"
	middlewareEcho "omailer/internal/middleware"
	"omailer/internal/scheduler"
	"omailer/pkg/constant"
	"omailer/pkg/log"
	"omailer/pkg/whatsapp"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// @title omailer
// @version 1.0.0
// @description This is a doc for omailer

func main() {

	log.Init()

	e := echo.New()

	middlewareEcho.Init(e)

	httpomailer.Init(e)

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// whatsapp & scheduler
	if err := whatsapp.Init(); err != nil {
		logrus.Fatalf("Failed to init WhatsApp client: %v", err)
	}
	scheduler.InitScheduler()

	go func() {
		addr := ":" + strconv.Itoa(constant.PORT)
		err := e.Start(addr)
		if err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatal(err)
			}
		}
	}()

	<-ch

	whatsapp.Close()

	logrus.Println("Shutting down server...")
	cancel()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	e.Shutdown(ctx2)
	logrus.Println("Server gracefully stopped")
}
