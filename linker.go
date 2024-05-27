package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	NONEFILE = iota
	ASSEMBLY
	APLUS
)

const (
	NOLINK = iota
	STATIC
	DYNAMIC
)

type link struct {
	name                 string
	isfile               bool
	headername           string
	codename             string
	filetype             int
	linktype             int
	assemblyrequirements []string
}

type links struct {
	alllinks []link
	count    int
}

func (ls *links) AddFile(name string, filename string, ftype string) {
	for index, val := range ls.alllinks {
		if val.name == name {
			if filename[len(filename)-3:] == "asm" || filename[len(filename)-2:] == "ap" {
				ls.alllinks[index].codename = filename
			} else if filename[len(filename)-3:] == "aph" {
				ls.alllinks[index].headername = filename
			}
			if ftype == "file" {
				ls.alllinks[index].isfile = true
			} else {
				ls.alllinks[index].isfile = false
			}
			return
		}
	}
	var newlink link
	if filename[len(filename)-3:] == "asm" || filename[len(filename)-2:] == "ap" {
		newlink.name = name
		newlink.codename = filename
	} else if filename[len(filename)-3:] == "aph" {
		newlink.name = name
		newlink.headername = filename
	}
	if ftype == "file" {
		newlink.isfile = true
	} else {
		newlink.isfile = false
	}
	ls.alllinks = append(ls.alllinks, newlink)
}

func (ls *links) AddType(name string, filetype string, linktype string) {
	for index, val := range ls.alllinks {
		if val.name == name {
			if filetype == "assembly" {
				ls.alllinks[index].filetype = ASSEMBLY
			} else if filetype == "aplus" {
				ls.alllinks[index].filetype = APLUS
			}
			if linktype == "static" {
				ls.alllinks[index].linktype = STATIC
			} else if filetype == "dynamic" {
				ls.alllinks[index].linktype = DYNAMIC
			}
			return
		}
	}
	var newlink link
	if filetype == "assembly" {
		newlink.filetype = ASSEMBLY
	} else if filetype == "aplus" {
		newlink.filetype = APLUS
	}
	if linktype == "static" {
		newlink.linktype = STATIC
	} else if filetype == "dynamic" {
		newlink.linktype = DYNAMIC
	}
	ls.alllinks = append(ls.alllinks, newlink)
}

func (ls *links) AddASMRequiements(name string, requirements []string) {
	for index, val := range ls.alllinks {
		if val.name == name {
			ls.alllinks[index].assemblyrequirements = append(ls.alllinks[index].assemblyrequirements, requirements...)
			return
		}
	}
	var newlink link
	newlink.assemblyrequirements = append(newlink.assemblyrequirements, requirements...)
	ls.alllinks = append(ls.alllinks, newlink)
}

func ParseAPMake(folder string) links {
	file, _ := os.ReadFile(folder + "/build.apmake")
	regexline := regexp.MustCompile("\r")
	file = regexline.ReplaceAll(file, nil)
	regexline = regexp.MustCompile(" ")
	file = regexline.ReplaceAll(file, nil)
	splitfile := strings.Split(string(file), "\n")
	index := 0
	var returnlinks links
	if splitfile[index] == "dir:" {
		for {
			index++
			if splitfile[index] == "buildtype:" {
				break
			}
			line := strings.Split(splitfile[index], ":")
			returnlinks.AddFile(line[0], line[1], line[2])
		}
		for {
			index++
			if splitfile[index] == "externalassemblyrequirements:" {
				break
			}
			line := strings.Split(splitfile[index], ":")
			returnlinks.AddType(line[0], line[1], line[2])
		}
		for {
			index++
			if index != len(splitfile) {
				break
			}
			line := strings.Split(splitfile[index], ":")
			returnlinks.AddASMRequiements(line[0], strings.Split(line[1], ","))
		}
	}
	return returnlinks
}

func Link(requirements []string, filename string, bl *buildlog) {
	imports := ""
	index := 12345
	indexmap := make(map[string]int, 0)
	for _, currfilename := range requirements {
		currfilename = "imports/" + currfilename
		bl.AddLog("Parsing \""+"imports/"+currfilename+"\"'s make file", 0)
		allimports := ParseAPMake(currfilename)
		bl.AddLog("Make file parsed", 0)
		for _, curlink := range allimports.alllinks {
			bl.AddLog("Modifying \""+currfilename+"/"+curlink.codename+"\" for link", 0)
			file, _ := os.ReadFile(currfilename + "/" + curlink.codename)

			regexline := regexp.MustCompile("\\$\\w+[^:\\s]*")
			wordocc := regexline.FindAll(file, -1)
			for _, word := range wordocc {
				_, isin := indexmap[string(word)]
				if !isin {
					indexmap[string(word)] = index
					regexline = regexp.MustCompile("\\" + string(word))
					file = regexline.ReplaceAll(file, []byte("_func_"+strconv.Itoa(index)))
					index++
				}
			}
			imports += string(file)
			bl.AddLog("\""+currfilename+"/"+curlink.codename+"\" modified", 0)
		}
	}
	file, _ := os.ReadFile("build/" + filename)
	regexline := regexp.MustCompile("{ASSEMBLYLINKHERE}")
	file = regexline.ReplaceAll(file, []byte(imports))
	f, _ := os.Create("build/" + filename)
	f.WriteString(string(file))
}
