package cloud

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMaintenanceWindow(t *testing.T) {
	t.Run("proper timestamp", func(t *testing.T) {
		l, _ := time.LoadLocation("EST")
		w, err := ParseMaintenanceWindow("02:00", l)
		assert.Nil(t, err)
		assert.Equal(t, "07:00", w)
	})
	t.Run("invalid formats", func(t *testing.T) {
		l, _ := time.LoadLocation("EST")
		_, err := ParseMaintenanceWindow("1234:00", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("02:1234", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("02:00:12", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("1:00pm", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("13:14 EST", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("1224", l)
		assert.NotNil(t, err)
		_, err = ParseMaintenanceWindow("oops", l)
		assert.NotNil(t, err)
	})
}
