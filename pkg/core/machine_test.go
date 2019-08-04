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

func TestMachine_Can(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			States: []State{State{Name: "a"}, State{Name: "b"}, State{Name: "c"}},
			Transitions: []Transition{
				Transition{From: "a", To: "b", Event: Event("a->b")},
				Transition{From: "b", To: "c", Event: Event("b->c")},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	tests := []struct {
		event    string
		status   string
		expected bool
	}{
		{event: "a->b", status: "a", expected: true},
		{event: "a->b", status: "b", expected: false},
		{event: "b->c", status: "a", expected: false},
		{event: "b->c", status: "b", expected: true},
	}

	for i, test := range tests {
		object := &obj{status: test.status}
		result := machine.Can(object, Event(test.event))
		if result != test.expected {
			t.Errorf("Expected %v for test %d but got %v", test.expected, i, result)
		}
	}
}

func TestMachine_SendEvent(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			States: []State{State{Name: "a"}, State{Name: "b"}, State{Name: "c"}},
			Transitions: []Transition{
				Transition{From: "a", To: "b", Event: Event("a->b")},
				Transition{From: "b", To: "c", Event: Event("b->c")},
				Transition{From: "c", To: "d", Event: Event("c->d")},
			},
		},
	}

	machine := NewMachine(context.Background(), md)

	object := &obj{status: "a"}

	steps := []struct {
		event       string
		statusAfter string
	}{
		{event: "a->b", statusAfter: "b"},
		{event: "b->c", statusAfter: "c"},
		{event: "c->d", statusAfter: "d"},
	}
	for i, step := range steps {
		err := machine.SendEvent(object, Event(step.event))
		if err != nil || object.Status() != step.statusAfter {
			t.Errorf("Failed test %d: expected status: %v; real err: %v, real status: %v",
				i,
				step.statusAfter,
				err,
				object.Status(),
			)
		}
	}

	// send an event which is not available for current state ("d")
	object.SetStatus("d") // just to prevent tests from breaking in case of changes
	err := machine.SendEvent(object, Event("d->e"))
	if err == nil || object.Status() != "d" {
		t.Errorf("Expected error and status 'd'; got error %v and status %v", err, object.Status())
	}

	// add second transition for state 'c' with the same event 'c->d' - this is error
	md.Schema.Transitions = append(md.Schema.Transitions, Transition{From: "c", To: "e", Event: Event("c->d")})
	object.SetStatus("c")
	err = machine.SendEvent(object, Event("c->d"))
	if err == nil || object.Status() != "c" {
		t.Errorf("Expected error and status 'c'; got error %v and status %v", err, object.Status())
	}
}
