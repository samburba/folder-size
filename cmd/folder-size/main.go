package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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
			size = GetFolderSize(directory + "/" + f.Name(), size)
		} else {
			return size
		}
	}
	return size
}

func main() {
	directoryPtr := flag.String("dir", ".", "the directory to view")
	flag.Parse()
	fmt.Println(*directoryPtr)
	files, err := ioutil.ReadDir(*directoryPtr)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			size := GetFolderSize(*directoryPtr + "/" + f.Name(), 0)
			fmt.Println(f.Name(), ByteCountSI(size))
		} else {
			fmt.Println(f.Name(), ByteCountSI(f.Size()))
		}
	}
}