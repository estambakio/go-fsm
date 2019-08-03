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

func TestMachine_AvailableTransitions(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			Transitions: []Transition{
				Transition{From: "a", To: "b", Event: "a->b"},
				Transition{From: "b", To: "c", Event: "b->c"},
				Transition{From: "b", To: "d", Event: "b->d"},
				Transition{From: "c", To: "d", Event: "c->d"},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	object := &obj{status: "b"}

	result, err := machine.AvailableTransitions(object, Event("b->c"))
	if len(result) != 1 || err != nil {
		t.Errorf("Failed to get 1 available transition: received %v and %v", result, err)
	}

	result, err = machine.AvailableTransitions(object)
	if len(result) != 2 || err != nil {
		t.Errorf("Failed to get 2 available transitions: received %v and %v", result, err)
	}

	// if object is in unknown state then return empty slice
	object.status = "notInList"
	result, err = machine.AvailableTransitions(object)
	if len(result) != 0 || err != nil {
		t.Errorf("Failed to get 0 available transitions: received %v and %v", result, err)
	}
}

func TestMachine_IsInFinalState(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			FinalStates: []State{
				State{Name: "one"},
				State{Name: "two"},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	tests := []struct {
		status   string
		expected bool
	}{
		{status: "one", expected: true},
		{status: "two", expected: true},
		{status: "three", expected: false},
	}

	for i, test := range tests {
		object := &obj{status: test.status}
		result := machine.IsInFinalState(object)
		if result != test.expected {
			t.Errorf("Expected %v for test %d but got %v", test.expected, i, result)
		}
	}
}

func TestMachine_AvailableStates(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			States: []State{
				State{Name: "one"},
				State{Name: "two"},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	expected := len(md.Schema.States)
	result := len(machine.AvailableStates())
	if result != expected {
		t.Errorf("Expected %d but got %v", expected, result)
	}
}

func TestMachine_IsRunning(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			FinalStates: []State{
				State{Name: "i_am_final"},
			},
			States: []State{
				State{Name: "one"},
				State{Name: "two"},
				State{Name: "i_am_final"},
				State{Name: "three"},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	tests := []struct {
		status   string
		expected bool
	}{
		{status: "one", expected: true},
		{status: "two", expected: true},
		{status: "i_am_final", expected: false},
		{status: "not_in_schema", expected: false},
	}

	for i, test := range tests {
		object := &obj{status: test.status}
		result := machine.IsRunning(object)
		if result != test.expected {
			t.Errorf("Expected %v for test %d but got %v", test.expected, i, result)
		}
	}
}
