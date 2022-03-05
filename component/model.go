package component

type Model struct {
	ModelName     string      //e.g. plan name, project name, watershed name, basin name
	Plugin        *Computable // a reference to the associated plugin
	ModelLinkages ModelLinks  //the connections of inputs to outputs
}

type ModelLinks struct {
	Links map[InputDataLocation]OutputDataLocation
}

type InputDataLocation struct {
	Parameter   string
	Format      string
	LinkInfo    string
	SourceModel *Model
}
type OutputDataLocation struct {
	Parameter       string
	Format          string
	LinkInfo        string
	GeneratingModel *Model
}
