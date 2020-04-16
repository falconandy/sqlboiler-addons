package model

type Entity struct {
	FileName        string
	Package         string
	Name            string
	Fields          []EntityField
	InsertSignature []string
	UpdateSignature []string
}

type EntityField struct {
	Name string
	Type string
}
