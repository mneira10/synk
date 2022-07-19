package internal

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	log "github.com/mneira10/synk/logger"
	"github.com/mneira10/synk/s3Storage"
)

const FILE_SEPARATOR = "/"

type File struct {
	Children map[string]*File
	IsFolder bool
	Name     string
	Path     string
}

type ByFilePath []string
type ByFile []File

func (strings ByFilePath) Len() int      { return len(strings) }
func (strings ByFilePath) Swap(i, j int) { strings[i], strings[j] = strings[j], strings[i] }
func (strings ByFilePath) Less(i, j int) bool {
	si := strings[i]
	sj := strings[j]

	return filePathLess(si, sj)
}

func filePathLess(si string, sj string) bool {
	spli := strings.Split(si, FILE_SEPARATOR)
	splj := strings.Split(sj, FILE_SEPARATOR)

	// are both files or folders
	if (len(spli) == 1 && len(splj) == 1) ||
		(len(spli) > 1 && len(splj) > 1) {
		return si < sj
	}

	// folders go first
	return len(spli) != 1
}

func (files ByFile) Len() int      { return len(files) }
func (files ByFile) Swap(i, j int) { files[i], files[j] = files[j], files[i] }
func (files ByFile) Less(i, j int) bool {
	si := files[i].Path
	sj := files[j].Path

	return filePathLess(si, sj)
}

func addFile(folder *File, file *File) {
	_, isFileInChildren := folder.Children[file.Name]

	if !isFileInChildren {
		folder.Children[file.Name] = file
	}
}

func prettifyFiles(folder *File, list *list.Writer) {
	var files []File
	for _, file := range folder.Children {
		files = append(files, *file)
	}

	sort.Sort(ByFile(files))

	for _, file := range files {
		// file := folder.Children[fileName]
		(*list).AppendItem(file.Name)
		if file.IsFolder {
			(*list).Indent()
			prettifyFiles(&file, list)
			(*list).UnIndent()
		}
	}

}

func GetFilePathsInLocalPath(path string) []string {
	var localFiles []string
	err := filepath.WalkDir(path,
		func(filePath string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path != filePath && !strings.HasSuffix(filePath, CONFIG_FILE_NAME) && !d.IsDir() {
				relativeFilepath, _ := filepath.Rel(path, filePath)
				localFiles = append(localFiles, relativeFilepath)
			}
			return nil
		})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return localFiles
}

func GetFilePathsInBucket(s3Client s3Storage.S3Storage) []string {
	var bucketFilePaths []string
	listObjectsData := s3Client.ListObjects()
	objects := listObjectsData.Contents
	for _, object := range objects {
		bucketFilePaths = append(bucketFilePaths, *object.Key)
	}
	return bucketFilePaths
}

func GetDiffFilePaths(localFiles *[]string, bucketFiles *[]string) []string {
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

func PrettifyFilePaths(filePaths *[]string) string {

	fileList := list.NewWriter()
	fileList.SetStyle(list.StyleConnectedBold)

	files := &File{
		Name:     "root",
		IsFolder: true,
		Children: make(map[string]*File),
	}

	for _, filePath := range *filePaths {

		split := strings.Split(filePath, FILE_SEPARATOR)
		parentFolder := files

		for i, part := range split {
			isFolder := i != len(split)-1
			currentFile := &File{
				Name:     part,
				Path:     filePath,
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
	return fileList.Render()
}
