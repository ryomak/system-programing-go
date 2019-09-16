package main

import (
	"bytes"
  "hash/crc32"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func textChunk(text string) io.Reader {
	byteData := []byte(text)
	var buffer  bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, int32(len(byteData)))
	buffer.WriteString("tEXt")
	buffer.Write(byteData)
	crc := crc32.NewIEEE()
	io.WriteString(crc, "tEXt")
	binary.Write(&buffer, binary.BigEndian, crc.Sum32())
	return &buffer
}

func dumpChunk(chunk io.Reader) {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length)
	buffer := make([]byte, 4)
	chunk.Read(buffer)
	if string(buffer) ==  "tEXt" {
		text := make([]byte,length)
		chunk.Read(text)
		fmt.Printf("chunk %s (%d bytes) :%s \n", string(buffer), length,string(text))
	}else {
		fmt.Printf("chunk %s (%d bytes)\n", string(buffer), length)
	}
}

func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader
	var offset int64 = 8

	file.Seek(8, 0)

	for {
		var length int32
		if err := binary.Read(file, binary.BigEndian, &length); err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))
		offset, _ = file.Seek(int64(length+8), 1)
	}
	return chunks
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newFile, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	chunks := readChunks(file)
	io.WriteString(newFile, "\x89PNG\r\n\x1a\n")
	io.Copy(newFile, chunks[0])
	io.Copy(newFile, textChunk(os.Args[3]))
	for _, v := range chunks[1:] {
		io.Copy(newFile, v)
	}
	newChunks := readChunks(newFile)
	for _, v := range newChunks {
		dumpChunk(v)
	}
}
