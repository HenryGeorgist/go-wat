package component

//ProgramOrder provides a structure to a list of Computable plugins
type ProgramOrder struct {
	Plugins []Computable `json:"plugins"`
}
