package lib

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	ErrRecoverFromPanic = errors.New("recover from panic")
)

// SafeGo spawns a go-routine to run {#fn} with panic handler.
func SafeGo(fn func()) {
	go func() {
		if err := Try(fn); err != nil {
			log.Errorf("panic in go-routine: %s", err.Error())
		}
	}()
}

// Try invokes #{fn} with panic handler.
// If #{fn} causes a panic, #{ErrRecoverFromPanic} will be returned. Otherwise, nil will be returned.
func Try(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrRecoverFromPanic
		}
	}()
	fn()
	return nil
}
