package editor

const Default = "vim"

type Editor interface {
	Open(filePath string) error
}

type factory = func() Editor

var factories = map[string]factory{}

func register(name string, f factory) {
	if _, ok := factories[name]; ok {
		panic("redundant editor: " + name)
	}
	factories[name] = f
}

func New(name string) (e Editor) {
	if f, ok := factories[name]; ok {
		return f()
	}
	return NewNone(name)
}
