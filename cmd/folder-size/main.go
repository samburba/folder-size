package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func GetFolderSize(directory string, size int64) int64 {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		size = f.Size() + size
		if f.IsDir() {
			size = GetFolderSize(directory+"/"+f.Name(), size)
		} else {
			return size
		}
	}
	return size
}

func main() {
	directoryPtr := flag.String("dir", ".", "the directory to view")
	visualizationPtr := flag.Bool("v", false, "visualization")
	flag.Parse()
	fmt.Println(*directoryPtr)
	files, err := ioutil.ReadDir(*directoryPtr)
	if err != nil {
		log.Fatal(err)
	}

	folders := make(map[string]int64)

	for _, f := range files {
		if f.IsDir() {
			size := GetFolderSize(*directoryPtr+"/"+f.Name(), 0)
			folders[f.Name()] = size
		} else {
			folders[f.Name()] = f.Size()
		}
	}

	names := make([]string, 0, len(folders))
	for k := range folders {
		names = append(names, k)
	}

	sort.Slice(names, func(i, j int) bool {
		return folders[names[i]] > folders[names[j]]
	})
	top := folders[names[0]]

	spacing := 30
	if *visualizationPtr {
		for _, name := range names {
			if folders[name] == 0 {
				fmt.Println(strings.Repeat(" ", spacing+2), name, ByteCountSI(folders[name]))
			} else {
				repeat := (spacing * int(folders[name]) / int(top)) + 1
				fmt.Println(strings.Repeat("#", repeat), strings.Repeat(" ", spacing+1-repeat), name, ByteCountSI(folders[name]))
			}
		}
	} else {
		for _, name := range names {
			fmt.Println(name, ByteCountSI(folders[name]))
		}
	}
}
