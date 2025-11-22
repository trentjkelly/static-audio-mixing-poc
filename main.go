package main

import (
	"fmt"
	"log"
	"os"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)


func main() {
	fmt.Println("Starting...")

	file1, file2, err := GetFileNames()
	if err != nil {
		log.Fatalf("Failed to get file names: %v\n", err)
	}

	mixedSamples, err :=MixAudioFiles(file1, file2)
	if err != nil {
		log.Fatalf("Failed to mix audio files: %v\n", err)
	}

	fmt.Println(mixedSamples)
	fmt.Println("Done!")
}

// GetFileNames retrieves the filenames from command line arguments
func GetFileNames() (string, string, error) {
	args := os.Args

	if len(args) != 3 {
		return "", "", fmt.Errorf("usage: go run main.go <file1> <file2>")
	}

	return args[1], args[2], nil
}

// MixAudioFiles mixes the two flac audio files
func MixAudioFiles(file1 string, file2 string) ([]int32, error) {
	fmt.Println("Mixing audio files...")

	samples1, err := GetSamples(file1)
	if err != nil {
		return nil, fmt.Errorf("failed to get samples from file1: %w", err)
	}

	samples2, err := GetSamples(file2)
	if err != nil {
		return nil, fmt.Errorf("failed to get samples from file2: %w", err)
	}

	mixedSamples := MixAudioSamples(samples1, samples2)

	return mixedSamples, nil
}

// GetSamples extracts samples from a flac file
func GetSamples(filePath string) ([]int32, error) {
	stream, err := flac.ParseFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flac file: %w", err)
	}

	samples := make([]int32, 0)

	for {
		frame, err := stream.ParseNext()
		if err != nil { // EOF, so just break the loop
			break
		}

		for _, subframe := range frame.Subframes {
			samples = append(samples, subframe.Samples...)
		}
	}

	return samples, nil
}

func MixAudioSamples(samples1 []int32, samples2 []int32) []int32 {
	var mixedSamples []int32

	samples1Len := len(samples1)
	samples2Len := len(samples2)

	for i := 0; i < samples1Len && i < samples2Len; i++ {
		mixedSample := (samples1[i] + samples2[i]) / 2
		mixedSamples = append(mixedSamples, mixedSample)
	}

	return mixedSamples
}

func EncodeFlacFile(newFilePath string, samples []int32) error {
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// TODO: Add the stram info metadata
	streamInfo := meta.StreamInfo{}

	encoder, err := flac.NewEncoder(file, &streamInfo)
	if err != nil {
		return fmt.Errorf("failed to create flac encoder: %w", err)
	}
	defer encoder.Close()

	return nil
}