/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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

		localFiles := internal.GetFilePathsInLocalPath(cfgFilePath)
		bucketFiles := internal.GetFilePathsInBucket(s3Client)

		diffFiles := internal.GetDiffFilePaths(&localFiles, &bucketFiles)

		output := internal.PrettifyFilePaths(&diffFiles)
		fmt.Println(output)

	},
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
