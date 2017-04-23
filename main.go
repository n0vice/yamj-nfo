package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	videoExtensions = []string{".mkv", ".avi"}
)

const (
	nfoExtension = ".nfo"
)

func main() {
	// init cli args
	kingpin.Parse()

	// read files
	fileList, err := ioutil.ReadDir(*directory)
	if err != nil {
		println("failed to read directory. ", err.Error())
		return
	}

	// move to map, ignore directories
	fileMap := make(map[string]os.FileInfo)
	for _, file := range fileList {
		if file.IsDir() {
			continue
		}
		fileMap[file.Name()] = file
	}
	regex, err := regexp.Compile(`(\d{4})`)
	if err != nil {
		println("failed to compile regular expression. ", err.Error())
		return
	}
	moviesWithoutNFO := make([]string, 0, len(fileMap))
	// get those without nfo file
	index := 0
	for name := range fileMap {
		fileNoExt, ok := isVideo(name)
		if !ok {
			continue
		}
		// try get nfo
		_, ok = fileMap[fileNoExt+nfoExtension]
		if !ok {
			// no nfo exist
			moviesWithoutNFO = append(moviesWithoutNFO, fileNoExt)
			year := regex.Find([]byte(fileNoExt))
			if year == nil {
				println("failed to match year from movie filename. ", fileNoExt)
				continue
			}
			indexOfYear := strings.Index(fileNoExt, string(year))
			fileNameToSearch := strings.Replace(fileNoExt[:indexOfYear-1]+" "+string(year), ".", " ", -1)
			if index != 0 {
				// sleep for a while before each request so we don't get banned by Google
				time.Sleep(time.Duration(*timeSleep) * time.Second)
			}
			println("Search for ", fileNameToSearch)
			response, err := searchMovie(fileNameToSearch)
			index++
			if err != nil {
				println(fmt.Errorf("failed to search a movie. %v", err.Error()))
				continue
			}
			if response == nil {
				println("search response is nil. ", fileNameToSearch)
				continue
			}
			if len(response.Items) > 0 {
				topItemLink := response.Items[0].Link
				println(" ", topItemLink)
				//create nfo file with that link
				err := ioutil.WriteFile(path.Join(*directory, fileNoExt+nfoExtension), []byte(topItemLink), 0644)
				if err != nil {
					fmt.Printf("failed to create nfo file %v. %v", fileNoExt+nfoExtension, err.Error())
					continue
				}
			}
		}
		index++
	}
}

func isVideo(name string) (string, bool) {
	for _, videoExt := range videoExtensions {
		fileExt := filepath.Ext(name)
		if fileExt == videoExt {
			return strings.TrimSuffix(name, fileExt), true
		}
	}
	return "", false
}
