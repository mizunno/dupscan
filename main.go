package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"crypto/md5"
	"encoding/hex"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func hashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}

func md5Hash(input io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, input); err != nil {
		log.Fatal(err)
		return "", err
	}

	return hashToString(h.Sum(nil)), nil
}

func worker(files <-chan string, results chan<- map[string][]string) {
	workerDuplicates := make(map[string][]string)

	// Open file
	for file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		// Compute MD5 hash
		hash, err := md5Hash(f)

		// Instead of closing the file with defer, we close it here to avoid
		// keeping the file open for too long
		f.Close()

		workerDuplicates[hash] = append(workerDuplicates[hash], file)
	}

	results <- workerDuplicates
}

func getPotentialDuplicates(dir string, minsize int) map[int64][]string {
	// Walk directory and find potencial duplicates grouping by size
	potentialDuplicates := make(map[int64][]string)

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.Size() < int64(minsize) {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		potentialDuplicates[info.Size()] = append(potentialDuplicates[info.Size()], path)

		return nil
	})

	return potentialDuplicates
}

func writeReport(duplicates [][]string) {

	now := time.Now().Format("20060102150405")

	path := fmt.Sprintf("duplicates-%s.txt", now)

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, v := range duplicates {
		for _, file := range v {
			if _, err := f.WriteString(file + ","); err != nil {
				log.Fatal(err)
			}
		}

		// Remove last comma
		f.Seek(-1, 1)

		_, err := f.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Report written to %s", path)
}

func main() {
	// Define flags
	dir := flag.String("dir", ".", "Directory to scan")
	minsize := flag.Int("minsize", 0, "Minimum file size to consider (bytes)")
	workers := flag.Int("workers", runtime.NumCPU(), "Number of workers")

	flag.Parse()

	// Get potential duplicates
	// Potential duplicates are files with the same size
	potentialDuplicates := getPotentialDuplicates(*dir, *minsize)

	// Files channel will be used to send files to workers
	files := make(chan string, 100)
	// Results channel will be used to collect results from workers
	results := make(chan map[string][]string, *workers)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(files, results)
		}()
	}

	// Send files to workers
	go func() {
		for _, v := range potentialDuplicates {
			// Skip groups with only one file
			if len(v) < 2 {
				continue
			}

			for _, file := range v {
				files <- file
			}
		}

		close(files)
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results from workers
	aggResults := make(map[string][]string)
	for workerResults := range results {
		for k, v := range workerResults {
			aggResults[k] = append(aggResults[k], v...)
		}
	}

	// Filter duplicates
	duplicates := make([][]string, 0)
	for _, v := range aggResults {
		if len(v) > 1 {
			duplicates = append(duplicates, v)
		}
	}

	writeReport(duplicates)
}
