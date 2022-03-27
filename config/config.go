package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TestSettings struct {
	UserHomeDir   string `json:"user_home_dir"`
	InputDataDir  string `json:"input_data_directoy"`
	OutputDataDir string `json:"output_data_directoy"`
	HydroModel    string `json:"hydro_model_info"`
	RasModel      string `json:"ras_model_info"`
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
	userRootDir := filepath.FromSlash(ts.UserHomeDir)
	ts.UserHomeDir = userRootDir
	ts.InputDataDir = filepath.FromSlash(fmt.Sprintf("%v/%v/", userRootDir, ts.InputDataDir))
	ts.OutputDataDir = filepath.FromSlash(fmt.Sprintf("%v/%v/", userRootDir, ts.OutputDataDir))
	ts.HydroModel = filepath.FromSlash(fmt.Sprintf("%v/%v", userRootDir, ts.HydroModel))
	ts.RasModel = filepath.FromSlash(fmt.Sprintf("%v/%v", userRootDir, ts.RasModel))

	return ts, nil

}
