package main

import (
	"flag"
	"fmt"
	"io"
	"ls-part-one-4large/functions"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	//Make flags for args (future assignment)

	//Instance variables
	flag.Parse()
	files := flag.Args()
	writer := io.Writer(os.Stdout)
	useColors := functions.IsTerminal(os.Stdout)

	//No args, no flags
	if len(files) == 0 {
		//get current directory
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		entriesUnfiltered, err2 := os.ReadDir(dir)
		if err2 != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		entries := functions.DirFilter(entriesUnfiltered)

		args := make([]string, len(entries))
		for i, entry := range entries {
			if entry == nil {
				continue
			}
			args[i] = filepath.Join(dir, entry.Name())
		}

		sort.Strings(args)

		functions.SimpleLS(writer, args, useColors)
	} else {
		//sort files by files, directories before passing them in to be printed
		var fileArgs []string
		var directories []string
		for _, file := range files {
			info, err := os.Lstat(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if info.IsDir() {
				directories = append(directories, file)
			} else {
				fileArgs = append(fileArgs, file)
			}
		}

		sort.Strings(fileArgs)
		sort.Strings(directories)

		files = append(fileArgs, directories...)

		//files in the input NOTE! (.) means use current working directory
		for _, str := range files {
			fileInfo, err := os.Lstat(str)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			if fileInfo.IsDir() {
				fmt.Fprintln(writer, str+":")
				dirRoot, err2 := os.ReadDir(str)
				dirRootFilter := functions.DirFilter(dirRoot)
				dirFiles := make([]string, len(dirRootFilter))
				if err2 != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}

				for i, file := range dirRootFilter {
					dirFiles[i] = filepath.Join(str, file.Name())
				}

				sort.Strings(dirFiles)
				functions.SimpleLS(writer, dirFiles, useColors)
				fmt.Fprintln(writer, "")
			} else {
				var singleArr []string
				singleArr = append(singleArr, str)
				functions.SimpleLS(writer, singleArr, useColors)
			}
		}
	}
}
