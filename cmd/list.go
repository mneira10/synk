/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sort"

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
		log.Info("Listing files in bucket:")
		s3Client := s3Storage.ConfigS3()
		listObjectsData := s3Client.ListObjects()
		objects := listObjectsData.Contents

		sort.Sort(s3Storage.ByFileName(objects))

		for _, object := range objects {
			fmt.Printf("Something %v\n", *object.Key)
		}
	},
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
