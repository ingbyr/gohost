package main

type Host interface {
	Name() string
	Content() []byte
	Labels() map[string]string
}
