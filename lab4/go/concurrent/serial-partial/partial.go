package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <file1> <file2> ...")
		return
	}

	// Map para armazenar as somas dos chunks de cada arquivo
	fileChunks := make(map[string][]int)

	// Processa cada arquivo de forma sequencial
	for _, path := range os.Args[1:] {
		chunks, err := readChunks(path)
		if err != nil {
			fmt.Printf("Erro ao processar o arquivo %s: %v\n", path, err)
			continue
		}
		fileChunks[path] = chunks
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
