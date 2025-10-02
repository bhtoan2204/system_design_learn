package model

type Handler struct {
	Name         string
	Accept       string
	Method       string
	Paths        []string
	Permission   string
	ServiceName  string
	ProtoFolder  string
	ProtoPackage string
	RequestType  string
	ResponseType string
	FolderPath   string
	Type         string
}
