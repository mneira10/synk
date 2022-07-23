/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mneira10/synk/internal"
	log "github.com/mneira10/synk/logger"
	"github.com/mneira10/synk/s3Storage"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Starting push...")

		config := internal.GetConfiguration(cfgFilePath)
		s3Client := s3Storage.ConfigS3(&config)

		localFiles := internal.GetFilePathsInLocalPath(cfgFilePath)
		bucketFiles := internal.GetFilePathsInBucket(s3Client)

		diffFiles := internal.GetDiffFilePaths(&localFiles, &bucketFiles)

		fmt.Println("These are the files that will be uploaded:")
		diffFileOutput := internal.PrettifyFilePaths(&diffFiles)
		fmt.Println(diffFileOutput)

		didUserConsent := getUserConsent()
		if !didUserConsent {
			os.Exit(1)
		}

		fmt.Println("Uploading files...")

		for _, filePathRelativeToCfgFolder := range diffFiles {
			localFilePath := filepath.Join(cfgFilePath, filePathRelativeToCfgFolder)
			fmt.Println("STARTING", localFilePath)
			s3Client.UploadFile(localFilePath, filePathRelativeToCfgFolder)
			fmt.Println("UPLOADED", localFilePath)
		}

		fmt.Println("Files uploaded.")
	},
}

func getUserConsent() bool {
	var userAnswer string
	fmt.Println("Do you wish to proceed? yes/no")
	fmt.Scanln(&userAnswer)

	if userAnswer != "yes" && userAnswer != "no" {
		fmt.Println("Please type 'yes' or 'no'")
		return getUserConsent()
	} else if userAnswer == "no" {
		fmt.Println("Did not upload files")
		return false
	}
	return true
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
