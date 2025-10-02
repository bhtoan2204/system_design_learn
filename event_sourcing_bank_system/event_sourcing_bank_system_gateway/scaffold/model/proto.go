package model

type Proto struct {
	FilePath string
	Package  string
	Module   string
	Services map[string]*Service
	Messages map[string]*Message
}
