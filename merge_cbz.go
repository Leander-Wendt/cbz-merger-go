package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FileInfo holds the parsed information from the filename
type FileInfo struct {
	FirstNum  int
	SecondNum float64
	Path      string
	Title     string
}

func main() {
	// Create the output directory
	outputDir := "merged"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// Get all .cbz files of the current directory
	files, err := filepath.Glob("*.cbz")

	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return
	}

	// Parse filenames and group by the first number
	fileGroups := make(map[int][]FileInfo)
	re := regexp.MustCompile(`^(\d+) - (\d+(?:\.\d+)?)\s*(.*)\.cbz$`)

	for _, file := range files {
		matches := re.FindStringSubmatch(file)
		if len(matches) < 4 {
			fmt.Printf("Skipping file with invalid format: %s\n", file)
			continue
		}

		firstNum, _ := strconv.Atoi(matches[1])
		secondNum, _ := strconv.ParseFloat(matches[2], 64)
		title := strings.TrimSpace(matches[3])

		fileGroups[firstNum] = append(fileGroups[firstNum], FileInfo{
			FirstNum:  firstNum,
			SecondNum: secondNum,
			Path:      file,
			Title:     title,
		})
	}

	// Process each group
	for firstNum, group := range fileGroups {
		// Sort by the second / volume number
		sort.Slice(group, func(i, j int) bool {
			return group[i].SecondNum < group[j].SecondNum
		})

		// Use the title from the first file in the group
		title := group[0].Title

		// Create a new zip file for the merged content
		outputFilename := filepath.Join(outputDir, fmt.Sprintf("%03d_%s.cbz", firstNum, title))
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			fmt.Printf("Failed to create output file %s: %v\n", outputFilename, err)
			continue
		}
		defer outputFile.Close()

		zipWriter := zip.NewWriter(outputFile)
		defer zipWriter.Close()

		// Merge all files in the group
		for _, fileInfo := range group {
			err := extractAndAddToZip(zipWriter, fileInfo.Path, fmt.Sprintf("%05.2f_", fileInfo.SecondNum))
			if err != nil {
				fmt.Printf("Failed to merge %s: %v\n", fileInfo.Path, err)
				continue
			}
		}

		fmt.Printf("Merged %s into %s\n", strings.Join(getFilenames(group), ", "), outputFilename)
	}
}

// extractAndAddToZip extracts the contents of a .cbz file and adds them to the zip writer with a prefix
func extractAndAddToZip(zipWriter *zip.Writer, cbzFilename string, prefix string) error {
	// Open the source .cbz file
	r, err := zip.OpenReader(cbzFilename)
	if err != nil {
		return err
	}
	defer r.Close()

	// Iterate through the files in the .cbz
	for _, f := range r.File {
		// Open the file inside the .cbz
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create a new file header in the output zip with a unique name
		header, err := zip.FileInfoHeader(f.FileInfo())
		if err != nil {
			return err
		}
		header.Name = prefix + f.Name

		// Write the file to the output zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copy the file contents
		_, err = io.Copy(writer, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

// getFilenames extracts filenames from a group of FileInfo
func getFilenames(group []FileInfo) []string {
	filenames := make([]string, len(group))
	for i, fileInfo := range group {
		filenames[i] = filepath.Base(fileInfo.Path)
	}
	return filenames
}
