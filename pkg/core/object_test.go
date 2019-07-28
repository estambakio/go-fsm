package core

// concrete implementation of Object interface for tests in this package
type obj struct {
	status  string // status field for interactions with FSM
	enabled bool   // some business data
}

func (o *obj) Status() string {
	return o.status
}

func (o *obj) SetStatus(s string) {
	o.status = s
}
