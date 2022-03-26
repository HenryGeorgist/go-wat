package plugins

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type FileLineID struct {
	Start int `json:"start"`
	Stop  int `json:"stop"`
}

type RasBCLine struct {
	BCName            string     `json:"bc_name"`
	Type              string     `json:"flow"`
	OridnatesLocation FileLineID `json:"ordinate_location"`
}

type RasFlowFile struct {
	Path    string      `json:"path"`
	BCLines []RasBCLine `json:"bclines"`
}

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

func newUfile(existingFlowFile, newFlowData string, linestart, linestop int) (string, error) {
	fileBytes, err := ioutil.ReadFile(existingFlowFile)
	if err != nil {
		return "", err
	}

	flowDataLines := strings.Split(string(fileBytes), "\n")
	newFlowFileData := ""
	for i, line := range flowDataLines {
		if i < linestart-1 {
			newFlowFileData += line
		} else if i == linestart {
			newFlowFileData += newFlowData
		} else if i > linestop {
			newFlowFileData += line
		} else {
			continue
		}

	}
	return newFlowFileData, nil

}
