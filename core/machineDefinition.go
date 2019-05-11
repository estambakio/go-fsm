package core

import (
	"fmt"
	"log"
)

type Object interface {
	Status() string
}

type Param struct {
	name  string
	value interface{}
}

type Guard struct {
	name   string
	params []Param
}

type Transition struct {
	from   string
	to     string
	event  string
	guards []Guard
}

type State struct {
	name string
}

type Schema struct {
	name         string
	initialState string
	finalStates  []string
	states       []State
	transitions  []Transition
}

type Condition struct {
	name string
	f    func(o Object) bool
}

type MachineDefinition struct {
	schema     Schema
	conditions []Condition
}

func (md *MachineDefinition) getAvailableStates() []State {
	return md.schema.states
}

func (md *MachineDefinition) getConditionByName(name string) (*Condition, error) {
	for _, cond := range md.conditions {
		if cond.name == name {
			return &cond, nil
		}
	}
	return nil, fmt.Errorf("Condition with name '%s' not found", name)
}

func (md *MachineDefinition) findAvailableTransitions(o Object) []Transition {
	availableT := []Transition{}

	status := o.Status()

	for _, t := range md.schema.transitions {
		if t.from != status {
			continue
		}

		allowed := func() bool {
			for _, guard := range t.guards {
				cond, err := md.getConditionByName(guard.name)
				if err != nil {
					log.Fatalf("Condition '%s' not found", guard.name)
				}
				// TODO handle case if condition is not found
				if cond.f(o) == false {
					return false
				}
			}
			return true
		}()

		if allowed {
			availableT = append(availableT, t)
		}
	}
	return availableT
}

func (s *State) String() string {
	return s.name
}
