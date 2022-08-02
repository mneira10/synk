/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mneira10/synk/internal"
	"github.com/mneira10/synk/s3Storage"
	"github.com/spf13/cobra"
)

const ASCII_NUKE string = `
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

		for _, bucketFile := range bucketFiles {
			fmt.Println("DELETING", bucketFile)
			s3Client.DeleteFile(bucketFile)
			fmt.Println("DELETED", bucketFile)
		}

		fmt.Println("Nuked everything in", config.BucketName)
		fmt.Println(ASCII_NUKE)
	},
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
