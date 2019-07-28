package core

import (
	"context"
	"testing"
)

func TestMachineDefinition_getAvailableStates(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			States: []State{
				State{Name: "one"},
				State{Name: "two"},
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
		Conditions: []Condition{
			Condition{Name: "one"},
			Condition{Name: "two"},
		},
	}
	cond, err := md.getConditionByName("one")
	if err != nil {
		t.Errorf("Failed to get condition '%s' with error %v", "one", err)
	}
	if cond.Name != "one" {
		t.Errorf("Wrong condition: expected '%s' but got '%s'", "one", cond.Name)
	}

	cond, err = md.getConditionByName("notfound")
	if err == nil {
		t.Errorf("Should've failed for unknown condition but didn't: %v", cond)
	}
}

func TestMachineDefinition_findAvailableTransitions(t *testing.T) {
	md := &MachineDefinition{
		Schema: Schema{
			Transitions: []Transition{
				Transition{From: "a", To: "b", Event: "a->b",
					Guards: []Guard{
						Guard{Name: "isEnabled"},
					},
				},
				Transition{From: "a", To: "c", Event: "a->c"},
				Transition{From: "b", To: "c", Event: "b->c"},
			},
		},
		Conditions: []Condition{
			Condition{
				Name: "isEnabled",
				F: func(ctx context.Context, o Object) bool {
					return o.(*obj).enabled
				},
			},
		},
	}

	ctx := context.Background()

	availableT, err := md.findAvailableTransitions(ctx, &obj{status: "a", enabled: true})
	if err != nil || len(availableT) != 2 {
		t.Errorf("Failed to find available transitions for 'a': got %v", availableT)
	}

	availableT, err = md.findAvailableTransitions(ctx, &obj{status: "a", enabled: false})
	if err != nil || len(availableT) != 1 {
		t.Errorf("Failed to find available transitions for 'a': got %v", availableT)
	}

	// if condition is not found findAvailableTransitions should return error
	md = &MachineDefinition{
		Schema: Schema{
			Transitions: []Transition{
				Transition{
					From:  "a",
					To:    "b",
					Event: "a->b",
					Guards: []Guard{
						Guard{
							Name: "isEnabled",
						},
					},
				},
			},
		},
	}

	availableT, err = md.findAvailableTransitions(ctx, &obj{status: "a"})
	if err == nil {
		t.Errorf("Condition not found, should've returned error, but returned %v", availableT)
	}
}
