package cron

import "github.com/robfig/cron/v3"

// WeakWrapper is a wrapper for a specified cron job
type WeakWrapper interface {
	// GetID: Returns the ID of the cron job
	GetID() cron.EntryID
	// SetID: Sets the ID of the cron job
	SetID(id cron.EntryID)
	// Schedule: The schedule of the cron job
	Schedule() string
	// Code: The code of the cron job
	Code() func()
	// IsEnabled: Whether or not the cron job is enabled
	IsEnabled() bool
	// IsDisabled: Whether or not the cron job is disabled
	IsDisabled() bool
	// SetStatus: Sets the status of the cron job
	SetStatus(status bool)
}

type weakWrapper struct {
	name string
	schedule string
	id cron.EntryID
	cmd func()
	enabled bool
}

func NewWeakWrapper(name, schedule string, id cron.EntryID, cmd func()) WeakWrapper {
	w := &weakWrapper{
		name: name,
		schedule: schedule,
		id: id,
		cmd: cmd,
		enabled: true,
	}

	return w
}

func (w *weakWrapper) GetID() cron.EntryID {
	return w.id
}

func (w *weakWrapper) SetID(id cron.EntryID) {
	w.id = id
}

func (w *weakWrapper) Schedule() string {
	return w.schedule
}

func (w *weakWrapper) Code() func() {
	return w.cmd
}

func (w *weakWrapper) IsEnabled() bool {
	return w.enabled
}

func (w *weakWrapper) IsDisabled() bool {
	return !w.enabled
}

func (w *weakWrapper) SetStatus(status bool) {
	w.enabled = status
}
