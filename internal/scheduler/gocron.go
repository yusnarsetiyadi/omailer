package scheduler

import (
	"omailer/pkg/general"
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

	TestGoCron()

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(8, 0, 0))),
		gocron.NewTask(AttendanceInWrapper),
	)
	if err != nil {
		logrus.Errorln(err.Error())
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(17, 0, 0))),
		gocron.NewTask(AttendanceOutWrapper),
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

func AttendanceOutWrapper() {
	now := general.NowLocal()
	if !isWorkday(*now) {
		logrus.Println("Attendance OUT skipped (weekend)")
		return
	}
	AttendanceOut()
}

func TestGoCron() {
	// test send message to me
	logrus.Info("test send message to me")
}

func AttendanceIn() {
	// send message to group
	logrus.Info("send message to group (attendance in)")
}

func AttendanceOut() {
	// send message to group
	logrus.Info("send message to group (attendance out)")
}
