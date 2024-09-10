package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo %s: %v\n", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, ch chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done() // finally marca a goroutine como concluída

	data, err := readFile(filePath)
	if err != nil {
		ch <- map[string]int{filePath: 0}
		return
	}

	_sum := 0
	for _, b := range data {
		_sum += int(b)
	}
l
	ch <- map[string]int{filePath: _sum}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <file1> <file2> ...")
		return
	}

	ch := make(chan map[string]int)
	var wg sync.WaitGroup

	for _, path := range os.Args[1:] {
		wg.Add(1)
		go sum(path, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	sums := make(map[int][]string)
	var totalSum int64

	for result := range ch {
		for file, _sum := range result {
			if _sum == 0 {
				fmt.Printf("Erro ao processar o arquivo %s\n", file)
				continue
			}

			totalSum += int64(_sum)

			sums[_sum] = append(sums[_sum], file)
		}
	}

	fmt.Printf("Soma total: %d\n", totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Arquivos com a mesma soma (%d): %v\n", sum, files)
		}
	}
}
