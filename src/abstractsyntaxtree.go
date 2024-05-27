package main

func MakePostfix(nodes []node) section {
	var out section
	var offstack stack
	var output queue
	var il idlist
	il.count = 0
	for _, value := range nodes {
		if value.token == IDENT || value.token == INT || value.token == FLOAT || value.token == FUNC {
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
	return out
}

func ConvertPostfix(postfixcode section) node {
	var sectionnode node
	for _, line := range postfixcode.lines {
		var nodestack stack
		for !line.IsEmpty() {
			currnode := line.Dequeue()
			if currnode.token == IDENT || currnode.token == INT || currnode.token == FLOAT {
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
		tempnode := nodestack.Pop()
		sectionnode.LinkNode(&tempnode)
	}
	return sectionnode
}
