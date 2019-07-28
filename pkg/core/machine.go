package core

import (
	"context"
	"fmt"
)

// Machine is the main structure which executes workflow
type Machine struct {
	ctx context.Context
	md  *MachineDefinition
}

// NewMachine returns new machine instance
func NewMachine(ctx context.Context, md *MachineDefinition) *Machine {
	return &Machine{ctx, md}
}

// Start sets object status to initial state
// TODO: add user, description variadic args like in findAvailableTransitions
func (m *Machine) Start(o Object) {
	o.SetStatus(m.md.Schema.InitialState.Name)
}

// CurrentState returns current state based on object's status
func (m *Machine) CurrentState(o Object) (state State, err error) {
	err = fmt.Errorf("state '%s' not found in schema", o.Status())
	for _, s := range m.md.Schema.States {
		if s.Name == o.Status() {
			state = s
			err = nil
			break
		}
	}
	return state, err
}

// AvailableTransitions returns transitions available for provided Object
func (m *Machine) AvailableTransitions(o Object, args ...interface{}) ([]Transition, error) {
	return m.md.findAvailableTransitions(m.ctx, o, args...)
}
