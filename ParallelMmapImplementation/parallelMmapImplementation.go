package parallelmmapimplementation

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/edsrzf/mmap-go"
	simpleprocess "github.com/gnana997/billion-rows-challenge/SimpleProcess"
)

type MemChunk struct {
	start int
	end   int
}

func splitMemChunks(m mmap.MMap, n int) []MemChunk {
	total := len(m)
	chunkSize := total / n
	chunks := make([]MemChunk, n)

	chunks[0].start = 0
	for i := 1; i < n; i++ {
		for j := i * chunkSize; j < i*chunkSize+50; j++ {
			if m[j] == '\n' {
				chunks[i-1].end = j
				chunks[i].start = j + 1
				break
			}
		}
	}
	chunks[n-1].end = total
	return chunks
}

func readMemChunk(ch chan map[string]*simpleprocess.SimpleMeasurements, m mmap.MMap, chunk MemChunk) {
	measurements := make(map[string]*simpleprocess.SimpleMeasurements)
	prev := chunk.start
	for i := chunk.start; i < chunk.end; i++ {
		if m[i] == '\n' {
			line := string(m[prev:i])
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
			prev = i + 1
		}
	}
	ch <- measurements
}

func ParallelMmapImplementation(filePath string) {
	defer func(t time.Time) {
		fmt.Println("Time taken:", time.Since(t))
	}(time.Now())

	maxGoroutines := min(runtime.NumCPU(), runtime.GOMAXPROCS(0))
	fmt.Println("Number of goroutines:", runtime.NumCPU(), runtime.GOMAXPROCS(0))

	dataFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer dataFile.Close()

	mmap, err := mmap.Map(dataFile, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer mmap.Unmap()

	chunks := splitMemChunks(mmap, maxGoroutines)

	ch := make(chan map[string]*simpleprocess.SimpleMeasurements, maxGoroutines)

	for i := 0; i < maxGoroutines; i++ {
		go readMemChunk(ch, mmap, chunks[i])
	}

	measurements := make(map[string]*simpleprocess.SimpleMeasurements)

	for i := 0; i < maxGoroutines; i++ {
		measurement := <-ch
		for key, value := range measurement {
			if _, ok := measurements[key]; !ok {
				measurements[key] = value
			} else {
				measurement := measurements[key]
				measurement.Min = min(measurement.Min, value.Min)
				measurement.Max = max(measurement.Max, value.Max)
				measurement.Sum += value.Sum
				measurement.Count += value.Count
			}
		}
	}

	simpleprocess.PrintMeasurements(measurements)
}
