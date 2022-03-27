package plugins

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/USACE/filestore"
	"github.com/USACE/mcat-ras/tools"
)

// These routines are not currently used
// HOEVER: They will be needed to read for scraping flow data from ras models
// Todo: implement data prep features using these to extract data.

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
func hecRasBCLineIndices(rm RasModel) ([]RasBCIndex, error) {

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
func hecRasBCs(rm RasModel) (RasBoundaryConditions, error) {

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

// // Placeholder for concurrent RAS sims--Cannot include now due to events dependency on each other

// func rasWorker(jobs <-chan ContainerParams, results chan<- string) {
// 	for n := range jobs {
// 		output, err := runSimInContainerPreview(n)
// 		// output, err := RunSimInContainer(n)
// 		if err != nil {
// 			results <- fmt.Sprintf("ERROR...%v, %v", n, output)
// 		} else {
// 			results <- fmt.Sprintf("Success!...%v, %v", n, output)
// 		}

// 	}
// }

// // Specify Number of simultaneous sims
// nWorkers := 2
// nJobs := len(requiredSims)
// message := fmt.Sprintf("Processing %v models with %v Max Concurrent Simulations:", nJobs, nWorkers)
// fmt.Println(message)
// jobsChan := make(chan ContainerParams, nWorkers)
// resultsChan := make(chan string, nJobs)

// // and outputs to the results channel
// for i := 0; i < nWorkers; i++ {
// 	go worker(jobsChan, resultsChan)
// }

// // send jobs to queue
// for j := 0; j < nJobs; j++ {
// 	jobsChan <- requiredSims[j]
// }

// close(jobsChan)

// // pull the results out of the results channel
// // since it is buffered, we do not have to close it, it will close when empty
// resultsCount := 0
// for k := 0; k < nJobs; k++ {
// 	r := <-resultsChan
// 	fmt.Println("r-->", r)
// 	resultsCount++
// }
// fmt.Println("Results count:", resultsCount)
