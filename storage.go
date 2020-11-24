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
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/krumberg/podbook/common"
)

type Config struct {
	URL string
}

type podBookStorage struct {
	temp   string
	books  string
	db     string
	config Config
}

func podbookDir(workdir string) string {
	return filepath.Join(workdir, ".podbook")
}

func podbookConfigPath(workdir string) string {
	return filepath.Join(podbookDir(workdir), "config.txt")
}

func makePodBookStorage(workdir string, config *Config) *podBookStorage {
	dotPodbook := podbookDir(workdir)

	podbook := podBookStorage{
		books: filepath.Join(workdir, "books"),
		temp:  filepath.Join(dotPodbook, "temp"),
		db:    filepath.Join(dotPodbook, "db"),
	}

	if config != nil {
		podbook.config = *config
	}

	return &podbook
}

func initPodBookStorage(workdir string, url string) (*podBookStorage, error) {
	podbook := makePodBookStorage(workdir, &Config{URL: url})

	dirs := []string{podbook.temp, podbook.books, podbook.db}
	for _, d := range dirs {
		err := os.MkdirAll(d, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	configPath := podbookConfigPath(workdir)

	f, err := os.Create(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create config file %w", err)
	}
	defer f.Close()

	err = toml.NewEncoder(f).Encode(podbook.config)
	if err != nil {
		return nil, fmt.Errorf("failed to write config %w", err)
	}

	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to write file %w", err)
	}

	return podbook, nil
}

func openPodBookStorage(workdir string) (*podBookStorage, error) {
	podbook := makePodBookStorage(workdir, nil)

	_, err := toml.DecodeFile(podbookConfigPath(workdir), &podbook.config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode podbook config")
	}

	return podbook, nil
}

func (podbook *podBookStorage) makeDBLink(linkPath, dbFilename string) error {
	linkPathDir := filepath.Dir(linkPath)

	linkContentDir, err := filepath.Rel(linkPathDir, podbook.db)
	if err != nil {
		return fmt.Errorf("failed to resolve symlink path for %s and %s", linkPathDir, podbook.db)
	}

	linkContentFile := filepath.Join(linkContentDir, dbFilename)

	err = os.Symlink(linkContentFile, linkPath)
	if err != nil {
		return fmt.Errorf("failed to symlink %s to %s", linkContentFile, linkPath)
	}

	return nil
}

func (podbook *podBookStorage) fixBrokenLinks() error {
	return filepath.Walk(podbook.books,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !common.IsSymlink(path) {
				return nil
			}

			target, err := os.Readlink(path)
			if err != nil {
				return err
			}

			// Prepend directory to symlink
			target = filepath.Join(filepath.Dir(path), target)

			if !common.IsFile(target) {
				fmt.Printf("Link for %s is broken\n", path)

				os.Remove(path)

				return podbook.makeDBLink(path, filepath.Base(target))
			}

			return nil
		})
}

func (podbook *podBookStorage) dbPathFromURL(url string) string {
	urlChecksumString := common.Sha256String(url)

	return path.Join(podbook.db, urlChecksumString+".mp3")
}

func (podbook *podBookStorage) hasURL(url string) bool {
	return common.IsFile(podbook.dbPathFromURL(url))
}

func (podbook *podBookStorage) importTempFile(filePath, url string) error {
	dbPath := podbook.dbPathFromURL(url)

	err := os.Rename(filePath, dbPath)
	if err != nil {
		return fmt.Errorf("failed to rename url %s", url)
	}

	linkPath := filepath.Join(podbook.books, filepath.Base(filePath))

	return podbook.makeDBLink(linkPath, filepath.Base(dbPath))
}
