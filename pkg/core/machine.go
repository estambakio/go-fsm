package core

import "context"

type Machine struct {
	ctx context.Context
	md  *MachineDefinition
}

func NewMachine(ctx context.Context, md *MachineDefinition) *Machine {
	return &Machine{ctx, md}
}

func (m *Machine) currentState(o Object) string {
	return o.Status()
}

func (m *Machine) AvailableTransitions(o Object, event Event) ([]Transition, error) {
	return m.md.findAvailableTransitions(m.ctx, o, event)
}
