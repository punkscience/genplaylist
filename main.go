package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var keywords []string

	var rootCmd = &cobra.Command{
		Use:   "genplaylist [folder]",
		Short: "Generate an m3u playlist from audio files in a folder",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			folder := args[0]
			playlistFile := "playlist.m3u"

			files, err := findAudioFiles(folder, keywords)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			err = writePlaylist(playlistFile, files)
			if err != nil {
				fmt.Printf("Error writing playlist: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Playlist created: %s\n", playlistFile)
		},
	}

	rootCmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "List of keywords to filter files")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func findAudioFiles(folder string, keywords []string) ([]string, error) {
	var audioFiles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isAudioFile(path) {
			if len(keywords) == 0 || containsKeyword(path, keywords) {
				audioFiles = append(audioFiles, path)
			}
		}
		return nil
	})
	return audioFiles, err
}

func isAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".mp3" || ext == ".flac" || ext == ".wav"
}

func containsKeyword(path string, keywords []string) bool {
	lowerPath := strings.ToLower(path)
	for _, keyword := range keywords {
		if strings.Contains(lowerPath, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func writePlaylist(filename string, files []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, file := range files {
		_, err := f.WriteString(file + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
