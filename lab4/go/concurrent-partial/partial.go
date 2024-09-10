package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Define o tamanho de cada chunk em bytes
const chunkSize = 1024

// readChunks lê um arquivo e retorna a soma de cada chunk em um slice.
func readChunks(filePath string) ([]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo %s: %v\n", filePath, err)
		return nil, err
	}
	defer file.Close()

	var chunkSums []int
	buffer := make([]byte, chunkSize)

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}

		// Calcula a soma de cada chunk
		chunkSum := 0
		for _, b := range buffer[:bytesRead] {
			chunkSum += int(b)
		}
		chunkSums = append(chunkSums, chunkSum)
	}

	return chunkSums, nil
}

// similarity calcula a similaridade entre dois arquivos com base nas somas dos chunks.
func similarity(chunks1, chunks2 []int) float64 {
	matches := 0
	minLength := len(chunks1)
	if len(chunks2) < minLength {
		minLength = len(chunks2)
	}

	// Compara as somas dos chunks
	for i := 0; i < minLength; i++ {
		if chunks1[i] == chunks2[i] {
			matches++
		}
	}

	// Calcula a porcentagem de similaridade
	return float64(matches) / float64(minLength) * 100
}

// worker function para ler chunks de arquivos e enviar resultados via canal
func worker(path string, resultCh chan<- map[string][]int, wg *sync.WaitGroup) {
	defer wg.Done()

	chunks, err := readChunks(path)
	if err != nil {
		fmt.Printf("Erro ao processar o arquivo %s: %v\n", path, err)
		resultCh <- nil
		return
	}
	// Envia a soma dos chunks pelo canal
	resultCh <- map[string][]int{path: chunks}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <file1> <file2> ...")
		return
	}

	// Canal para armazenar resultados e WaitGroup para sincronização
	resultCh := make(chan map[string][]int, len(os.Args)-1)
	var wg sync.WaitGroup

	// Inicia uma goroutine para cada arquivo
	for _, path := range os.Args[1:] {
		wg.Add(1)
		go worker(path, resultCh, &wg)
	}

	// Fecha o canal quando todas as goroutines terminarem
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Map para armazenar as somas dos chunks de cada arquivo
	fileChunks := make(map[string][]int)

	// Coleta os resultados dos chunks conforme as goroutines vão enviando
	for result := range resultCh {
		if result != nil {
			for path, chunks := range result {
				fileChunks[path] = chunks
			}
		}
	}

	// Compara os arquivos e imprime a similaridade
	paths := os.Args[1:]
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			sim := similarity(fileChunks[paths[i]], fileChunks[paths[j]])
			fmt.Printf("Similaridade entre %s e %s: %.5f%%\n", filepath.Base(paths[i]), filepath.Base(paths[j]), sim)
		}
	}
}
