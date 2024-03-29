/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mneira10/synk/internal"
	log "github.com/mneira10/synk/logger"
	"github.com/mneira10/synk/s3Storage"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all files in the s3 compatible storage",
	Long:  `This will list all of the existing files in the bucket specified in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Listing files in bucket")

		config := internal.GetConfiguration(cfgFilePath)
		s3Client := s3Storage.ConfigS3(&config)

		fmt.Println("Bucket:", s3Client.BucketName)
		fmt.Println("Url:", s3Client.Url)

		objects := s3Client.ListObjects()

		numberOnly, _ := cmd.Flags().GetBool("numberOnly")

		if numberOnly {
			fmt.Println(len(objects), "in", s3Client.BucketName)
			return
		}

		formatAndPrintObjects(objects)

	},
}

func formatAndPrintObjects(objects []s3Types.Object) {

	if len(objects) == 0 {
		fmt.Println("No files in this bucket yet!")
		return
	}

	var filePaths []string

	for _, object := range objects {
		filePath := object.Key
		filePaths = append(filePaths, *filePath)
	}

	output := internal.PrettifyFilePaths(&filePaths)
	fmt.Println(output)
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().BoolP("numberOnly", "n", false, "Only display the number of files in the bucket")
}
