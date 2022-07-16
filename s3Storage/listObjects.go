package s3Storage

import (
	"sort"
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

	// both are files in the same dir
	if len(spli) == 1 && len(splj) == 1 {
		stringSlice := sort.StringSlice{si, sj}
		return stringSlice.Less(0, 1)
	}

	// files in the same dir with dir names
	if len(spli) == len(splj) {
		return fileStringLess(si[len(spli[0])+1:], sj[len(splj[0])+1:])
	}

	// different dirs or file and dir
	if len(spli) == 1 {
		return false
	} else if len(splj) == 1 {
		return true
	}

	// different folders
	stringSlice := sort.StringSlice{si, sj}
	return stringSlice.Less(0, 1)

}
