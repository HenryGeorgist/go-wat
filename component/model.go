package component

//Model provides interface to a specific model that a plugin can compute
type Model interface {
	ModelName() string         //e.g. plan name, project name, watershed name, basin name
	Plugin() Computable        // a reference to the associated plugin
	ModelLinkages() ModelLinks //the connections of inputs to outputs
}

//ModelLinks provide a way to describe how models are linked for inputs and outputs.
type ModelLinks struct {
	Links map[InputDataLocation]OutputDataLocation
}

//InputDataLocations define where a model will get input from the format, parameter and the link information
type InputDataLocation struct {
	Parameter   string
	Format      string
	LinkInfo    string
	SourceModel *Model
}

//OutputDataLocations define where a model can produce output the format, parameter and the link information
type OutputDataLocation struct {
	Parameter       string
	Format          string
	LinkInfo        string
	GeneratingModel *Model
}
