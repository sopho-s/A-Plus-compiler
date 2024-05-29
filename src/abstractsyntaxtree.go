package main

func MakePostfix(nodes []node) section {
	var out section
	var offstack stack
	var output queue
	var il idlist
	il.count = 0
	for _, value := range nodes {
		if value.token == IDENT || value.token == INT || value.token == FLOAT || value.token == BOOL || value.token == FUNC {
			if value.token == IDENT {
				value.id = il.GetId()
			}
			output.Enqueue(value)
		} else if value.token == SEMI {
			for !offstack.IsEmpty() {
				if !(offstack.Peek().token == OPENBRACKET) || !(offstack.Peek().token == CLOSEBRACKET) {
					output.Enqueue(offstack.Pop())
				} else {
					_ = offstack.Pop()
				}
			}
			out.lines = append(out.lines, output)
			output.Clear()
		} else if value.token == OPENCBRACKET {
			for !offstack.IsEmpty() {
				if !(offstack.Peek().token == OPENBRACKET) || !(offstack.Peek().token == CLOSEBRACKET) {
					output.Enqueue(offstack.Pop())
				} else {
					_ = offstack.Pop()
				}
			}
			out.lines = append(out.lines, output)
			output.Clear()
		} else if value.token == CLOSECBRACKET {
			output.Enqueue(value)
			out.lines = append(out.lines, output)
			output.Clear()
		} else if offstack.IsEmpty() {
			value.id = il.GetId()
			offstack.Push(value)
		} else if value.token == OPENBRACKET {
			offstack.Push(value)
		} else if value.token == CLOSEBRACKET {
			for {
				popval := offstack.Pop()
				if popval.token == OPENBRACKET {
					break
				}
				output.Enqueue(popval)
				topstack := offstack.Peek()
				if topstack.token == OPENBRACKET {
					_ = offstack.Pop()
					break
				}
				if offstack.IsEmpty() || precedence[int(topstack.token)] < precedence[int(value.token)] {
					break
				}
			}
		} else {
			topstack := offstack.Peek()
			value.id = il.GetId()
			if topstack.token == OPENBRACKET {
				offstack.Push(value)
			} else if precedence[int(topstack.token)] < precedence[int(value.token)] {
				offstack.Push(value)
			} else {
				for {
					popval := offstack.Pop()
					output.Enqueue(popval)
					topstack = offstack.Peek()
					if topstack.token == OPENBRACKET {
						offstack.Push(value)
						break
					}
					if offstack.IsEmpty() || precedence[int(topstack.token)] < precedence[int(value.token)] {
						offstack.Push(value)
						break
					}
				}
			}
		}
	}
	if len(out.lines) == 0 {
		for !offstack.IsEmpty() {
			if !(offstack.Peek().token == OPENBRACKET) || !(offstack.Peek().token == CLOSEBRACKET) {
				output.Enqueue(offstack.Pop())
			} else {
				_ = offstack.Pop()
			}
		}
		out.lines = append(out.lines, output)
		output.Clear()
	}
	return out
}

func ConvertPostfix(postfixcode section) node {
	var sectionnode node
	var tempnodestack stack
	var sectionnodestack stack
	isinsection := false
	for _, line := range postfixcode.lines {
		var nodestack stack
		for !line.IsEmpty() {
			currnode := line.Dequeue()
			if currnode.token == IDENT || currnode.token == INT || currnode.token == FLOAT || currnode.token == BOOL || currnode.token == FUNC {
				nodestack.Push(currnode)
			} else if currnode.token == IF || currnode.token == DO {
				nodestack.Push(currnode)
			} else {
				var left node
				var right node
				left = nodestack.Pop()
				right = nodestack.Pop()
				currnode.LinkNode(&right)
				currnode.LinkNode(&left)
				nodestack.Push(currnode)
			}
		}
		if !isinsection && (nodestack.Peek().token != IF && nodestack.Peek().token != DO) {
			tempnode := nodestack.Pop()
			sectionnode.LinkNode(&tempnode)
		} else {
			currnode := nodestack.Pop()
			if currnode.token == IF || currnode.token == DO {
				sectionnodestack.Push(currnode)
				var openbrac node
				openbrac.token = OPENCBRACKET
				tempnodestack.Push(openbrac)
			}
			isinsection = true
			if currnode.token == CLOSECBRACKET {
				sectionno := sectionnodestack.Pop()
				for {
					tempnode := tempnodestack.Pop()
					if tempnode.token == OPENCBRACKET {
						if sectionnodestack.IsEmpty() {
							tempnodestack.Push(sectionno)
							isinsection = false
						} else {
							sectionnodestack.stack[sectionnodestack.count-1].LinkNode(&sectionno)
						}
						break
					} else if tempnode.token != CLOSECBRACKET {
						sectionno.LinkNode(&tempnode)
					}
				}
			}
			if sectionnodestack.IsEmpty() {
				tempnode := tempnodestack.Pop()
				sectionnode.LinkNode(&tempnode)
			} else if currnode.token != IF && currnode.token != DO {
				tempnodestack.Push(currnode)
			}
		}
	}
	return sectionnode
}
