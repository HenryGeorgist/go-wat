package component

//Model provides interface to a specific model that a plugin can compute
type Model interface {
	ModelName() string         //e.g. plan name, project name, watershed name, basin name
	PluginName() string        // a reference to the associated plugin
	ModelLinkages() ModelLinks //the connections of inputs to outputs
}

//ModelLinks provide a way to describe how models are linked for inputs and outputs.
type ModelLinks struct {
	Links []Link `json:"links"`
}
