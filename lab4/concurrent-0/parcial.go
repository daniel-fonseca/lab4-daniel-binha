package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
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

// worker function to calculate sum and send the result via channel
func worker(path string, resultCh chan<- map[int]string) {
	_sum, err := sum(path)
	if err != nil {
		// send an empty result if there is an error
		resultCh <- nil
		return
	}

	// Send the result as a map with the sum and the file path
	resultCh <- map[int]string{_sum: path}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	sums := make(map[int][]string)
	totalSum := int64(0)
	resultCh := make(chan map[int]string)

	// Launch a goroutine for each file
	for _, path := range os.Args[1:] {
		go worker(path, resultCh)
	}

	// Collect results from the channel
	for i := 0; i < len(os.Args)-1; i++ {
		result := <-resultCh
		if result == nil {
			continue
		}

		// Process the result
		for _sum, file := range result {
			totalSum += int64(_sum)
			sums[_sum] = append(sums[_sum], file)
		}
	}

	// Print the total sum
	fmt.Println(totalSum)

	// Print files with the same sum
	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}
}
