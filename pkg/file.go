package pkg

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type chunk struct {
	bufferSize int
	offset     int64
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

const KB = 1024
const MB = 1024 * KB
const BufferSize = 1 * MB

func main() {
	srcPath := "./pkg/all_data.csv"
	destPath := "./pkg/all_data_copy.csv"
	//concurrentCopy(srcPath, destPath)
	data, _ := concurrentRead(srcPath)
	_ = ConcurrentWrite(destPath, data)
}

func ConcurrentWrite(path string, data []byte) error {
	_ = os.Remove(path)
	destFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer destFile.Close()

	filesize := len(data)

	concurrentChunksLen := int(filesize) / BufferSize
	if remainder := filesize % BufferSize; remainder != 0 {
		concurrentChunksLen++
	}

	chunks := make([]chunk, concurrentChunksLen)
	for i := 0; i < concurrentChunksLen; i++ {
		chunks[i].bufferSize = BufferSize
		chunks[i].offset = int64(BufferSize * i)
	}

	var wg sync.WaitGroup
	wg.Add(concurrentChunksLen)

	for i := 0; i < concurrentChunksLen; i++ {
		currentChunk := chunks[i]
		go func(currentChunk chunk, isLast bool) error {
			defer wg.Done()
			if isLast {
				_, err = destFile.WriteAt(data[currentChunk.offset:], currentChunk.offset)
			} else {
				_, err = destFile.WriteAt(data[currentChunk.offset:currentChunk.offset+BufferSize], currentChunk.offset)
			}
			if err != nil {
				return err
			}
			return nil
		}(currentChunk, i == concurrentChunksLen-1)
	}

	wg.Wait()
	return nil
}

func concurrentRead(path string) ([]byte, error) {
	srcFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return nil, err
	}
	filesize := info.Size()

	concurrentChunksLen := int(filesize) / BufferSize
	if remainder := filesize % BufferSize; remainder != 0 {
		concurrentChunksLen++
	}

	chunks := make([]chunk, concurrentChunksLen)
	for i := 0; i < concurrentChunksLen; i++ {
		chunks[i].bufferSize = BufferSize
		chunks[i].offset = int64(BufferSize * i)
	}

	data := make([]byte, filesize)

	var wg sync.WaitGroup
	wg.Add(concurrentChunksLen)

	for i := 0; i < concurrentChunksLen; i++ {
		currentChunk := chunks[i]
		go func(currentChunk chunk) error {
			defer wg.Done()
			buffer := make([]byte, currentChunk.bufferSize)

			bytesRead, err := srcFile.ReadAt(buffer, currentChunk.offset)
			copy(data[currentChunk.offset:int(currentChunk.offset)+bytesRead], buffer[:bytesRead])
			if err != nil && err != io.EOF {
				return err
			}
			return nil
		}(currentChunk)
	}

	wg.Wait()
	return data, nil
}

func concurrentCopy(srcPath string, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_ = os.Remove(destPath)
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	filesize := info.Size()

	concurrentChunksLen := int(filesize) / BufferSize
	if remainder := filesize % BufferSize; remainder != 0 {
		concurrentChunksLen++
	}

	chunks := make([]chunk, concurrentChunksLen)
	for i := 0; i < concurrentChunksLen; i++ {
		chunks[i].bufferSize = BufferSize
		chunks[i].offset = int64(BufferSize * i)
	}

	var wg sync.WaitGroup
	wg.Add(concurrentChunksLen)

	for i := 0; i < concurrentChunksLen; i++ {
		currentChunk := chunks[i]
		go readWriteChunk(currentChunk, &wg, file, destFile)
	}

	wg.Wait()
	return nil
}

func readWriteChunk(currentChunk chunk, wg *sync.WaitGroup, srcFile *os.File, destFile *os.File) error {
	defer wg.Done()
	buffer := make([]byte, currentChunk.bufferSize)

	bytesRead, err := srcFile.ReadAt(buffer, currentChunk.offset)
	if err != nil && err != io.EOF {
		return err
	}

	_, err = destFile.WriteAt(buffer[:bytesRead], currentChunk.offset)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
