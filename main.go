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
	"flag"
	"fmt"
	"os"
)

func dieErr(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}

func die(msg string) {
	dieErr(fmt.Errorf("%s", msg))
}

func showHelpAndExit() {
	die("expected 'init', 'fix', 'get' or 'rss' subcommands")
}

func main() {
	if len(os.Args) < 2 {
		showHelpAndExit()
	}

	subCmds := map[string]func(string){
		"init": runInit,
		"get":  runGet,
		"fix":  runFix,
		"rss":  runRss,
	}

	subCmd := os.Args[1]

	if fun, ok := subCmds[subCmd]; ok {
		fun(subCmd)
	} else {
		showHelpAndExit()
	}
}

func runInit(subcmd string) {
	initCmd := flag.NewFlagSet(subcmd, flag.ExitOnError)

	err := initCmd.Parse(os.Args[2:])
	if err != nil {
		die("failed to parse init args")
	}

	if len(initCmd.Args()) != 1 {
		die("Failed to parse url")
	}

	baseurl := initCmd.Args()[0]

	_, err = initPodBookStorage(".", baseurl)
	if err != nil {
		dieErr(err)
	}
}

func runGet(subcmd string) {
	podbook, err := openPodBookStorage(".")
	if err != nil {
		die("failed to open podbook storage. Did you execute podbook init <url>?")
	}

	getCmd := flag.NewFlagSet(subcmd, flag.ExitOnError)
	rssFile := getCmd.String("rss", "", "rss")

	err = getCmd.Parse(os.Args[2:])
	if err != nil {
		die("failed to parse get args")
	}

	fmt.Println(*rssFile)

	urls := getCmd.Args()

	if len(urls) == 0 {
		die("no inputs supplied to get command")
	}

	downloader := NewYoutubeDownloader(podbook)

	downloadResults := make(chan BookResult)

	go DownloadBooks(downloader, podbook, urls, downloadResults)

	for result := range downloadResults {
		fmt.Printf("Finished %v with result %v\n", result.uri, result.err)

		if result.err == nil && len(*rssFile) > 0 {
			fmt.Println("Writing ", *rssFile)

			err := writeRss(podbook, *rssFile)
			if err != nil {
				dieErr(err)
			}
		}
	}
}

func runFix(subcmd string) {
	podbook, err := openPodBookStorage(".")
	if err != nil {
		die("failed to open podbook storage. Did you execute podbook init <url>?")
	}

	fixCmd := flag.NewFlagSet(subcmd, flag.ExitOnError)

	err = fixCmd.Parse(os.Args[2:])
	if err != nil {
		die("failed to parse fix args")
	}

	if len(fixCmd.Args()) > 0 {
		die("error: Fix command takes no arguments")
	}

	err = podbook.fixBrokenLinks()
	if err != nil {
		dieErr(fmt.Errorf("failed to fix symlinks: %w", err))
	}
}

func runRss(subcmd string) {
	podbook, err := openPodBookStorage(".")
	if err != nil {
		die("failed to open podbook storage. Did you execute podbook init <url>?")
	}

	rssCmd := flag.NewFlagSet(subcmd, flag.ExitOnError)

	err = rssCmd.Parse(os.Args[2:])
	if err != nil {
		die("failed to parse rss args")
	}

	if len(rssCmd.Args()) != 1 {
		die("no output rss supplied")
	}

	outfile := rssCmd.Args()[0]

	err = writeRss(podbook, outfile)
	if err != nil {
		dieErr(err)
	}
}
