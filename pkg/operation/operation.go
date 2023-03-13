package operation

type Operation int

const (
	Login Operation = iota
	Logout
	Save
	Retrieve
)

func (o Operation) String() string {
	return [...]string{"login", "logout", "save", "retrieve"}[o]
}
