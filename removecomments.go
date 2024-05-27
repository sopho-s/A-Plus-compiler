package main

import (
	"os"
	"regexp"
	"strings"
)

func RemoveComments(infile string) {
	file, _ := os.ReadFile(infile)
	regexfile := regexp.MustCompile("//.*")
	newfile := regexfile.ReplaceAll(file, nil)
	regexfile = regexp.MustCompile(`\/\*[^*]*\*+(?:[^/*][^*]*\*+)*\/`)
	newfile = regexfile.ReplaceAll(newfile, nil)
	filesplit := strings.Split(string(newfile), "\n")
	returnfile := ""
	for _, val := range filesplit {
		if len(val) > 0 {
			returnfile += val + "\n"
		}
	}
	returnfile = returnfile[:len(returnfile)-1]
	f, _ := os.Create("build/build.ap")
	f.Write([]byte(returnfile))
}
