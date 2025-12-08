package scheduler

import (
	"context"
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

	CodeIdJob(scheduler)

	scheduler.Start()
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
