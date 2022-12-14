package cron

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const OneSecond = 1*time.Second + 50*time.Millisecond
const TwoSecond = 2*OneSecond + 50*time.Millisecond

var ctx = context.Background()

func wait(wg *sync.WaitGroup) chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}

func TestAdd(t *testing.T) {
	t.Run("Works", func(t *testing.T) {
		c := NewManager(ctx, false)
		defer c.Stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd: func() {
				wg.Done()
			},
		})

		assert.NoError(t, err)

		c.Start()

		select {
		case <-time.After(OneSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}

		assert.NoError(t, c.Remove("test"))
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		c := NewManager(ctx, false)
		defer c.Stop()

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd:  func() {},
		})

		assert.NoError(t, err)

		err = c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd:  func() {},
		})

		assert.Equal(t, ErrAlreadyExists, err)

		assert.NoError(t, c.Remove("test"))
	})

	t.Run("ErrInvalidSpec", func(t *testing.T) {
		c := NewManager(ctx, false)
		defer c.Stop()

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "invalid",
			Cmd:  func() {},
		})

		assert.Error(t, err)
	})
}

func TestRemove(t *testing.T) {
	t.Run("Can remove", func(t *testing.T) {
		c := NewManager(ctx, false)
		defer c.Stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd: func() {
				wg.Done()
			},
		})

		assert.NoError(t, err)

		c.Start()

		select {
		case <-time.After(TwoSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}

		assert.NoError(t, c.Remove("test"))

		wg.Add(1)

		select {
		case <-time.After(TwoSecond):
		case <-wait(wg):
			t.Fatal("should not have run")
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		c := NewManager(ctx, false)
		defer c.Stop()

		assert.Equal(t, ErrNotFound, c.Remove("test"))
	})
}

func TestEnable(t *testing.T) {
	t.Run("Works", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd: func() {
				wg.Done()
			},
		})

		assert.NoError(t, err)

		c.Start()

		select {
		case <-time.After(OneSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}

		assert.NoError(t, c.Disable("test"))

		wg.Add(1)

		select {
		case <-time.After(OneSecond):
		case <-wait(wg):
			t.Fatal("should not have run")
		}

		assert.NoError(t, c.Enable("test"))

		wg.Add(1)

		select {
		case <-time.After(TwoSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		assert.Equal(t, ErrNotFound, c.Enable("test"))
	})

	t.Run("ErrAlreadyEnabled", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd:  func() {},
		})

		assert.NoError(t, err)

		assert.Equal(t, ErrAlreadyEnabled, c.Enable("test"))
	})
}

func TestEnabled(t *testing.T) {
	t.Run("Already enabled", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd: func() {
				wg.Done()
			},
		})

		assert.NoError(t, err)

		c.Start()

		select {
		case <-time.After(TwoSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}

		assert.Equal(t, ErrAlreadyEnabled, c.Enable("test"))

		c.Stop()
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		assert.Equal(t, ErrNotFound, c.Enable("test"))
	})
}

func TestDisabled(t *testing.T) {
	t.Run("Already disabled", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)

		err := c.Add(CronOptions{
			Name: "test",
			Spec: "@every 1s",
			Cmd: func() {
				wg.Done()
			},
		})

		assert.NoError(t, err)

		c.Start()

		select {
		case <-time.After(TwoSecond):
			t.Fatal("timed out")
		case <-wait(wg):
		}

		assert.NoError(t, c.Disable("test"))

		assert.Equal(t, ErrAlreadyDisabled, c.Disable("test"))
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		c := NewManager(ctx, true)
		defer c.Stop()

		assert.Equal(t, ErrNotFound, c.Disable("test"))
	})
}

func TestContext(t *testing.T) {
	tCtx, cancel := context.WithCancel(ctx)

	c := NewManager(tCtx, false)
	defer c.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := c.Add(CronOptions{
		Name: "test",
		Spec: "@every 1s",
		Cmd: func() {
			wg.Done()
		},
	})

	assert.NoError(t, err)

	c.Start()

	select {
	case <-time.After(OneSecond):
		t.Fatal("timed out")
	case <-wait(wg):
	}

	cancel()

	wg.Add(1)

	select {
	case <-time.After(OneSecond):
	case <-wait(wg):
		t.Fatal("should not have run")
	}
}

func TestRunAtCreation(t *testing.T) {
	c := NewManager(ctx, false)
	defer c.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	err := c.Add(CronOptions{
		Name: "test",
		Spec: "@every 1s",
		Cmd: func() {
			wg.Done()
		},
		RunNow: true,
	})

	assert.NoError(t, err)

	c.Start()

	select {
	case <-time.After(OneSecond):
		t.Fatal("timed out")
	case <-wait(wg):
	}

	assert.NoError(t, c.Remove("test"))
}
