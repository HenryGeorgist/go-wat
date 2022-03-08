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
type Link struct {
	InputDataLocation  `json:"input"`
	OutputDataLocation `json:"output"`
}

//InputDataLocations define where a model will get input from the format, parameter and the link information
type InputDataLocation struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Format    string `json:"format"`
}

//OutputDataLocations define where a model can produce output the format, parameter and the link information
type OutputDataLocation struct {
	Name                 string   `json:"name"`
	Parameter            string   `json:"parameter"`
	Format               string   `json:"format"`
	LinkInfo             LinkInfo `json:"link_info"`
	GeneratingPluginName string   `json:"generating_plugin_name"`
	GeneratingModelName  string   `json:"generating_model_name"`
}

type LinkInfo interface {
	Type() LinkType
}

type LinkType byte

const (
	CloudStore LinkType = 0
	LocalStore LinkType = 1
	Virtual    LinkType = 2
)

type LocalCSVLink struct {
	Path string `json:"Path"`
}

func (lcsv LocalCSVLink) Type() LinkType {
	return LocalStore
}
