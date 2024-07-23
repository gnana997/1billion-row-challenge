package simpleprocess

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type SimpleMeasurements struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func SimplProcessFunc(filePath string) {

	defer func(t time.Time) {
		fmt.Println("Time taken:", time.Since(t))
	}(time.Now())

	dataFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer dataFile.Close()

	// Store the process measurements
	measurements := make(map[string]*SimpleMeasurements)

	//Read the data from the file
	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		cityTemp := strings.Split(line, ";")
		cityName := cityTemp[0]
		temp, err := strconv.ParseFloat(cityTemp[1], 64)
		if err != nil {
			continue
		}

		// Update the measurements
		if _, ok := measurements[cityName]; !ok {
			measurements[cityName] = &SimpleMeasurements{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		} else {
			measurement := measurements[cityName]
			measurement.Min = min(measurement.Min, temp)
			measurement.Max = max(measurement.Max, temp)
			measurement.Sum += temp
			measurement.Count++
		}
	}

	PrintMeasurements(measurements)
}

func PrintMeasurements(measurements map[string]*SimpleMeasurements) {
	cityNames := make([]string, 0, len(measurements))
	for cityName, _ := range measurements {
		cityNames = append(cityNames, cityName)
	}

	slices.Sort(cityNames)

	fmt.Printf("{")
	for idx, cityName := range cityNames {
		measurement := measurements[cityName]
		avg := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f", cityName, measurement.Min, avg, measurement.Max)

		if idx < len(cityNames)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}\n")
}
