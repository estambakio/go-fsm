package core

import (
	"context"
	"fmt"
)

// Param describes a single param for guard's condition function
type Param struct {
	Name  string
	Value interface{}
}

// Guard holds configuration for condition call
type Guard struct {
	Name   string
	Params []Param
}

// Event is a reason for transition
type Event string

// Transition is a single path between two states
type Transition struct {
	From   string
	To     string
	Event  Event
	Guards []Guard
}

// State marks a node in workflow's graph
type State struct {
	Name string
}

// Schema is a workflow configuration, TODO: add serialize/parse methods (to/from JSON)
type Schema struct {
	Name         string
	InitialState State
	FinalStates  []State
	States       []State
	Transitions  []Transition
}

// Condition wraps a function which defines if certain condition is passed for provided object or not
type Condition struct {
	Name string
	F    func(ctx context.Context, o Object) bool
}

// MachineDefinition is a configuration for finite states machine
type MachineDefinition struct {
	Schema     Schema
	Conditions []Condition
}

func (md *MachineDefinition) getAvailableStates() []State {
	return md.Schema.States
}

func (md *MachineDefinition) getConditionByName(name string) (*Condition, error) {
	for _, cond := range md.Conditions {
		if cond.Name == name {
			return &cond, nil
		}
	}
	return nil, fmt.Errorf("Condition with name '%s' not found", name)
}

func (md *MachineDefinition) findAvailableTransitions(ctx context.Context, o Object, args ...interface{}) ([]Transition, error) {
	transitions := []Transition{}

	// handle dynamic assigning of optional args based on passed types (yay arbitrary order)
	var event Event

	for _, arg := range args {
		switch arg := arg.(type) {
		case Event:
			event = arg
		default:
			return nil, fmt.Errorf("unknown type %T, value %v in findAvailableTransitions call", arg, arg)
		}
	}

	for _, t := range md.Schema.Transitions {
		if t.From != o.Status() {
			continue
		}

		// if event does matter for search then narrow down transitions to only those which contain this event
		if event != "" && t.Event != event {
			continue
		}

		allowed, err := func() (bool, error) {
			for _, guard := range t.Guards {
				cond, err := md.getConditionByName(guard.Name)
				if err != nil {
					return false, err
				}
				if cond.F(ctx, o) == false {
					return false, nil
				}
			}
			return true, nil
		}()

		if err != nil {
			return nil, err
		}

		if allowed {
			transitions = append(transitions, t)
		}
	}
	return transitions, nil
}
