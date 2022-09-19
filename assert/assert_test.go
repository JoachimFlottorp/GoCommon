package assert

import (
	"errors"
	"testing"
)

func TestAssert(t *testing.T) {
	t.Run("ErrAssert runs", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		
		err := errors.New("T")

		Error(err, "T")
	})

	t.Run("ErrAssert fails", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code did panic")
			}
		}()
		
		var err error

		Error(err, "F")
	})

	t.Run("Assert runs", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("The code did not panic")
			}
		}()
		
		IsTrue(true, "T")
	})

	t.Run("Assert fails", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did panic")
			}
		}()
		
		IsTrue(false, "F")
	})
}