package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type RasModelInfo struct {
	ProjectFilePath string `json:"project_file_path"`
	UFilePath       string `json:"ufile"`
	BasePath        string `json:"base_path"`
}

type TestSettings struct {
	InputDataDir  string       `json:"input_data_directoy"`
	OutputDataDir string       `json:"output_data_directoy"`
	RasModel      RasModelInfo `json:"ras_model_info"`
}

func LoadTestSettings() (TestSettings, error) {

	var ts TestSettings
	jsonFile, err := os.Open("../config/test-config.json")
	if err != nil {
		return ts, nil
	}

	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return ts, err
	}

	json.Unmarshal(jsonData, &ts)
	return ts, nil

}
