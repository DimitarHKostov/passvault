package operation

type Operation int

const (
	Login Operation = iota
	Save
	Retrieve
	Update
)

func (o Operation) String() string {
	return [...]string{"login", "save", "retrieve", "update"}[o]
}
