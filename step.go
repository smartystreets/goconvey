package goconvey

import (
	"testing"
)

func catch(t *testing.T) func() {
	return func() {
		if r := recover(); r != nil {
			t.Errorf("%v", r)
		}
	}
}
