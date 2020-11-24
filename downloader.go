// Copyright (c) 2020 Kristian Rumberg (kristianrumberg@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"sync"
)

// Downloader interface provides a source for audio books
type Downloader interface {
	// Download an audio book
	DownloadBook(podbooks *podBookStorage, uri string) error
}

// BookResult holds the download result for the book given by uri
type BookResult struct {
	uri string
	err error
}

// DownloadBooks downloads multiple books at the same time. Expected to run as a goroutine
func DownloadBooks(downloader Downloader, podbooks *podBookStorage, uris []string, outchan chan BookResult) {
	var wg sync.WaitGroup

	for _, uri := range uris {
		wg.Add(1)

		go func(uri string) {
			defer wg.Done()

			err := downloader.DownloadBook(podbooks, uri)
			outchan <- BookResult{uri, err}
		}(uri)
	}

	wg.Wait()

	close(outchan)
}
