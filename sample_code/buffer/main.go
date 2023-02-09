package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	err := loadAndProcess(os.Args[1], func(data []byte) {
		fmt.Print(string(data))
	})
	if err != nil {
		fmt.Println("error:", err)
	}
}

func loadAndProcess(fileName string, process func([]byte)) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, 100)
	for {
		count, err := file.Read(data)
		process(data[:count])
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if count == 0 {
			return nil
		}
	}
}
