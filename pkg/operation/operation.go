package operation

type Operation int

const (
	Login Operation = iota
	Save
	Retrieve
)

func (o Operation) String() string {
	return [...]string{"login", "save", "retrieve"}[o]
}
