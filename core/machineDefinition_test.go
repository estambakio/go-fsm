package core

import (
	"testing"
)

func TestMachineDefinition_getAvailableStates(t *testing.T) {
	md := &MachineDefinition{
		schema: Schema{
			states: []State{
				State{
					name: "one",
				},
				State{
					name: "two",
				},
			},
		},
	}
	states := md.getAvailableStates()
	if len(states) != 2 {
		t.Errorf("Failed with %v", states)
	}
}

func TestMachineDefinition_getConditionByName(t *testing.T) {
	md := &MachineDefinition{
		conditions: []Condition{
			Condition{
				name: "one",
			},
			Condition{
				name: "two",
			},
		},
	}
	cond, err := md.getConditionByName("one")
	if err != nil {
		t.Errorf("Failed to get condition '%s' with error %v", "one", err)
	}
	if cond.name != "one" {
		t.Errorf("Wrong condition: expected '%s' but got '%s'", "one", cond.name)
	}

	cond, err = md.getConditionByName("notfound")
	if err == nil {
		t.Errorf("Should've failed for unknown condition but didn't: %v", cond)
	}
}

type obj struct {
	status  string
	enabled bool
}

func (o obj) Status() string {
	return o.status
}

func TestMachineDefinition_findAvailableTransitions(t *testing.T) {
	md := &MachineDefinition{
		schema: Schema{
			transitions: []Transition{
				Transition{
					from:  "a",
					to:    "b",
					event: "a->b",
					guards: []Guard{
						Guard{
							name: "isEnabled",
						},
					},
				},
				Transition{
					from:  "a",
					to:    "c",
					event: "a->c",
				},
				Transition{
					from:  "b",
					to:    "c",
					event: "b->c",
				},
			},
		},
		conditions: []Condition{
			Condition{
				name: "isEnabled",
				f: func(o Object) bool {
					return true // TODO use object somehow
				},
			},
		},
	}

	o := obj{status: "a"}
	availableT := md.findAvailableTransitions(o)
	if len(availableT) != 2 {
		t.Errorf("Failed to find available transitions for 'a': got %v", availableT)
	}
}

func TestState_String(t *testing.T) {
	state := State{
		name: "bestOfTheBest",
	}
	if state.String() != "bestOfTheBest" {
		t.Errorf("Failed with %s", state.String())
	}
}
