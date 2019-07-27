package core

import (
	"context"
	"fmt"
)

type Param struct {
	Name  string
	Value interface{}
}

type Guard struct {
	Name   string
	Params []Param
}

type Event string

type Transition struct {
	From   string
	To     string
	Event  Event
	Guards []Guard
}

type State struct {
	Name string
}

type Schema struct {
	Name         string
	InitialState State
	FinalStates  []State
	States       []State
	Transitions  []Transition
}

type Condition struct {
	Name string
	F    func(ctx context.Context, o Object) bool
}

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
	availableT := []Transition{}

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

		// if event does matter for search then narrow transitions to only containing this event
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
			availableT = append(availableT, t)
		}
	}
	return availableT, nil
}
