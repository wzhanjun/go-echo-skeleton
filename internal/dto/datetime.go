package dto

import (
	"reflect"
	"time"

	"github.com/wzhanjun/go-echo-skeleton/internal/enum"
)

type Date struct {
	time.Time
}

func (s *Date) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}
	now, err := time.ParseInLocation(`"`+enum.DateTpl+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*s = Date{now}
	return
}

func (s Date) MarshalJSON() ([]byte, error) {
	if s.Time.IsZero() {
		return []byte(`""`), nil
	}
	b := make([]byte, 0, len(enum.DateTpl)+2)
	b = append(b, '"')
	b = s.Time.AppendFormat(b, enum.DateTpl)
	b = append(b, '"')
	return b, nil
}

func (s *Date) UnmarshalParam(data string) (err error) {
	if len(data) == 0 {
		return nil
	}
	now, err := time.ParseInLocation(enum.DateTpl, data, time.Local)
	if err != nil {
		return err
	}
	*s = Date{now}
	return
}

type DateTime struct {
	time.Time
}

func (s *DateTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 || reflect.DeepEqual(data, []byte("\"\"")) {
		return nil
	}
	now, err := time.ParseInLocation(`"`+enum.DateTimeTpl+`"`, string(data), time.Local)
	if err != nil {
		return err
	}
	*s = DateTime{now}
	return
}

func (s DateTime) MarshalJSON() ([]byte, error) {
	if s.Time.IsZero() {
		return []byte(`""`), nil
	}
	b := make([]byte, 0, len(enum.DateTimeTpl)+2)
	b = append(b, '"')
	b = s.Time.AppendFormat(b, enum.DateTimeTpl)
	b = append(b, '"')
	return b, nil
}

func (s *DateTime) UnmarshalParam(data string) (err error) {
	if len(data) == 0 {
		return nil
	}
	now, err := time.ParseInLocation(enum.DateTimeTpl, data, time.Local)
	if err != nil {
		return err
	}
	*s = DateTime{now}
	return
}
