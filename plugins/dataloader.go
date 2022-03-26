package plugins

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

//plugin implementation
//plugin helper function.
func NewHydrographScalerModel(s string) (HydrographScalerModel, error) {

	var hsm HydrographScalerModel = HydrographScalerModel{
		// Todo: FlowFrequency not unmarshalled, need to read from source
		FlowFrequency: statistics.LogPearsonIIIDistribution{
			Mean: 1.0, StandardDeviation: .01, Skew: .02, EquivalentYearsOfRecord: 10}}

	jsonFile, err := os.Open(s)
	if err != nil {
		return hsm, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return hsm, err
	}

	json.Unmarshal(jsonData, &hsm)
	return hsm, nil
}

func NewRasModel(s string) (RasModel, error) {

	var rm RasModel

	jsonFile, err := os.Open(s)
	if err != nil {
		return rm, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return rm, err
	}

	json.Unmarshal(jsonData, &rm)
	return rm, nil
}
