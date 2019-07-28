package core

// Object is an interface for business object which is a subject of workflow
type Object interface {
	Status() string
	SetStatus(string)
}
