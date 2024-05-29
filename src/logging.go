package main

import (
	"strings"
	"time"
)

type loggingconversion struct {
	originalcodenum int
	originalcode    code
	assemblycode    code
}

type textlog struct {
	text      string
	time      string
	iswarning bool
	iserr     bool
}

type buildlog struct {
	logs  []textlog
	count int
}

func SetLoggingConversion(apcode string, lc []*loggingconversion) {
	splitcode := strings.Split(apcode, "\n")
	for index, val := range lc {
		lc[index].originalcode.AddStringCode(splitcode[val.originalcodenum])
	}
}

func (bl *buildlog) AddLog(text string, flag int) {
	var log textlog
	log.text = text
	if flag == 1 {
		log.iswarning = true
	} else if flag == 2 {
		log.iserr = true
	}
	dt := time.Now()
	log.time = dt.Format("2006 01 2 15:04:05.000000")
	bl.logs = append(bl.logs, log)
	bl.count += 1
}

func (bl *buildlog) ReturnLogs() string {
	var returnstring string
	for index, val := range bl.logs {
		if index != 0 {
			returnstring += "\n"
		}
		returnstring += val.time + "\n"
		if val.iswarning {
			returnstring += "WARNING:  "
		} else if val.iserr {
			returnstring += "ERROR: "
		}
		returnstring += val.text + "\n"
	}
	returnstring = returnstring[:len(returnstring)-1]
	return returnstring
}
