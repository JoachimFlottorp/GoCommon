package cron

import "testing"

// 100% coverage is really important ;-)

func TestInfo(t *testing.T) {
	l.Info(l{}, "something")
}

func TestError(t *testing.T) {
	l.Error(l{}, nil, "something")
}
