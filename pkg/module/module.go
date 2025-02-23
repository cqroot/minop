package module

type Module interface {
	Run() (string, error)
}

type Doc interface {
	Desc() string
	Args() []Arg
}
