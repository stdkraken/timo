package main

import (
	"io"
	"log"
	"os"
)

func ReadStdin() []byte {
	buffer := make([]byte, BUFFER_SIZE)
	// read from os.Stdin, until err != nil or n != 0 or
	n, err := os.Stdin.Read(buffer)
	if err != nil {
		log.Fatalln(err.Error())
	}
	buffer = BufferCut(buffer, n-1) // remove last byte, because this is 0xa, and not very helpful
	return buffer
}

func ReadAllMaxSize(reader io.Reader, maxSize int, bufferSize int) (data []byte, err error) {
	err = nil
	data = []byte{}
	n := 0
	for err == nil {
		buffer := make([]byte, bufferSize)
		n, err = reader.Read(buffer)
		if len(data)+n <= maxSize {
			data = append(data, BufferCut(buffer, n)...)
		} else {
			// don't read any further => the maxSize is already reached
		}
	}
	return data, err
}
