/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mneira10/synk/internal"
	log "github.com/mneira10/synk/logger"
	"github.com/mneira10/synk/s3Storage"
	"github.com/spf13/cobra"
)

const ASCII_NUKE = `
      _.-^^---....,,--_
   _--                  --_
  <                        >)
  |                         |
   \._                   _./
      '''--. . , ; .--'''
            | |   |
         .-=||  | |=-.
         '-=#$%&%$#=-'
            | ;  :|
   _____.,-#%&$@%#&#~,._____`

const NUM_CONCURRENT_DELETES = 10

type nukeData struct {
	bucketFilePath string
}

type nukeResultData struct {
	nukeData nukeData
	err      error
}

// nukeCmd represents the nuke command
var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "This will permanently delete all of the files in the specified bucket.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.GetConfiguration(cfgFilePath)
		fmt.Println("WARNING: This will delete ALL of the files in the",
			config.BucketName, "bucket.")

		fmt.Println()
		fmt.Println("If you are sure you want to do this, type the bucket name:")
		var userAnswer string
		fmt.Scanln(&userAnswer)

		if userAnswer != config.BucketName {
			fmt.Println("Aborting nuke.")
			os.Exit(1)
		}

		fmt.Println("Nuking all files in", config.BucketName)
		s3Client := s3Storage.ConfigS3(&config)
		bucketFiles := internal.GetFilePathsInBucket(s3Client)

		numFilesToNuke := len(bucketFiles)
		filesToNuke := make(chan nukeData, numFilesToNuke)
		results := make(chan nukeResultData, numFilesToNuke)

		for i := 0; i < NUM_CONCURRENT_DELETES; i++ {
			go nukeFileWorker(filesToNuke, results, s3Client)
		}

		for _, bucketFile := range bucketFiles {
			nukeData := nukeData{bucketFilePath: bucketFile}
			filesToNuke <- nukeData
		}

		close(filesToNuke)

		var nukeErrorFiles []string

		for i := 0; i < numFilesToNuke; i++ {
			nukeResult := <-results
			if nukeResult.err != nil {
				nukeErrorFiles = append(
					nukeErrorFiles,
					nukeResult.nukeData.bucketFilePath,
				)

				log.WithFields(log.Fields{
					"error": nukeResult.err,
					"file":  nukeResult.nukeData.bucketFilePath,
				}).Error("Could not nuke file")
			}
		}

		if len(nukeErrorFiles) != 0 {
			errorFilesOutput := internal.PrettifyFilePaths(&nukeErrorFiles)
			fmt.Println("Could not nuke some files:")
			fmt.Println(errorFilesOutput)
			os.Exit(1)
		} else {
			fmt.Println("Nuked everything in", config.BucketName)
			fmt.Println(ASCII_NUKE)
		}

		// for _, bucketFile := range bucketFiles {
		// 	fmt.Println("DELETING", bucketFile)
		// 	s3Client.DeleteFile(bucketFile)
		// 	fmt.Println("DELETED", bucketFile)
		// }

	},
}

func nukeFileWorker(filesToNuke <-chan nukeData, results chan<- nukeResultData, s3Client *s3Storage.S3Object) {
	for dataToNuke := range filesToNuke {
		fmt.Println("DELETING", dataToNuke.bucketFilePath)
		err := s3Client.DeleteFile(dataToNuke.bucketFilePath)
		fmt.Println("DELETED", dataToNuke.bucketFilePath)

		nukeResult := nukeResultData{nukeData: dataToNuke}

		if err != nil {
			nukeResult.err = err
		}
		results <- nukeResult
	}
}

func init() {
	rootCmd.AddCommand(nukeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nukeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nukeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
