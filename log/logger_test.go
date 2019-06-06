package log

import (
	"testing"
)

func TestKit(t *testing.T) {
	l := WithField(Fields{
		"k_string": "v1",
		"k_int":    100,
	})

	l.Info("hello world!")
	l.Info("this is msg.",
		Fields{
			"k_string_inner": "v1",
			"k_int_inner":    100,
		},
	)

	l.Debug("this is another msg.",
		Fields{
			"k_string_other": "v1",
			"k_int_other":    100,
		},
	)

	t.Logf("kit ok!")
}
