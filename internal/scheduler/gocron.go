package scheduler

import (
	"context"
	"omailer/pkg/general"
	"omailer/pkg/whatsapp"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

func InitScheduler() {
	jakartaTime, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Errorln(err.Error())
	}

	scheduler, err := gocron.NewScheduler(gocron.WithLocation(jakartaTime))
	if err != nil {
		logrus.Errorln(err.Error())
	}

	go WaitUntilWAReadyThenRun("TestGoCron", TestGoCron)

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0))),
		gocron.NewTask(AttendanceInWrapper),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(17, 0, 0))),
		gocron.NewTask(AttendanceOutWrapper, false),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(20, 0, 0))),
		gocron.NewTask(AttendanceOutWrapper, true),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	scheduler.Start()
}

func isWorkday(t time.Time) bool {
	w := t.Weekday()
	return w != time.Saturday && w != time.Sunday
}

func AttendanceInWrapper() {
	now := general.NowLocal()
	if !isWorkday(*now) {
		logrus.Println("Attendance IN skipped (weekend)")
		return
	}
	AttendanceIn()
}

func AttendanceOutWrapper(isNight bool) {
	now := general.NowLocal()
	if !isWorkday(*now) {
		logrus.Println("Attendance OUT skipped (weekend)")
		return
	}
	AttendanceOut(isNight)
}

func WaitUntilWAReadyThenRun(name string, fn func()) {
	go func() {
		logrus.Infof("Waiting WhatsApp ready before running %s...", name)

		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		if err := whatsapp.WaitUntilReady(ctx); err != nil {
			logrus.Errorf("Failed waiting WA ready for %s: %v", name, err)
			return
		}

		logrus.Infof("WhatsApp ready → running %s", name)
		fn()
	}()
}

func AutomatedMessage(mainText string) string {
	footer := "_Automated message.\nCoffee keeps this running ☕_"
	return mainText + "\n\n\n" + footer
}

func TestGoCron() {
	// test send message to me
	logrus.Info("test send message")

	err := whatsapp.SendText(
		context.Background(),
		"6281398447822",
		AutomatedMessage("Cron Test: Hello from WhatsMeow!"),
	)

	if err != nil {
		logrus.Error("Gagal kirim pesan:", err)
	}
}

func AttendanceIn() {
	// send message to group
	logrus.Info("send message to group (attendance in)")

	err := whatsapp.SendText(
		context.Background(),
		"grup masjid lt 20",
		AutomatedMessage("Haloo mas/mba semua, JANGAN LUPA ABSEN MASUK YA PAGI INI. Terima kasih!"),
	)

	if err != nil {
		logrus.Error("Gagal kirim pesan ke grup:", err)
	}
}

func AttendanceOut(isNight bool) {
	// send message to group
	logrus.Info("send message to group (attendance out)")

	var err error
	if isNight {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN KELUAR YA GES, UDAH MALEM. Terima kasih!"),
		)
	} else {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN KELUAR YA GES, UDAH SORE. Terima kasih!"),
		)
	}

	if err != nil {
		logrus.Error("Gagal kirim pesan ke grup:", err)
	}
}
