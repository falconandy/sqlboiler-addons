package model

type FunctionSignature struct {
	Name       string
	Parameters []FunctionParameter
	Return     string
}

type FunctionParameter struct {
	Names []string
	Type  string
}
