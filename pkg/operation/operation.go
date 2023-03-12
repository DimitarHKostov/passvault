package operation

type Operation int

const (
	Login Operation = iota
	Logout
	Upload
	Retrieve
)

func (o Operation) String() string {
	return [...]string{"login", "logout", "upload", "retrieve"}[o]
}
