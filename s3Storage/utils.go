package s3Storage

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

// basically all copied from:
// https://github.com/gabrielhora/s3md5
const CHUNK_SIZE int = 15

// TODO: normie md5 files
func GetEtagFromFile(fileName *string) string {
	// This will return the ETag for multipart files in S3. I think
	// that means for files >5GB. Have to figure out how to calculate the
	// regular md5 hash for normie files.

	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	reader, err := os.Open(*fileName)
	defer reader.Close()
	if err != nil {
		panic(err.Error())
	}

	chunkSizeInBytes := 1024 * 1024 * CHUNK_SIZE
	buffer := make([]byte, chunkSizeInBytes)
	hasher := md5.New()

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(buffer, chunkSizeInBytes)
	scanner.Split(splitByBufferSize)

	totalChunks := 0
	var md5bytes []byte

	for scanner.Scan() {
		_, err := hasher.Write(scanner.Bytes())
		if err != nil {
			panic(err.Error())
		}
		md5bytes = append(md5bytes, hasher.Sum(nil)...)
		totalChunks++
		hasher.Reset()
	}

	if scanner.Err() != nil {
		panic(scanner.Err().Error())
	}

	hasher.Write(md5bytes)
	filemd5 := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("%s-%d\n", filemd5, totalChunks)
}

func splitByBufferSize(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return len(data), data[0:], nil
}
