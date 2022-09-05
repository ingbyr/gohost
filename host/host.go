package host

type Host interface {
	Name() string
	Content() []byte
	Desc() string
}
