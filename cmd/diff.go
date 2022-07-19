/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/mneira10/synk/internal"
	"github.com/mneira10/synk/s3Storage"

	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		config := internal.GetConfiguration(cfgFilePath)
		s3Client := s3Storage.ConfigS3(&config)

		fmt.Println("Calculating differences...")

		localFiles := internal.GetFilesInLocalPath(cfgFilePath)
		bucketFiles := getFilesInBucket(s3Client)

		diffFiles := getDiffFiles(&localFiles, &bucketFiles)

		fmt.Println("Diff:")
		for _, file := range diffFiles {
			fmt.Println(file)
		}

	},
}

func getDiffFiles(localFiles *[]string, bucketFiles *[]string) []string {
	var diffFiles []string
	sort.Strings(*bucketFiles)

	for _, localFile := range *localFiles {
		// Binary search the sorted bucket files
		i := sort.SearchStrings(*bucketFiles, localFile)
		isLocalFileInBucket := i < len(*bucketFiles) && (*bucketFiles)[i] == localFile

		if !isLocalFileInBucket {
			diffFiles = append(diffFiles, localFile)
		}
	}
	return diffFiles

}

func getFilesInBucket(s3Client s3Storage.S3Storage) []string {
	var bucketFiles []string
	listObjectsData := s3Client.ListObjects()
	objects := listObjectsData.Contents
	for _, object := range objects {
		bucketFiles = append(bucketFiles, *object.Key)
	}
	return bucketFiles
}

func init() {
	rootCmd.AddCommand(diffCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
