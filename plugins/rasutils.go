package plugins

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"go-wat/config"

	"github.com/USACE/filestore"
	"github.com/USACE/mcat-ras/tools"
)

type RasBCIndex struct {
	BCLineIDX   int `json:"bcline_idx"`
	IntervalIDX int `json:"interval_idx"`
	// StepsIDX           int `json:"steps_idx"`
	HydrographStartIDX int `json:"hydrograph_start_idx"`
	HydrographStopIDX  int `json:"hydrograph_stop_idx"`
}

type RasBoundaryConditions struct {
	BCLine     string    `json:"bc_line"`
	Interval   float64   `json:"interval"`
	Steps      int       `json:"steps"`
	Hydrograph []float64 `json:"hydrograph"`
}

// Need some converter to pull this from text
// Using this as a place holder, which will fail on any model that has anything other than 1Hour
var rasIntervals map[string]float64 = map[string]float64{"1HOUR": 1}

func extractTimeInterval(s string) (float64, error) {

	rawText := strings.Trim(s, "\r")
	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return 0, fmt.Errorf("extractTimeInterval error: insufficient data from text file line")
	}

	if _, found := rasIntervals[textLineParts[1]]; !found {
		return 0, fmt.Errorf("extractTimeInterval error: unknown timestep, please add to `rasIntervals` in `ras.go`")
	}

	numericInterval := rasIntervals[textLineParts[1]]
	return numericInterval, nil

}

func extractNumberTimeSteps(s string) (int, error) {

	rawText := strings.Trim(s, "\r")

	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return 0, fmt.Errorf("extractNumberTimeSteps error: unrecognized data from text file line")
	}

	stepsNumeric, err := strconv.Atoi(strings.Trim(textLineParts[1], " "))
	if err != nil {
		return 0, err
	}

	return stepsNumeric, nil
}

func extractBCName(s string) (string, error) {

	rawText := strings.Trim(s, "\r")
	textLineParts := strings.Split(rawText, "=")
	if len(textLineParts) < 2 {
		return "", fmt.Errorf("extractBCName error: unrecognized data from text file line")
	}

	lineValues := strings.Split(textLineParts[1], ",")
	var fullBCName string = strings.Trim(lineValues[0], " ")
	for i, text := range lineValues {
		// Todo: Need to verify how BC's  are stored / nomenclature convention for this
		// currently this function strips white space from `Boundary Location` line and concatenates values with textData (skipping empty spaces)
		// i.e. `Boundary Location=White           ,Muncie          ,15696.24,        ,                ,                ,                , `
		// is returned as `White-Muncie-15696.24`
		textData := strings.Trim(text, " ")
		if i > 0 && textData != "" {
			fullBCName += "-" + textData
		}
	}
	return fullBCName, nil
}

func extractHydrograph(ss []string) ([]float64, error) {
	var start int
	var stride int = 8
	var ordinates int = 10
	var data []float64

	for _, lineValues := range ss {
		start = 0
		for i := 1; i <= ordinates; i++ {
			if len(string(lineValues)) > start+stride {

				val, err := strconv.ParseFloat(strings.TrimSpace(lineValues[start:start+stride]), 64)
				if err != nil {
					return data, err
				}
				data = append(data, val)
				start += stride
			} else {
				continue
			}
		}
	}
	return data, nil
}

// hecRasBCs is a placeholder utility funciton for reading data from models
func hecRasBCLineIndices(rm config.RasModelInfo) ([]RasBCIndex, error) {

	var rbcidx []RasBCIndex
	fs, err := filestore.NewFileStore(filestore.BlockFSConfig{})
	if err != nil {
		return rbcidx, err
	}

	modelData, err := tools.NewRasModel(rm.ProjectFilePath, fs)
	if err != nil {
		return rbcidx, err
	}

	for _, file := range modelData.Metadata.PlanFiles {

		lineBytes, err := ioutil.ReadFile(rm.BasePath + "." + file.FlowFile)
		if err != nil {
			return rbcidx, err
		}

		lines := strings.Split(string(lineBytes), "\n")

		for i, line := range lines {
			match, err := regexp.MatchString("=", line)
			if err != nil {
				return rbcidx, err
			}

			if match {
				lineData := strings.Split(line, "=")
				var bcInfo RasBCIndex

				switch lineData[0] {
				// Todo: make this work on any model, not just muncie!
				case "Boundary Location":
					nextLine := strings.Split(lines[i+1], "=")[0]

					if nextLine == "Interval" {

						stepsNumeric, err := extractNumberTimeSteps(lines[i+2])
						if err != nil {
							return rbcidx, err
						}

						bcInfo.BCLineIDX = i
						bcInfo.IntervalIDX = i + 1
						// bcInfo.StepsIDX = i + 2
						bcInfo.HydrographStartIDX = i + 3
						bcInfo.HydrographStopIDX = int(math.Ceil(float64(stepsNumeric) / 10))
						rbcidx = append(rbcidx, bcInfo)

					} else {
						continue
					}
				}
			}
		}

	}

	// Configured to use the muncie model for demo purposes only
	// grabbing the firstBC info from u01 | p04, using hard coded flows
	return rbcidx, nil

}

// hecRasBCs is a placeholder utility funciton for reading data from models
// Modify to ingest line indices
func hecRasBCs(rm config.RasModelInfo) (RasBoundaryConditions, error) {

	var rbc RasBoundaryConditions
	fs, err := filestore.NewFileStore(filestore.BlockFSConfig{})
	if err != nil {
		return rbc, err
	}

	modelData, err := tools.NewRasModel(rm.ProjectFilePath, fs)
	if err != nil {
		return rbc, err
	}

	var rbcs []RasBoundaryConditions

	for _, file := range modelData.Metadata.PlanFiles {

		lineBytes, err := ioutil.ReadFile(rm.BasePath + "." + file.FlowFile)
		if err != nil {
			return rbc, err
		}

		lines := strings.Split(string(lineBytes), "\n")

		for i, line := range lines {
			match, err := regexp.MatchString("=", line)
			if err != nil {
				return rbc, err
			}

			if match {
				lineData := strings.Split(line, "=")

				switch lineData[0] {
				// Todo: make this work on any model, not just muncie!
				case "Boundary Location":
					nextLine := strings.Split(lines[i+1], "=")[0]

					if nextLine == "Interval" {

						bcLineName, err := extractBCName(line)
						if err != nil {
							return rbc, err
						}

						numericInterval, err := extractTimeInterval(lines[i+1])
						if err != nil {
							return rbc, err
						}

						stepsNumeric, err := extractNumberTimeSteps(lines[i+2])
						if err != nil {
							return rbc, err
						}

						// Parse block of hydrograph ords from text file
						startLine := i + 3
						nTextLines := math.Ceil(float64(stepsNumeric) / 10)
						hydrograph, err := extractHydrograph(lines[startLine : startLine+int(nTextLines)])
						if err != nil {
							return rbc, err
						}

						rbc.BCLine = bcLineName
						rbc.Interval = numericInterval
						rbc.Steps = stepsNumeric
						rbc.Hydrograph = hydrograph

						rbcs = append(rbcs, rbc)

					} else {
						continue
					}
				}
			}
		}

	}

	// Configured to use the muncie model for demo purposes only
	// grabbing the firstBC info from u01 | p04, using hard coded flows
	return rbcs[0], nil

}

// Todo: Fix--this is not quite right
func rightJustify(s string, n int, fill string) string {
	if len(s) < n {
		padLevel := n - len(s)
		return strings.Repeat(fill, padLevel) + s
	} else {
		return s[:n-1]
	}
}

func hydroArrayToRasFormat(sb []byte) (string, error) {
	var blockArray []string       // Temporary holder for formatted string ordinates
	var outputBlock string        // Formatted output block for insertion into ras file
	var rasOrdinateWidth int = 8  // RAS spec
	var ordinatesPerLine int = 10 // RAS spec

	for i, lineText := range strings.Split(string(sb), "\n") {
		// example lineText = `2019-01-03 16:01:01.000000001 -0500 EST,50046.20488749354`

		lineParts := strings.Split(lineText, ",")

		if i == 0 {
			if lineParts[1] != "Flow" {
				return outputBlock, fmt.Errorf("expected Flow in header not found")
			}
		} else if len(lineParts[0]) == 0 {
			// Todo: add a data check here or elsewhere, currently skips eval from extra line at EOF
			continue

		} else if i > 0 {

			ordinate, err := strconv.ParseFloat(lineParts[1], 32)
			if err != nil {
				return outputBlock, err
			}

			formattedNumStr := strconv.FormatFloat(ordinate, 'f', -1, 64)
			justifiedNumStr := rightJustify(formattedNumStr, rasOrdinateWidth, " ")
			blockArray = append(blockArray, justifiedNumStr)
		}

	}

	var rowIDX int = 0
	for _, ordinate := range blockArray {

		if rowIDX > ordinatesPerLine-1 {
			outputBlock += ordinate + "\n"
			rowIDX = 0
		} else {
			outputBlock += ordinate
			rowIDX += 1
		}

	}
	outputBlock += "\n"
	return outputBlock, nil
}
