package main

import (
	"fmt"
	"os"

	"genplaylist/playlist"

	"github.com/spf13/cobra"
)

func main() {
	var (
		keywords []string
		format   string
		output   string
	)

	rootCmd := &cobra.Command{
		Use:   "genplaylist [folder]",
		Short: "Generate a playlist from audio files in a folder",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			folder := args[0]

			cfg, err := parseConfig(format, output, keywords)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			outFile, err := os.Create(output)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			defer outFile.Close()

			if err := playlist.Generate(os.DirFS(folder), outFile, cfg); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Playlist created: %s\n", output)
		},
	}

	rootCmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "Keywords to filter files by name")
	rootCmd.Flags().StringVarP(&format, "format", "f", "m3u", "Playlist format: m3u or pls")
	rootCmd.Flags().StringVarP(&output, "output", "o", "playlist.m3u", "Output file path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseConfig(format, output string, keywords []string) (playlist.Config, error) {
	var f playlist.Format
	switch format {
	case "m3u":
		f = playlist.FormatM3U
	case "pls":
		f = playlist.FormatPLS
	default:
		return playlist.Config{}, fmt.Errorf("unknown format %q: must be m3u or pls", format)
	}
	return playlist.Config{Format: f, Keywords: keywords}, nil
}
