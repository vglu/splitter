package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	version   string
	buildTime string
)

// go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" -o splitter.exe
// go run -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date -u '+%Y-%m-%dT%H:%M:%SZ')" splitter.go
func main() {

	fmt.Printf("Version: %s\n", version)
	fmt.Printf("Build Time: %s\n", buildTime)

	var input_file = flag.String("input_file", "", "Input file to split")
	var lines_per_file = flag.Int("lines_per_file", 0, "Number of lines per file")
	var useHeader = flag.Bool("use_header", true, "Use header in each file default is true")
	var output_dir = flag.String("output_dir", "", "Output directory by default same like input file")
	var help = flag.Bool("help", false, "Show help")
	flag.Parse()

	if *help || *input_file == "" || *lines_per_file == 0 {
		flag.PrintDefaults()
		fmt.Println("Example:")
		fmt.Println("splitter -input_file=source.csv -lines_per_file=1000 -use_header=true")
		return
	}

	file, err := os.Open(*input_file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	//scanner := bufio.NewScanner(file)
	reader := csv.NewReader(file)
	reader.Comma = ','       // Fields are separated by commas
	reader.LazyQuotes = true // Handle multi-line fields

	var header []string

	if *useHeader {

		// Read the header
		header, err = reader.Read()
		if err != nil {
			fmt.Println("Error reading header:", err)
			return
		}
	}
	startTime := time.Now()

	fileCount := 1
	lineCount := 0
	fieldCnt := len(header)

	baseName := strings.TrimSuffix(filepath.Base(*input_file), filepath.Ext(*input_file))
	ext := filepath.Ext(*input_file)
	if *output_dir == "" {
		*output_dir = filepath.Dir(*input_file)
	}

	var chunkCounter int
	for {
		chunkCounter++
		outputFileName := fmt.Sprintf("%s\\%s_%d%s", *output_dir, baseName, fileCount, ext)
		chunkFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("Error creating chunk file:", err)
			return
		}
		defer chunkFile.Close()

		writer := bufio.NewWriter(chunkFile)

		if *useHeader {
			for _, field := range header {
				_, err := writer.WriteString(field) // Add a newline for each record

				fieldCnt--
				if fieldCnt > 0 {
					writer.WriteString(",")
				}

				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
			writer.WriteString("\r\n")
		}

		for i := 0; i < *lines_per_file; i++ {
			fields, err := reader.Read()
			if err == io.EOF {
				writer.Flush()
				elapsedTime := time.Since(startTime)
				fmt.Printf("Processed %d lines into %d files in %s\n", lineCount, fileCount-1, elapsedTime)
				fmt.Printf("Processing completed.\n")
				fmt.Printf("Total files created: %d\n", chunkCounter)
				fmt.Printf("Total records processed: %d\n", lineCount)
				fmt.Printf("Total processing time: %v\n", elapsedTime)
				fmt.Printf("Processing speed: %.2f records/second\n", float64(lineCount)/elapsedTime.Seconds())
				return // End of file
			} else if err != nil {
				fmt.Println("Error reading record:", err)
				return
			}
			fieldCnt = len(fields)
			for _, field := range fields {

				// Check if the record contains a newline character
				if strings.Contains(field, "\n") {
					field = fmt.Sprintf("\"%s\"", field) // Enclose the line in double quotes
				}

				_, err := writer.WriteString(field) // Add a newline for each record

				fieldCnt--
				if fieldCnt > 0 {
					writer.WriteString(",")
				}

				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
			writer.WriteString("\r\n")
			lineCount++
		}
		writer.Flush()
		fileCount++
	}
}
