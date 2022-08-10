package pkg

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Timer timer
)

const (
	FormatDatetime = "2006-01-02 15:04:05"
	FormatDate     = "2006-01-02"
	FormatTime     = "15:04:05"
)

type timer struct{}

func (timer) NowString() string {
	return time.Now().Format(FormatDatetime)
}

func (timer) Time2String(t time.Time) string {
	return t.Format(FormatDatetime)
}

func (timer) String2Time(s string) time.Time {
	t, err := time.ParseInLocation(FormatDatetime, s, time.Local)
	if err != nil {
		logrus.WithField("String2Time", s).Warn(err)
	}
	return t
}
