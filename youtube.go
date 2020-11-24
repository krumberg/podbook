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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/krumberg/podbook/common"
)

type youtubeDownloader struct {
	podbooks *podBookStorage
}

// NewYoutubeDownloader creates a youtube audiobook downloader
func NewYoutubeDownloader(podbooks *podBookStorage) Downloader {
	return &youtubeDownloader{podbooks: podbooks}
}

func (downloader *youtubeDownloader) downloadBookToDirectory(url, dir string) (string, error) {
	c := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", url)
	c.Dir = dir

	_, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("failed to download %s, cleaning up %s", url, dir)
	}

	return common.FindSingleMatchInDir(dir, "*.mp3")
}

func (downloader *youtubeDownloader) DownloadBook(podbooks *podBookStorage, url string) error {
	if podbooks.hasURL(url) {
		return fmt.Errorf("book %s already exists: %s", url, podbooks.dbPathFromURL(url))
	}

	fmt.Println("Downloading", url)

	// Create tempdir
	tempdir, err := ioutil.TempDir(podbooks.temp, "book")
	if err != nil {
		return fmt.Errorf("failed to create tempdir for %s", url)
	}

	defer os.RemoveAll(tempdir)

	filename, err := downloader.downloadBookToDirectory(url, tempdir)
	if err != nil {
		return fmt.Errorf("failed to download book: %w", err)
	}

	return podbooks.importTempFile(filename, url)
}
