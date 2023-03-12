package operation

type Operation int

const (
	Login Operation = iota
	Logout
	Refresh
	Upload
	Retrieve
)

func (o Operation) String() string {
	return [...]string{"login", "logout", "refresh", "upload", "retrieve"}[o]
}
