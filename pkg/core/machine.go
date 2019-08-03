// Package core provides core features for building finite state machine.
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

// AvailableTransitions returns transitions available for provided Object
func (m *Machine) AvailableTransitions(o Object, args ...interface{}) ([]Transition, error) {
	return m.md.findAvailableTransitions(m.ctx, o, args...)
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

// IsInFinalState returns true if Object.Status() is a name of a final state
func (m *Machine) IsInFinalState(o Object) bool {
	for _, state := range m.md.Schema.FinalStates {
		if o.Status() == state.Name {
			return true
		}
	}
	return false
}

// AvailableStates returns all states available in machine's definition
func (m *Machine) AvailableStates() []State {
	return m.md.getAvailableStates()
}

// IsRunning returns true if object's status matches non-final state of machine
func (m *Machine) IsRunning(o Object) bool {
	for _, state := range m.AvailableStates() {
		if state.Name == o.Status() {
			return !m.IsInFinalState(o)
		}
	}
	return false
}
