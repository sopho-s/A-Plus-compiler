package main

import (
	"strconv"
	"strings"
)

func AddingPrecodeToFunctions(functions code, data code) code {
	var outcode code
	outcode.AddStringCode("global _main")
	outcode.AddStringCode("section .data")
	outcode.AddCode(data)
	outcode.AddStringCode("section .text")
	outcode.AddStringCode("{ASSEMBLYLINKHERE}")
	outcode.AddCode(functions)
	return outcode
}

func ConvertToNASM(intcode string, funcname string, offsetmap *map[string]int, floatcountmap *map[string]int) (code, code, []*loggingconversion) {
	splitcode := strings.Split(intcode, "\n")
	var precode code
	var outcode code
	var data code
	var log []*loggingconversion
	if funcname == "main" {
		funcname = "_main"
	}
	outcode.AddStringCode(funcname + ":")
	startindex, _ := strconv.ParseInt(strings.Split(splitcode[0], " ")[0], 10, 32)
	for _, value := range splitcode {
		linesplit := strings.Split(value, " ")
		index, _ := strconv.ParseInt(linesplit[0], 10, 32)
		if len(log) <= int(index)-1 {
			var singlelog loggingconversion
			log = append(log, &singlelog)
		}
		switch linesplit[1] {
		case "PUSH":
			if string(linesplit[2][0]) == "I" {
				outcode.AddStringCode("PUSH " + registers[linesplit[2]])
				log[index-startindex].assemblycode.AddStringCode("PUSH " + registers[linesplit[2]])
			} else if string(linesplit[2][0]) == "F" {
				outcode.AddStringCode("SUB ESP, 4")
				log[index-startindex].assemblycode.AddStringCode("SUB ESP, 4")
				outcode.AddStringCode("MOVSS DWORD [ESP], " + registers[linesplit[2]])
				log[index-startindex].assemblycode.AddStringCode("MOVSS DWORD [ESP], " + registers[linesplit[2]])
			}
			break
		case "POP":
			if string(linesplit[2][0]) == "I" {
				outcode.AddStringCode("POP " + registers[linesplit[2]])
				log[index-startindex].assemblycode.AddStringCode("POP " + registers[linesplit[2]])
			} else if string(linesplit[2][0]) == "F" {
				outcode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", DWORD [ESP]")
				log[index-startindex].assemblycode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", DWORD [ESP]")
				outcode.AddStringCode("ADD ESP, 4")
				log[index-startindex].assemblycode.AddStringCode("ADD ESP, 4")
			}
			break
		case "IADD":
			outcode.AddStringCode("ADD " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("ADD " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "ISUB":
			_, isin := registers[linesplit[3]]
			if isin {
				outcode.AddStringCode("SUB " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			} else {
				outcode.AddStringCode("SUB " + registers[linesplit[2]] + ", " + linesplit[3])
			}
			log[index-startindex].assemblycode.AddStringCode("SUB " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "IMUL":
			outcode.AddStringCode("IMUL " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("IMUL " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "IDIV":
			outcode.AddStringCode("XOR EDX, EDX")
			log[index-startindex].assemblycode.AddStringCode("XOR EDX, EDX")
			outcode.AddStringCode("IDIV " + registers[linesplit[2]])
			log[index-startindex].assemblycode.AddStringCode("IDIV " + registers[linesplit[2]])
			break
		case "IMOV":
			_, isreg := registers[linesplit[3]]
			if isreg {
				outcode.AddStringCode("MOV " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
				log[index-1-startindex].assemblycode.AddStringCode("MOV " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			} else {
				if len(linesplit) == 5 {
					val, _ := (*offsetmap)[linesplit[3]]
					outcode.AddStringCode("MOV " + registers[linesplit[2]] + ", [EBP-" + strconv.Itoa(val) + "]")
					log[index-startindex].assemblycode.AddStringCode("MOV " + registers[linesplit[2]] + ", [EBP-" + strconv.Itoa(val) + "]")
				} else {
					outcode.AddStringCode("MOV " + registers[linesplit[2]] + ", " + linesplit[3])
					log[index-startindex].assemblycode.AddStringCode("MOV " + registers[linesplit[2]] + ", " + linesplit[3])
				}
			}
			break
		case "ISTR":
			val, isinoffset := (*offsetmap)[linesplit[2]]
			if !isinoffset {
				offset := 0
				offset = len((*offsetmap))
				offset += 4
				(*offsetmap)[linesplit[2]] = offset
				outcode.AddStringCode("MOV [EBP-" + strconv.Itoa(offset) + "], " + registers[linesplit[3]])
				log[index-startindex].assemblycode.AddStringCode("MOV [EBP-" + strconv.Itoa(offset) + "], " + registers[linesplit[3]])
			} else {
				outcode.AddStringCode("MOV [EBP-" + strconv.Itoa(val) + "], " + registers[linesplit[3]])
				log[index-startindex].assemblycode.AddStringCode("MOV [EBP-" + strconv.Itoa(val) + "], " + registers[linesplit[3]])
			}
			break
		case "FADD":
			outcode.AddStringCode("ADDSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("ADDSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "FSUB":
			outcode.AddStringCode("SUBSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("SUBSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "FMUL":
			outcode.AddStringCode("MULSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("MULSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			break
		case "FDIV":
			outcode.AddStringCode("DIVSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			log[index-startindex].assemblycode.AddStringCode("DIVSS " + registers[linesplit[2]])
			break
		case "FMOV":
			_, isreg := registers[linesplit[3]]
			if isreg {
				outcode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
				log[index-startindex].assemblycode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", " + registers[linesplit[3]])
			} else {
				if len(linesplit) == 5 {
					val, _ := (*offsetmap)[linesplit[3]]
					outcode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", [EBP-" + strconv.Itoa(val) + "]")
					log[index-1-startindex].assemblycode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", [EBP-" + strconv.Itoa(val) + "]")
				} else {
					count := 0
					val, isin := (*floatcountmap)[linesplit[3]]
					if !isin {
						count = len((*floatcountmap)) + 1
						(*floatcountmap)[linesplit[3]] = count
					} else {
						count = val
					}
					data.AddStringCode(".LFV" + strconv.Itoa(count) + ":")
					data.AddStringCode("DD " + linesplit[3])
					outcode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", DWORD [.LFV" + strconv.Itoa(count) + "]")
					log[index-startindex].assemblycode.AddStringCode("MOVSS " + registers[linesplit[2]] + ", DWORD [.LFV" + strconv.Itoa(count) + "]")
				}
			}
			break
		case "FSTR":
			val, isinoffset := (*offsetmap)[linesplit[2]]
			if !isinoffset {
				offset := 0
				offset = len((*offsetmap))
				offset += 4
				(*offsetmap)[linesplit[2]] = offset
				outcode.AddStringCode("MOVSS [EBP-" + strconv.Itoa(offset) + "], " + registers[linesplit[3]])
				log[index-startindex].assemblycode.AddStringCode("MOVSS [EBP-" + strconv.Itoa(offset) + "], " + registers[linesplit[3]])
			} else {
				outcode.AddStringCode("MOVSS [EBP-" + strconv.Itoa(val) + "], " + registers[linesplit[3]])
				log[index-startindex].assemblycode.AddStringCode("MOVSS [EBP-" + strconv.Itoa(val) + "], " + registers[linesplit[3]])
			}
			break
		}
	}
	outcode.AddStringCode("ret")
	precode.AddCode(outcode)
	return precode, data, log
}
