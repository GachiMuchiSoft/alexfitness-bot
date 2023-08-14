package user

type User interface {
	GetID() uint
	IsCoach() bool
	IsStudent() bool
	GetName() string
}

type Coach interface {
	User
}

type Student interface {
	User
}
