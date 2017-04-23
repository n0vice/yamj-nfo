package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	googleAPIkey = kingpin.Flag("apikey", "Google API Key").Required().Short('a').String()
	googleAPICX  = kingpin.Flag("cx", "Google Custom Search CX ID").Required().Short('c').String()
	proxy        = kingpin.Flag("proxy", "Proxy server URL in form of http://address:port").Short('p').String()
	directory    = kingpin.Flag("dir", "Directory to parse").Default(getWd()).Short('d').String()
	timeSleep    = kingpin.Flag("sleep", "Time to sleep before each HTTP request").Default("5").Short('s').Int()
)

func getWd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get working directory. %v", err.Error()))
	}
	return dir
}
