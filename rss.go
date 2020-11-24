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
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/krumberg/podbook/common"
)

func writeRss(podbook *podBookStorage, outpath string) error {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Audiobooks",
		Link:        &feeds.Link{Href: podbook.config.URL},
		Description: "A list of Audiobooks",
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	err := filepath.Walk(podbook.books,
		func(path string, info os.FileInfo, err error) error {
			if common.IsFile(path) {
				title := filepath.Base(path)
				title = strings.TrimSuffix(title, filepath.Ext(title))

				uri := podbook.config.URL
				for _, p := range strings.Split(path, "/") {
					uri += "/" + url.PathEscape(p)
				}

				st, err := os.Stat(path)
				if err != nil {
					return fmt.Errorf("failed to stat path %s", path)
				}

				item := feeds.Item{
					Title:     title,
					Link:      &feeds.Link{Href: uri},
					Enclosure: &feeds.Enclosure{Url: uri, Length: strconv.FormatInt(st.Size(), 10), Type: "audio/mpeg"},
					Created:   now,
				}

				feed.Items = append(feed.Items, &item)
			}

			return nil
		})
	if err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outpath, []byte(rss), 0o600)
}
