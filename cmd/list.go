/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jedib0t/go-pretty/v6/list"
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

		printBucketFiles(s3Client)

	},
}

func printBucketFiles(s3Client s3Storage.S3Storage) {
	listObjectsData := s3Client.ListObjects()
	objects := listObjectsData.Contents

	// for _, object := range objects {
	// 	fmt.Printf("%v\n", *object.Key)
	// }

	formatAndPrintObjects(&objects)

}

func formatAndPrintObjects(objects *[]types.Object) {
	sort.Sort(s3Storage.ByFileName(*objects))

	if len(*objects) == 0 {
		fmt.Println("No files in this bucket yet!")
		return
	}

	fileList := list.NewWriter()
	fileList.SetStyle(list.StyleConnectedBold)

	var fileNameStack []string

	for _, object := range *objects {
		fileName := object.Key

		split := strings.Split(*fileName, "/")

		for i, part := range split {

			if i != len(split)-1 {

				if len(fileNameStack) <= i {
					fileList.AppendItem(part)
					fileList.Indent()

					fileNameStack = append(fileNameStack, part)
				} else if fileNameStack[i] != part {

					fileNameStack = split[:i]
					for j := 0; j < len(split)-i; j++ {
						fileList.UnIndent()
					}

					fileList.AppendItem(part)
					fileList.Indent()
					fileNameStack = append(fileNameStack, part)
				}

				continue
			}

			if len(split) == 1 {
				for j := 0; j < len(fileNameStack); j++ {
					fileList.UnIndent()
				}
				fileNameStack = []string{}
			}

			fileList.AppendItem(part)
		}
	}

	fileList.Length()
	fmt.Println(fileList.Render())
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
