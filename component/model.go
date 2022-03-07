package component

//Model provides interface to a specific model that a plugin can compute
type Model interface {
	ModelName() string         //e.g. plan name, project name, watershed name, basin name
	Plugin() Computable        // a reference to the associated plugin
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

/*func (ml ModelLinks) MarshalJSON() ([]byte, error) {
	s := "{\"model_links\":{\""
	for i, o := range ml.Links {
		input, _ := json.Marshal(i)
		output, _ := json.Marshal(o)
		s += string(input) + "\":" + string(output) + ","
	}
	s = strings.TrimRight(s, ",")
	//check if the last value is a string, if so, we need to clse
	s += "}}"
	fmt.Println(s)
	return []byte(s), nil
}
*/
//InputDataLocations define where a model will get input from the format, parameter and the link information
type InputDataLocation struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Format    string `json:"format"`
}

//OutputDataLocations define where a model can produce output the format, parameter and the link information
type OutputDataLocation struct {
	Parameter       string `json:"parameter"`
	Format          string `json:"format"`
	LinkInfo        string `json:"link_info"`
	GeneratingModel *Model `json:"-"` //`json:"generating_model"`
}
