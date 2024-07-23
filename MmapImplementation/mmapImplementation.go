package mmapimplementation

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/edsrzf/mmap-go"
	simpleprocess "github.com/gnana997/billion-rows-challenge/SimpleProcess"
)

func MmapImplementationFunc(filePath string) {
	defer func(t time.Time) {
		fmt.Println("Time taken:", time.Since(t))
	}(time.Now())

	dataFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	// Load the file into memory using mmap
	mmap, err := mmap.Map(dataFile, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer mmap.Unmap()

	measurements := make(map[string]*simpleprocess.SimpleMeasurements)

	// run throught the mmap to process the data
	prev := 0
	for i := 0; i < len(mmap); {
		if mmap[i] == '\n' {
			// process the line
			line := string(mmap[prev:i])
			cityTemp := strings.Split(line, ";")
			temp, err := strconv.ParseFloat(cityTemp[1], 64)
			if err != nil {
				continue
			}
			cityName := cityTemp[0]

			// Update the measurements
			if _, ok := measurements[cityName]; !ok {
				measurements[cityName] = &simpleprocess.SimpleMeasurements{
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

			// reset the prev
			prev = i + 1
		}
		i++
	}

	// Print the results
	for cityName, measurement := range measurements {
		avg := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f", cityName, measurement.Min, avg, measurement.Max)
	}
}
