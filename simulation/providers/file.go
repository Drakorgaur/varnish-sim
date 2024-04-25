//  Copyright 2024 Mark Barzali
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0

package providers

import (
	"bufio"
	"os"
)

const fileProviderName = "file"

type FileProvider struct {
	Files     []string
	Formatter func(string) (string, int)
}

func (f *FileProvider) SetFormatter(frmt func(string) (string, int)) {
	if frmt == nil {
		f.Formatter = defaultFormatter
	}
	f.Formatter = frmt
}

func (f *FileProvider) String() string {
	return fileProviderName
}

func (f *FileProvider) Channel() <-chan *Request {
	if f.Formatter == nil {
		f.Formatter = defaultFormatter
	}

	ch := make(chan *Request)

	go func() {
		defer close(ch)

		for _, file := range f.Files {
			if openAndProvide(file, f.Formatter, ch) {
				// bad file
				continue
			}
		}
		// send nil to indicate end of data
		ch <- nil
	}()

	return ch
}

func openAndProvide(file string, frmt func(string) (string, int), ch chan *Request) bool {
	// open file
	fd, err := os.Open(file)
	if err != nil {
		// bad file
		return true
	}
	defer fd.Close()

	pipeReaderChannel(bufio.NewReader(fd), frmt, ch)
	return false
}

func pipeReaderChannel(r *bufio.Reader, frmt func(string) (string, int), ch chan *Request) {
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		url, size := frmt(line)
		ch <- &Request{url, size}
	}
}
