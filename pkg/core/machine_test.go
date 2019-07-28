package core

import (
	"context"
	"testing"
)

func TestMachine_Start(t *testing.T) {
	initialState := "k2h3ih38olhdo32ydo93hlf34"

	md := &MachineDefinition{
		Schema: Schema{
			InitialState: State{Name: initialState},
		},
	}

	machine := NewMachine(context.Background(), md)

	object := &obj{}

	machine.Start(object)

	result := object.Status()
	if result != initialState {
		t.Errorf("Failed to set initial state: expected %s but got %v", initialState, result)
	}
}

func TestMachine_CurrentState(t *testing.T) {
	status := "wfnblho439p28yr"

	object := &obj{}
	object.SetStatus(status)

	state := State{Name: status}

	// test positive path

	machine := NewMachine(
		context.Background(),
		&MachineDefinition{
			Schema: Schema{
				States: []State{state},
			},
		},
	)

	s, err := machine.CurrentState(object)
	if err != nil || s != state {
		t.Errorf("Failed to get current state: expected %s but got %v", status, s)
	}

	// test error case if obkect's status doesn't match any state in machine
	machine = NewMachine(
		context.Background(),
		&MachineDefinition{
			Schema: Schema{
				States: []State{
					State{Name: "some other name"},
				},
			},
		},
	)

	s, err = machine.CurrentState(object)
	if err == nil {
		t.Error("Should've failed: expected error but got nil")
	}
}
