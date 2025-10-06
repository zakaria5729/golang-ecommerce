package timeutil

import (
	"time"

	"gorm.io/gorm"
)

func NowUTC() time.Time {
	return time.Now().UTC()
}

func GormNowUTC() *gorm.DeletedAt {
	return &gorm.DeletedAt{
		Time:  time.Now().UTC(),
		Valid: true,
	}
}

func ParseTimeUTC(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func ParseTimeInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func FormatTimeUTC(t time.Time, layout string) string {
	return t.UTC().Format(layout)
}

func IsTimeExpired(t time.Time) bool {
	return t.Before(NowUTC())
}

func IsTimeValid(t time.Time) bool {
	return t.After(NowUTC())
}

func AddDurationUTC(d time.Duration) time.Time {
	return NowUTC().Add(d)
}

func AddSecondsUTC(seconds int64) time.Time {
	return NowUTC().Add(time.Duration(seconds) * time.Second)
}

func AddHoursUTC(hours int) time.Time {
	return NowUTC().Add(time.Duration(hours) * time.Hour)
}

func AddDaysUTC(days int) time.Time {
	return NowUTC().Add(time.Duration(days) * 24 * time.Hour)
}

func ConvertUTCToLocal(t time.Time) time.Time {
	return t.In(time.Local)
}

func ConvertLocalToUTC(t time.Time) time.Time {
	return t.In(time.UTC)
}
