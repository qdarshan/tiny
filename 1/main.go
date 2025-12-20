package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"
)

type result struct {
	response string
	status   string
	duration time.Duration
}

func main() {
	fmt.Println("started")

	if len(os.Args) == 2 {
		loadTest(os.Args[1], 10)
	} else if len(os.Args) == 3 {
		iterations, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: The iteration count must be a valid number.")
			os.Exit(1)
		}
		loadTest(os.Args[1], iterations)
	}

	fmt.Println("end")
}

func loadTest(url string, iterations int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	start := time.Now()

	var output map[int]result
	output = make(map[int]result)
	for i := range iterations {
		j := i
		wg.Go(func() {
			eachExecution := time.Now()
			resp, err := http.Get(url)

			if err != nil {
				fmt.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			mu.Lock()

			res := result{
				response: string(body),
				status:   resp.Status,
				duration: time.Since(eachExecution),
			}

			output[j] = res
			mu.Unlock()
		})
	}
	wg.Wait()

	fmt.Printf("total time it took execute all %d iterations in parallel: %s\n", iterations, time.Since(start))
	printOutput(output)
}

func printOutput(output map[int]result) {

	keys := make([]int, 0, len(output))

	var minDuration, maxDuration, totalTime time.Duration
	var durations = make([]time.Duration, 0, len(output))

	for key, result := range output {
		keys = append(keys, key)
		durations = append(durations, result.duration)

		if minDuration == 0 || result.duration < minDuration {
			minDuration = result.duration
		}

		if result.duration > maxDuration {
			maxDuration = result.duration
		}

		totalTime += result.duration
	}

	sort.Ints(keys)
	slices.Sort(durations)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Key\tStatus\tDuration")
	fmt.Fprintln(w, "----\t-------\t---------")

	for _, key := range keys {
		fmt.Fprintf(w, "%d\t%s\t%s\n", key, output[key].status, output[key].duration)
	}
	w.Flush()

	fmt.Println("avg time it took to execute: ", totalTime/time.Duration(len(output)))
	fmt.Println("P95: ", durations[int(float64(len(durations))*0.95)])
}
