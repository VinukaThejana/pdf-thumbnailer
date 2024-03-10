/*
Copyright © 2024 Vinuka Kodituwakku <vinuka.t@pm.me>
*/

// Package cmd combines all the commands
package cmd

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"pdf-thumbnailer/lib"
	"strings"
	"sync"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/fatih/color"
	"github.com/gen2brain/go-fitz"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pdf-thumbnailer",
	Short: "Generate thumbnail images to the pdfs provided",
	Long:  `Create thumbnail images for pdf files to easily distinguish the pdf file without opening each file`,
	Run: func(cmd *cobra.Command, _ []string) {
		path, err := cmd.Flags().GetString("path")
		if err != nil || path == "" {
			color.Red("Provide a valid source to get PDF files")
			return
		}

		dest, err := cmd.Flags().GetString("destination")
		if err != nil || dest == "" {
			color.Red("Provide a destination to store the image thumbnails")
			return
		}

		fileInfo, err := os.Stat(path)
		if err != nil || !fileInfo.IsDir() {
			color.Red("The path provided is not a valid ")
			return
		}
		fileInfo, err = os.Stat(dest)
		if err != nil || !fileInfo.IsDir() {
			color.Red("The destination provided is not a vaid destination")
			return
		}

		files := lib.GetAllPDFFiles(path)

		var wg sync.WaitGroup
		wg.Add(len(files))

		for _, path := range files {
			go func(path string) {
				defer wg.Done()
				filename := filepath.Base(path)

				doc, err := fitz.New(path)
				if err != nil {
					logger.Error(err)
					return
				}
				defer doc.Close()

				img, err := doc.Image(0)
				if err != nil {
					logger.Error(err)
					return
				}

				f, err := os.Create(filepath.Join(dest, fmt.Sprintf("%s.png", strings.TrimSuffix(filename, filepath.Ext(path)))))
				if err != nil {
					logger.Error(err)
					return
				}
				defer f.Close()

				err = png.Encode(f, img)
				if err != nil {
					logger.Error(err)
					return
				}
			}(path)
		}

		wg.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("path", "p", "", "Path to pdf files")
	rootCmd.Flags().StringP("destination", "d", "", "Path to store the extracted thumbnails")
}
