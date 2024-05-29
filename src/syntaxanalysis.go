package main

import "fmt"

func SeperateFunctions(nodes []node, bl *buildlog) ([]nodefunction, bool) {
	var functionlist []nodefunction
	infunc := false
	index := 0
	var function nodefunction
	for {
		if index == len(nodes) {
			functionlist = append(functionlist, function)
			break
		}
		if nodes[index].token == FUNC {
			if infunc {
				functionlist = append(functionlist, function)
				function.nodes = make([]node, 0)
				function.parameters = make([]node, 0)
				function.returns = make([]node, 0)
			} else {
				infunc = true
			}
			index++
			function.name = nodes[index].value
			index++
			if nodes[index].token != OPENBRACKET {
				bl.AddLog("Function parameter brackets not set correctly", 2)
				return nil, false
			}
			for {
				index++
				if nodes[index].token == CLOSEBRACKET {
					break
				}
				val, isin := types[nodes[index].value]
				if isin {
					index++
					nodes[index].variable.varname = nodes[index].value
					nodes[index].variable.vartype = val
					function.parameters = append(function.parameters, nodes[index])
				}
				index++
				if nodes[index].token == CLOSEBRACKET {
					break
				} else if nodes[index].token == COMMA {
					continue
				} else {
					bl.AddLog("Function parameter brackets not set correctly", 2)
					return nil, false
				}
			}
			index++
			if nodes[index].token != OPENBRACKET {
				bl.AddLog("Function return brackets not set correctly", 2)
				return nil, false
			}
			for {
				index++
				if nodes[index].token == CLOSEBRACKET {
					break
				}
				val, isin := types[nodes[index].value]
				if isin {
					nodes[index].variable.vartype = val
					function.returns = append(function.returns, nodes[index])
				}
				index++
				if nodes[index].token == CLOSEBRACKET {
					break
				} else if nodes[index].token == COMMA {
					continue
				} else {
					bl.AddLog("Function return brackets not set correctly", 2)
					return nil, true
				}
			}
			index++
			bl.AddLog("Found function \""+function.name+"\"", 0)
		} else {
			function.nodes = append(function.nodes, nodes[index])
			index++
		}
	}
	return functionlist, true
}

func SeperateSections(nodes []node, JMPlabel *int, bl *buildlog) []node {
	index := 0
	returnnodes := make([]node, 0)
	lastconditionindex := 0
	isgettingsection := false
	for {
		if index == len(nodes) {
			break
		}
		if !isgettingsection && nodes[index].token == IF {
			isgettingsection = true
			lastconditionindex = index
			index++
		} else if isgettingsection && nodes[index].token != OPENCBRACKET {
			nodes[lastconditionindex].condition = append(nodes[lastconditionindex].condition, nodes[index])
			index++
		} else if isgettingsection {
			isgettingsection = false
			tempnodes := nodes[lastconditionindex].condition
			outqueue := MakePostfix(tempnodes)
			AST := ConvertPostfix(outqueue).children[0]
			newcode, _ := MakeIntermediate(AST, JMPlabel)
			nodes[lastconditionindex].conditionintcode = &newcode
			returnnodes = append(returnnodes, nodes[lastconditionindex])
			returnnodes = append(returnnodes, nodes[index])
			index++
		} else {
			returnnodes = append(returnnodes, nodes[index])
			index++
		}
	}
	return returnnodes
}

func SyntaxAnalysis(nodes []node, bl buildlog, noerror bool, predefinedvar variablelist, predefinedfunctions definedfunctions) ([]node, bool, int, buildlog) {
	returnval := true
	shouldcontinue := CheckBrackets(nodes, &bl, noerror)
	if !shouldcontinue {
		returnval = false
	}
	nodes, shouldcontinue, vl := MakeType(nodes, &bl, noerror, predefinedvar)
	if !shouldcontinue {
		returnval = false
	}
	nodes, shouldcontinue = CheckVariables(nodes, vl, &bl, noerror, predefinedfunctions)
	if !shouldcontinue {
		returnval = false
	}
	if noerror && !returnval {
		fmt.Println("Lmao imagine making a syntax error, I am not gonna tell you where it is LOL, Gl nerd")
		bl.AddLog(fmt.Sprintf("Lmao imagine making a syntax error, I am not gonna tell you where it is LOL, Gl nerds"), 2)
	}
	return nodes, returnval, vl.count, bl
}

func CheckVariables(nodes []node, vl variablelist, bl *buildlog, noerror bool, predefinedfunctions definedfunctions) ([]node, bool) {
	returnval := true
	for index, value := range nodes {
		if value.token == IDENT {
			var vari variable
			vari.varname = value.value
			if !vl.CheckVariableInList(vari) {
				if !predefinedfunctions.CheckFunctionExists(vari.varname) {
					returnval = false
					if !noerror {
						fmt.Printf("Variable \""+value.value+"\" is undefined, line:%d, column:%d\n", value.linenumber, value.columnnumber)
						logval := fmt.Sprintf("Variable \""+value.value+"\" is undefined, line:%d, column:%d", value.linenumber, value.columnnumber)
						bl.AddLog(logval, 2)
					}
				} else {
					nodes[index].token = FUNC
					if nodes[index+1].token == OPENBRACKET {
						if nodes[index+2].token == CLOSEBRACKET {
							nodes[index].isbeingcalled = true
						}
					}
				}
			} else {
				nodes[index].variable.vartype = vl.GetVariableType(vari)
				nodes[index].variable.varname = value.value
			}
		}
	}
	return nodes, returnval
}

func MakeType(nodes []node, bl *buildlog, noerror bool, predefinedvar variablelist) ([]node, bool, variablelist) {
	var vl variablelist
	vl.AddList(predefinedvar)
	var returnnodes []node
	returnval := true
	for index, value := range nodes {
		val, isin := types[value.value]
		if isin {
			if nodes[index+1].token != IDENT {
				returnval = false
				if !noerror {
					fmt.Printf("Variable type identifier given to non variable, line:%d, column:%d\n", value.linenumber, value.columnnumber)
					logval := fmt.Sprintf("Variable type identifier given to non variable, line:%d, column:%d", value.linenumber, value.columnnumber)
					bl.AddLog(logval, 2)
				}
			} else {
				nodes[index+1].variable.vartype = val
				nodes[index+1].variable.varname = nodes[index+1].value
				if vl.CheckVariableInList(nodes[index+1].variable) {
					returnval = false
					if !noerror {
						fmt.Printf("Redefinition of variable, line:%d, column:%d\n", value.linenumber, value.columnnumber)
						logval := fmt.Sprintf("Redefinition of variable, line:%d, column:%d", value.linenumber, value.columnnumber)
						bl.AddLog(logval, 2)
					}
				} else {
					vl.Add(nodes[index+1].variable)
				}
			}
		} else {
			returnnodes = append(returnnodes, value)
		}
	}
	return returnnodes, returnval, vl
}

func CheckBrackets(nodes []node, bl *buildlog, noerror bool) bool {
	bracketcount := 0
	returnval := true
	var brackstack stack
	for _, value := range nodes {
		if value.token == OPENBRACKET {
			bracketcount++
			brackstack.Push(value)
		} else if value.token == CLOSEBRACKET {
			bracketcount--
			brackstack.Pop()
		}
		if bracketcount < 0 && value.token == CLOSEBRACKET {
			bracketcount++
			returnval = false
			if !noerror {
				fmt.Printf("Closing bracket without opening bracket, line:%d, column:%d\n", value.linenumber, value.columnnumber)
				logval := fmt.Sprintf("Closing bracket without opening bracket, line:%d, column:%d", value.linenumber, value.columnnumber)
				bl.AddLog(logval, 2)
			}
		}
	}
	if bracketcount > 0 {
		for !brackstack.IsEmpty() {
			unclosedbracket := brackstack.Pop()
			returnval = false
			if !noerror {
				fmt.Printf("Open bracket not closed, line:%d, column:%d\n", unclosedbracket.linenumber, unclosedbracket.columnnumber)
				logval := fmt.Sprintf("Open bracket not closed, line:%d, column:%d", unclosedbracket.linenumber, unclosedbracket.columnnumber)
				bl.AddLog(logval, 2)
			}
		}
		returnval = false
	}
	return returnval
}
