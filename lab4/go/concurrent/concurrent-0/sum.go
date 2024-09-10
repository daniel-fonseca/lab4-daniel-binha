package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

// readFile reads a file from a given path and returns the content as a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func sum(filePath string) (int, error) {
	data, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}

	return _sum, nil
}

func worker(id int, jobs <-chan string, results chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	for filePath := range jobs {
		_sum, err := sum(filePath)
		if err != nil {
			results <- map[string]int{filePath: 0}
			fmt.Printf("Error processing file %s: %v\n", filePath, err)
		} else {
			results <- map[string]int{filePath: _sum}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers)
	results := make(chan map[string]int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		for _, path := range os.Args[1:] {
			jobs <- path
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	sums := make(map[int][]string)
	var totalSum int64

	for result := range results {
		for file, _sum := range result {
			if _sum == 0 {
				continue
			}

			totalSum += int64(_sum)
			sums[_sum] = append(sums[_sum], file)
		}
	}

	fmt.Printf("%d\n", totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Files with the same sum (%d): %v\n", sum, files)
		}
	}
}
