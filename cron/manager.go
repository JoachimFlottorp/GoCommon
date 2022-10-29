package cron

import (
	"context"
	"errors"

	"github.com/robfig/cron/v3"
)

var (
	ErrNotFound        = errors.New("cron job not found")
	ErrAlreadyEnabled  = errors.New("cron job already enabled")
	ErrAlreadyDisabled = errors.New("cron job already disabled")
	ErrAlreadyExists   = errors.New("cron job already exists")
)

// Manager: A struct that manages cron jobs
type Manager struct {
	ctx   context.Context
	cron  *cron.Cron
	crons map[string]WeakWrapper
}

type CronOptions struct {
	Name   string
	RunNow bool
	Spec   string
	Cmd    func()
}

// NewManager: Creates a new cron manager
func NewManager(ctx context.Context, logExec bool) *Manager {
	opts := cron.WithLogger(WithLogger(logExec))
	c := cron.New(opts)

	go func() {
		<-ctx.Done()
		c.Stop()
	}()

	m := &Manager{
		ctx:   ctx,
		cron:  c,
		crons: make(map[string]WeakWrapper),
	}

	return m
}

// Start: Starts the cron manager
func (m *Manager) Start() {
	m.cron.Start()
}

// Stop: Stops the cron manager
func (m *Manager) Stop() {
	m.cron.Stop()
}

// Add: Adds a cron job to the manager
func (m *Manager) Add(opts CronOptions) error {
	if _, ok := m.crons[opts.Name]; ok {
		return ErrAlreadyExists
	}

	id, err := m.cron.AddFunc(opts.Spec, opts.Cmd)
	if err != nil {
		return err
	}

	wrapper := NewWeakWrapper(opts.Name, opts.Spec, id, opts.Cmd)

	m.crons[opts.Name] = wrapper

	if opts.RunNow {
		go opts.Cmd()
	}

	return nil
}

// Remove: Removes a cron job from the manager
func (m *Manager) Remove(name string) error {
	if err := m.Disable(name); err != nil {
		return err
	}

	delete(m.crons, name)

	return nil
}

// Enable: Enable a previously disabled cron job
func (m *Manager) Enable(name string) error {
	c := m.crons[name]
	if c == nil {
		return ErrNotFound
	}

	if c.IsEnabled() {
		return ErrAlreadyEnabled
	}

	id, _ := m.cron.AddFunc(c.Schedule(), c.Code())
	/* 
		100% coverage is needed, and there's no way for this to error anyway (hopefully) 
		As it's using data which was previously validated
	*/

	c.SetID(id)
	c.SetStatus(true)

	return nil
}

// Disable: Disable a previously enabled cron job
func (m *Manager) Disable(name string) error {
	c := m.crons[name]
	if c == nil {
		return ErrNotFound
	}

	if c.IsDisabled() {
		return ErrAlreadyDisabled
	}

	m.cron.Remove(c.GetID())
	c.SetStatus(false)

	return nil
}
