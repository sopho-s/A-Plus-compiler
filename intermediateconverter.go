package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func MakeIntermediate(AST *node) (code, int) {
	var returncode code
	if AST.IsLeaf() {
		if AST.token == IDENT {
			if AST.variable.vartype == INTEGER {
				returncode.store = strconv.Itoa(AST.linenumber) + " IMOV IR3 " + AST.value + " IDENT\n" + strconv.Itoa(AST.linenumber) + " PUSH IR3"
			} else if AST.variable.vartype == FLOATING {
				returncode.store = strconv.Itoa(AST.linenumber) + " FMOV FR3 " + AST.value + " IDENT\n" + strconv.Itoa(AST.linenumber) + " PUSH FR3"
			}
			returncode.linecount = 2
			return returncode, AST.variable.vartype
		} else {
			if AST.token == INT {
				returncode.store = strconv.Itoa(AST.linenumber) + " IMOV IR3 " + AST.value + "\n" + strconv.Itoa(AST.linenumber) + " PUSH IR3"
			} else if AST.token == FLOAT {
				floatval, _ := strconv.ParseFloat(AST.value, 32)
				hexstring := fmt.Sprintf("%X", math.Float32bits(float32(floatval)))
				returncode.store = strconv.Itoa(AST.linenumber) + " FMOV FR3 0x" + hexstring + "\n" + strconv.Itoa(AST.linenumber) + " PUSH FR3"
			}
			returncode.linecount = 2
			if AST.token == INT {
				return returncode, INTEGER
			} else if AST.token == FLOAT {
				return returncode, FLOATING
			}
			return returncode, VOID
		}
	}
	switch AST.token {
	case ADD:
		leftcode, lefttype := MakeIntermediate(AST.children[0])
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 4
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " POP IR1\n" + strconv.Itoa(AST.linenumber) + " IADD IR2 IR1\n" + strconv.Itoa(AST.linenumber) + " PUSH IR2"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " POP FR1\n" + strconv.Itoa(AST.linenumber) + " FADD FR2 FR1\n" + strconv.Itoa(AST.linenumber) + " PUSH FR2"
		}
		leftcode.AddCode(rightcode)
		leftcode.AddCode(tempcode)
		returncode.AddCode(leftcode)
		return returncode, lefttype
	case SUB:
		leftcode, lefttype := MakeIntermediate(AST.children[0])
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 4
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " POP IR1\n" + strconv.Itoa(AST.linenumber) + " ISUB IR1 IR2\n" + strconv.Itoa(AST.linenumber) + " PUSH IR1"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " POP FR1\n" + strconv.Itoa(AST.linenumber) + " FSUB FR1 FR2\n" + strconv.Itoa(AST.linenumber) + " PUSH FR2"
		}
		leftcode.AddCode(rightcode)
		leftcode.AddCode(tempcode)
		returncode.AddCode(leftcode)
		return returncode, lefttype
	case IMUL:
		leftcode, lefttype := MakeIntermediate(AST.children[0])
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 4
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " POP IR1\n" + strconv.Itoa(AST.linenumber) + " IMUL IR2 IR1\n" + strconv.Itoa(AST.linenumber) + " PUSH IR2"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " POP FR1\n" + strconv.Itoa(AST.linenumber) + " FMUL FR2 FR1\n" + strconv.Itoa(AST.linenumber) + " PUSH FR2"
		}
		leftcode.AddCode(rightcode)
		leftcode.AddCode(tempcode)
		returncode.AddCode(leftcode)
		return returncode, lefttype
	case IDIV:
		leftcode, lefttype := MakeIntermediate(AST.children[0])
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 4
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " POP IR1\n" + strconv.Itoa(AST.linenumber) + " IDIV IR2\n" + strconv.Itoa(AST.linenumber) + " PUSH IR1"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " POP FR1\n" + strconv.Itoa(AST.linenumber) + " FDIV FR1 FR2\n" + strconv.Itoa(AST.linenumber) + " PUSH FR1"
		}
		leftcode.AddCode(rightcode)
		leftcode.AddCode(tempcode)
		returncode.AddCode(leftcode)
		return returncode, lefttype
	case ASSIGN:
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 2
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " ISTR " + AST.children[0].value + " IR2"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " FSTR " + AST.children[0].value + " FR2"
		}
		returncode.AddCode(rightcode)
		returncode.AddCode(tempcode)
		return returncode, righttype
	case RETURN:
		rightcode, righttype := MakeIntermediate(AST.children[1])
		var tempcode code
		tempcode.linecount = 2
		if righttype == INTEGER {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP IR2\n" + strconv.Itoa(AST.linenumber) + " PUSH IR2\n" + strconv.Itoa(AST.linenumber) + " ISUB SP 4\n" + strconv.Itoa(AST.linenumber) + " RET"
		} else if righttype == FLOATING {
			tempcode.store = strconv.Itoa(AST.linenumber) + " POP FR2\n" + strconv.Itoa(AST.linenumber) + " PUSH FR2\n" + strconv.Itoa(AST.linenumber) + " ISUB SP 4\n" + strconv.Itoa(AST.linenumber) + " RET"
		}
		returncode.AddCode(rightcode)
		returncode.AddCode(tempcode)
		return returncode, righttype
	default:
		break
	}
	return returncode, 0
}

func OptimiseIntermediate(intcode string) string {
	intcode = RemoveAdjacentPushPops(intcode)
	intcode = RemoveRandomMovs(intcode)
	intcode = RemovePointlessPushPops(intcode)
	intcode = RemoveRandomMovs(intcode)
	intcode = RemoveSelfReferenceMovs(intcode)
	intcode = RemoveSwitching(intcode)
	intcode = RemoveMovsForStr(intcode)
	return intcode
}

func RemoveAdjacentPushPops(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	var returncode code
	i := 0
	for i < len(splitcode)-1 {
		line1 := strings.Split(splitcode[i], " ")
		line2 := strings.Split(splitcode[i+1], " ")
		if line1[1] == "PUSH" && line2[1] == "POP" {
			if string(line2[2][0]) == "I" {
				returncode.AddStringCode(line2[0] + " " + "IMOV " + line2[2] + " " + line1[2])
			} else if string(line2[2][0]) == "F" {
				returncode.AddStringCode(line2[0] + " " + "FMOV " + line2[2] + " " + line1[2])
			}
			i++
		} else {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	if i == len(splitcode)-1 {
		returncode.AddStringCode(splitcode[i])
	}
	return returncode.store
}

func RemoveRandomMovs(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	var returncode code
	i := 0
	for i < len(splitcode)-1 {
		line1 := strings.Split(splitcode[i], " ")
		line2 := strings.Split(splitcode[i+1], " ")
		if i < len(splitcode)-2 {
			line3 := strings.Split(splitcode[i+2], " ")
			if line3[1] == "ISUB" || line3[1] == "IDIV" || line3[1] == "FSUB" || line3[1] == "FDIV" {
				returncode.AddStringCode(splitcode[i])
				i++
				continue
			}
		}
		if (line1[1] == "IMOV" && line2[1] == "IMOV") || (line1[1] == "FMOV" && line2[1] == "FMOV") {
			if line1[2] == line2[3] {
				if len(line1) == 4 {
					returncode.AddStringCode(line2[0] + " " + line2[1] + " " + line2[2] + " " + line1[3])
				} else {
					returncode.AddStringCode(line2[0] + " " + line2[1] + " " + line2[2] + " " + line1[3] + " " + line1[4])
				}
				i++
			} else {
				returncode.AddStringCode(splitcode[i])
			}
		} else {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	if i == len(splitcode)-1 {
		returncode.AddStringCode(splitcode[i])
	}
	return returncode.store
}

func RemoveSelfReferenceMovs(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	var returncode code
	i := 0
	for i < len(splitcode) {
		line := strings.Split(splitcode[i], " ")
		if !((line[1] == "IMOV" && (line[2] == line[3])) || (line[1] == "FMOV" && (line[2] == line[3]))) {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	if i == len(splitcode)-1 {
		returncode.AddStringCode(splitcode[i])
	}
	return returncode.store
}

func RemovePointlessPushPops(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	ignore := make([]bool, len(splitcode))
	var returncode code
	i := 0
	for i < len(splitcode)-1 {
		line := strings.Split(splitcode[i], " ")
		if line[1] == "POP" {
			reg := line[2]
			t := i
			for t > 0 {
				t--
				checkline := strings.Split(splitcode[t], " ")
				if checkline[2] == reg {
					break
				} else if len(checkline) >= 4 {
					if checkline[3] == reg {
						break
					}
				} else if checkline[1] == "PUSH" {
					if string(checkline[2][0]) == "I" {
						splitcode[t] = line[0] + " " + "IMOV " + line[2] + " " + checkline[2]
					} else if string(checkline[2][0]) == "F" {
						splitcode[t] = line[0] + " " + "FMOV " + line[2] + " " + checkline[2]
					}
					ignore[i] = true
					break
				}
			}
		}
		i++
	}
	i = 0
	for i <= len(splitcode)-1 {
		if !ignore[i] {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	return returncode.store
}

func RemoveSwitching(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	ignore := make([]bool, len(splitcode))
	var returncode code
	i := 0
	for i < len(splitcode)-1 {
		line := strings.Split(splitcode[i], " ")
		if len(line) >= 4 {
			if i >= 2 {
				if line[1] != "ISTR" && line[1] != "IMOV" && line[1] != "FSTR" && line[1] != "FMOV" {
					prevline1 := strings.Split(splitcode[i-1], " ")
					prevline2 := strings.Split(splitcode[i-2], " ")
					if prevline1[1] == "POP" {
						if prevline2[1] == "IMOV" || prevline2[1] == "FMOV" {
							if prevline1[2] != prevline2[2] && prevline1[2] == prevline2[3] {
								ignore[i-2] = true
								splitcode[i-1] = prevline1[0] + " " + prevline1[1] + " " + prevline2[2]
							}
						}
					} else if prevline1[1] == "IMOV" || prevline1[2] == "FMOV" {
						if prevline2[1] == "IMOV" || prevline2[2] == "FMOV" {
							if prevline1[2] != prevline2[2] && prevline1[2] == prevline2[3] {
								ignore[i-2] = true
								if len(prevline1) == 4 {
									splitcode[i-1] = prevline1[0] + " " + prevline1[1] + " " + prevline2[2] + " " + prevline1[3]
								} else {
									splitcode[i-1] = prevline1[0] + " " + prevline1[1] + " " + prevline2[2] + " " + prevline1[3] + " " + prevline1[4]
								}
							}
						}
					}
				}
			}
		}
		i++
	}
	i = 0
	for i <= len(splitcode)-1 {
		if !ignore[i] {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	return returncode.store
}

func RemoveMovsForStr(intcode string) string {
	splitcode := strings.Split(intcode, "\n")
	var returncode code
	i := 0
	for i < len(splitcode)-1 {
		line1 := strings.Split(splitcode[i], " ")
		line2 := strings.Split(splitcode[i+1], " ")
		if (line1[1] == "IMOV" && line2[1] == "ISTR") || (line1[1] == "FMOV" && line2[1] == "FSTR") {
			if line1[3] == "IR1" || line1[3] == "IR2" || line1[3] == "IR3" || line1[3] == "IR4" || line1[3] == "FR1" || line1[3] == "FR2" || line1[3] == "FR3" || line1[3] == "FR4" {
				if line1[2] == line2[3] {
					if len(line1) == 4 {
						if string(line1[3][0]) == "I" {
							returncode.AddStringCode(line2[0] + "ISTR " + line2[2] + ", " + line1[3])
						} else if string(line1[3][0]) == "F" {
							returncode.AddStringCode(line2[0] + "FSTR " + line2[2] + ", " + line1[3])
						}
						i++
					} else {
						returncode.AddStringCode(splitcode[i])
					}
				} else {
					returncode.AddStringCode(splitcode[i])
				}
			} else {
				returncode.AddStringCode(splitcode[i])
			}
		} else {
			returncode.AddStringCode(splitcode[i])
		}
		i++
	}
	if i == len(splitcode)-1 {
		returncode.AddStringCode(splitcode[i])
	}
	return returncode.store
}
