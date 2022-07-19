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

	files := &File{Name: "root",
		IsFolder: true, Children: make(map[string]*File)}

	// var fileNameStack []string

	for _, object := range *objects {
		fileName := object.Key

		split := strings.Split(*fileName, "/")

		parentFolder := files

		for i, part := range split {
			isFolder := i != len(split)-1
			currentFile := &File{
				Name:     part,
				IsFolder: isFolder,
				Children: make(map[string]*File),
			}

			addFile(parentFolder, currentFile)

			if isFolder {
				parentFolder = parentFolder.Children[currentFile.Name]
			}

		}

	}

	prettifyFiles(files, &fileList)
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

type File struct {
	Children map[string]*File
	IsFolder bool
	Name     string
}

func addFile(folder *File, file *File) {

	_, isFileInChildren := folder.Children[file.Name]

	if !isFileInChildren {
		folder.Children[file.Name] = file
	}
}

func prettifyFiles(folder *File, list *list.Writer) {
	var keys []string
	for k, _ := range folder.Children {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, fileName := range keys {
		file := folder.Children[fileName]
		(*list).AppendItem(fileName)
		if file.IsFolder {
			(*list).Indent()
			prettifyFiles(file, list)
			(*list).UnIndent()
		}
	}

}
