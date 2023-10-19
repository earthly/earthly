package flagutil

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Duration implements cli.GenericFlag methods to support time.Duration with days, e.g. 1d
type Duration time.Duration

func (d *Duration) String() string {
	return (time.Duration(*d)).String()
}

func (d *Duration) Set(value string) error {
	if value == "" {
		return nil
	}
	daysToHours := false
	if strings.HasSuffix(value, "d") {
		value = fmt.Sprintf("%s%s", strings.TrimSuffix(value, "d"), "h")
		daysToHours = true
	}
	dur, err := time.ParseDuration(value)
	if err != nil {
		return errors.New("parse error")
	}

	if daysToHours {
		dur *= 24
	}
	*d = Duration(dur)
	return nil
}
