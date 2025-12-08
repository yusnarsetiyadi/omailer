package scheduler

import (
	"context"
	"omailer/pkg/general"
	"omailer/pkg/whatsapp"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/sirupsen/logrus"
)

func CodeIdJob(scheduler gocron.Scheduler) {
	// AttendanceIn
	_, err := scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0))),
		gocron.NewTask(AttendanceInWrapper, true),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(11, 0, 0))),
		gocron.NewTask(AttendanceInWrapper, false),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	// AttendanceOut
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(17, 0, 0))),
		gocron.NewTask(AttendanceOutWrapper, true),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(20, 0, 0))),
		gocron.NewTask(AttendanceOutWrapper, false),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}
}

func isWorkday(t time.Time) bool {
	w := t.Weekday()
	return w != time.Saturday && w != time.Sunday
}

func AttendanceInWrapper(isFirst bool) {
	now := general.NowLocal()
	if !isWorkday(*now) {
		logrus.Println("Attendance IN skipped (weekend)")
		return
	}
	AttendanceIn(isFirst)
}

func AttendanceOutWrapper(isFirst bool) {
	now := general.NowLocal()
	if !isWorkday(*now) {
		logrus.Println("Attendance OUT skipped (weekend)")
		return
	}
	AttendanceOut(isFirst)
}

func AttendanceIn(isFirst bool) {
	// send message to group
	logrus.Info("send message to group (attendance in)")

	var err error
	if isFirst {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN MASUK YA GUYS, UDAH PAGI!"),
		)
	} else {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN MASUK YA GUYS, UDAH SIANG!"),
		)
	}
	if err != nil {
		logrus.Error("Gagal kirim pesan ke grup:", err)
	}
}

func AttendanceOut(isFirst bool) {
	// send message to group
	logrus.Info("send message to group (attendance out)")

	var err error
	if isFirst {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN KELUAR YA GUYS, UDAH SORE!"),
		)
	} else {
		err = whatsapp.SendText(
			context.Background(),
			"grup masjid lt 20",
			AutomatedMessage("JANGAN LUPA ABSEN KELUAR YA GUYS, UDAH MALEM!"),
		)
	}
	if err != nil {
		logrus.Error("Gagal kirim pesan ke grup:", err)
	}
}
