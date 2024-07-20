package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// CityStats holds statistics for a city
type CityStats struct {
	MinTemp float64
	MaxTemp float64
	AvgTemp float64
	Count   int
	Total   float64
	RWMutex sync.RWMutex
}

// processChunk processes a chunk of CSV rows
func processChunk(rows [][]string, wg *sync.WaitGroup, ch chan<- CityTemperature) {
	defer wg.Done()
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		cityName := row[0]
		temp, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			fmt.Println("Error parsing temperature:", err)
			continue
		}
		ch <- CityTemperature{CityName: cityName, Temp: temp}
	}
}

// updateCityStats updates the statistics for a given city
func updateCityStats(cityStatsMap *sync.Map, result CityTemperature) {
	statsInterface, _ := cityStatsMap.LoadOrStore(result.CityName, &CityStats{
		MinTemp: result.Temp,
		MaxTemp: result.Temp,
		AvgTemp: result.Temp,
		Count:   1,
		Total:   result.Temp,
	})

	stats := statsInterface.(*CityStats)

	stats.RWMutex.Lock()
	defer stats.RWMutex.Unlock()

	if result.Temp < stats.MinTemp {
		stats.MinTemp = result.Temp
	}
	if result.Temp > stats.MaxTemp {
		stats.MaxTemp = result.Temp
	}
	stats.Total += result.Temp
	stats.Count++
	stats.AvgTemp = stats.Total / float64(stats.Count)
}

// consumeResults processes messages from the channel
func consumeResults(results <-chan CityTemperature, cityStatsMap *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	d := time.Duration(100) * time.Millisecond
	t := time.NewTimer(d)

	for {
		select {
		case result, ok := <-results:
			if !ok {
				continue
			}
			count++
			updateCityStats(cityStatsMap, result)
			if count%1000 == 0 {
				log.Println("Processed", count, "records")
			}
			if !t.Stop() {
				<-t.C
			}
			t.Reset(d)
		case <-t.C:
			// Break the loop if no new messages for 50ms
			return
		}
	}
}

func main() {
	defer func(t time.Time) {
		fmt.Println("Time taken:", time.Since(t))
	}(time.Now())

	file, err := os.Open("data/weather_stations.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	// Channel to collect results
	results := make(chan CityTemperature, 1000000)

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Map to store city statistics
	cityStatsMap := &sync.Map{}

	consumerSetup := false

	// Process CSV file in chunks
	const chunkSize = 1000000
	for {
		// Read a chunk of rows
		rows, err := readN(reader, chunkSize)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error reading rows:", err)
			break
		}

		// Launch a goroutine to process the chunk
		wg.Add(1)
		go processChunk(rows, &wg, results)

		if !consumerSetup {
			go func() {
				const numConsumers = 100
				for i := 0; i < numConsumers; i++ {
					wg.Add(1)
					go consumeResults(results, cityStatsMap, &wg)
				}
			}()
			consumerSetup = true
		}
	}

	// Wait for all consumer goroutines to finish
	wg.Wait()
	close(results)

	// Print city statistics
	count := 0
	cityStatsMap.Range(func(key, value interface{}) bool {
		if count >= 10 {
			return false
		}
		city := key.(string)
		stats := value.(*CityStats)

		stats.RWMutex.RLock()
		fmt.Printf("%s - Min: %.2f, Max: %.2f, Avg: %.2f\n", city, stats.MinTemp, stats.MaxTemp, stats.AvgTemp)
		stats.RWMutex.RUnlock()

		count++
		return true
	})
}

// ReadN reads n lines from the CSV reader
func readN(reader *csv.Reader, n int) ([][]string, error) {
	var records [][]string
	for i := 0; i < n; i++ {
		record, err := reader.Read()
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}
	return records, nil
}
