package s3Storage

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const FILE_SEPARATOR = "/"

type ByFileName []types.Object

func (objects ByFileName) Len() int      { return len(objects) }
func (objects ByFileName) Swap(i, j int) { objects[i], objects[j] = objects[j], objects[i] }
func (objects ByFileName) Less(i, j int) bool {
	si := *objects[i].Key
	sj := *objects[j].Key

	return fileStringLess(si, sj)
}

func fileStringLess(si string, sj string) bool {
	spli := strings.Split(si, FILE_SEPARATOR)
	splj := strings.Split(sj, FILE_SEPARATOR)

	// are both files or folders
	if (len(spli) == 1 && len(splj) == 1) || (len(spli) > 1 && len(splj) > 1) {
		return si < sj
	}

	// folders go first
	return len(spli) != 1

}
