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

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// Sha256String computes the SHA256 of string
func Sha256String(s string) string {
	sum := sha256.Sum256([]byte(s))

	return hex.EncodeToString(sum[:])
}

// IsFile checks if file exists
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !s.IsDir()
}

// IsSymlink checks if path is a symlink
func IsSymlink(path string) bool {
	linkInfo, err := os.Lstat(path)
	if err == nil {
		if linkInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			return true
		}
	}

	return false
}

// FindSingleMatchInDir checksif there is a single file matching pattern in the directory or returns an error
func FindSingleMatchInDir(dir string, pattern string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.mp3"))
	if err != nil {
		return "", fmt.Errorf("failed to list any mp3 file in %s", dir)
	}

	if len(files) != 1 {
		return "", fmt.Errorf("invalid number of mp3 files, found %s", files)
	}

	return files[0], nil
}
