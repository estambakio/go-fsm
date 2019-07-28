package core

type Object interface {
	Status() string
	SetStatus(string)
}
