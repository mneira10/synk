package internal

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/mneira10/synk/logger"
)

func GetFilesInLocalPath(path string) []string {
	var localFiles []string
	err := filepath.Walk(path,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// fmt.Println(path, info.Size())
			if path != filePath && !strings.HasSuffix(filePath, CONFIG_FILE_NAME) {
				localFiles = append(localFiles, filePath)
			}
			return nil
		})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return localFiles
}
