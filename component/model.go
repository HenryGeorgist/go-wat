package component

type Model struct {
	ModelName     string      //e.g. plan name, project name, watershed name, basin name
	Plugin        *Computable // a reference to the associated plugin
	ModelLinkages ModelLinks  //the connections of inputs to outputs
}

type ModelLinks struct {
	Links map[string]string
}
