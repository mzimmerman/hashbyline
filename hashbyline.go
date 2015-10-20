// hashbyline.go
package main

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"os"
	"runtime"
	"sync"
)

func main() {
	work := make(chan int)
	text := make([]string, 0, 10000)
	lineReader := bufio.NewScanner(bufio.NewReader(os.Stdin))
	for {
		if !lineReader.Scan() {
			break
		}
		text = append(text, lineReader.Text())
	}
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	go func() {
		for x := range text {
			work <- x
		}
		close(work)
	}()
	hashes := make([]string, len(text))
	for x := 0; x < runtime.NumCPU(); x++ {
		go func() {
			for x := range work {
				sum := md5.Sum([]byte(text[x]))
				hashes[x] = hex.EncodeToString(sum[:])
			}
			wg.Done()
		}()
	}
	wg.Wait()
	csvWriter := csv.NewWriter(bufio.NewWriter(os.Stdout))
	csvWriter.Comma = '\t'
	for x := range text {
		csvWriter.Write([]string{text[x], hashes[x]})
	}
	csvWriter.Flush()
}
