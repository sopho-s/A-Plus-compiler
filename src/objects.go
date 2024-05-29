package main

import "io"

type buf []byte

func (b *buf) Read(p []byte) (n int, err error) {
	n = copy(p, *b)
	*b = (*b)[n:]
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (b *buf) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *buf) String() string { return string(*b) }

type node struct {
	id               int
	token            token
	value            string
	size             int
	children         []*node
	childrencount    int
	variable         variable
	linenumber       int
	columnnumber     int
	isloop           bool
	looptype         int
	isbeingcalled    bool
	condition        []node
	conditionintcode *code
	loopnum          int
}

type nodefunction struct {
	name       string
	nodes      []node
	parameters []node
	returns    []node
}

type definedfunctions struct {
	functions map[string]nodefunction
}

type idlist struct {
	count int
}

type variablelist struct {
	variables []variable
	count     int
}
type code struct {
	store     string
	node      node
	linecount int
}

type variable struct {
	id      int
	vartype int
	varname string
	size    int
	offset  int
}

type section struct {
	lines []queue
	name  string
}

type stack struct {
	stack []node
	count int
}

type queue struct {
	queue []node
	count int
}

type registertable struct {
	values []string
	filled []bool
}

func (df *definedfunctions) AddFunction(nf nodefunction) {
	df.functions[nf.name] = nf
}

func (df *definedfunctions) CheckFunctionExists(funcname string) bool {
	_, isin := df.functions[funcname]
	return isin
}

func (df *definedfunctions) CountFunctionParameters(funcname string) int {
	val, _ := df.functions[funcname]
	return len(val.parameters)
}

func (df *definedfunctions) CountFunctionReturns(funcname string) int {
	val, _ := df.functions[funcname]
	return len(val.returns)
}

func (il *idlist) GetId() int {
	il.count += 1
	return il.count - 1
}

func (c *code) AddCode(val code) {
	if val.store != "" {
		if c.linecount > 0 {
			c.store += "\n" + val.store
		} else {
			c.store = val.store
		}
		c.linecount += val.linecount
	}
}

func (c *code) AddStringCode(val string) {
	if c.linecount > 0 {
		c.store += "\n" + val
	} else {
		c.store = val
	}
	c.linecount += 1
}

func (n *node) IsLeaf() bool {
	if n.size == 0 {
		return true
	} else {
		return false
	}
}

func (vl *variablelist) Add(val variable) {
	for _, value := range vl.variables {
		if value.varname == val.varname {
			return
		}
	}
	vl.variables = append(vl.variables, val)
	vl.count += 1
}

func (vl *variablelist) CheckVariableInList(val variable) bool {
	for _, value := range vl.variables {
		if value.varname == val.varname {
			return true
		}
	}
	return false
}
func (vl *variablelist) GetVariableType(val variable) int {
	for _, value := range vl.variables {
		if value.varname == val.varname {
			return value.vartype
		}
	}
	return -1
}

func (vl *variablelist) AddList(val variablelist) {
	vl.variables = append(vl.variables, val.variables...)
	vl.count += val.count
}

func (vl *variablelist) PopBack() variable {
	if vl.count > 0 {
		variable := vl.variables[len(vl.variables)-1]
		vl.variables = vl.variables[:len(vl.variables)-1]
		vl.count -= 1
		return variable
	} else {
		var temp variable
		return temp
	}
}

func (n *node) LinkNode(val *node) {
	n.children = append(n.children, val)
	n.size++
}

func (n *nodefunction) RemoveStartAndEnd() {
	n.nodes = n.nodes[1 : len(n.nodes)-1]
}

func (s *stack) Push(val node) {
	s.stack = append(s.stack, val)
	s.count += 1
}

func (s *stack) Pop() node {
	if s.count > 0 {
		node := s.stack[len(s.stack)-1]
		s.stack = s.stack[:len(s.stack)-1]
		s.count -= 1
		return node
	} else {
		var temp node
		return temp
	}
}

func (s *stack) Peek() node {
	if s.count > 0 {
		node := s.stack[len(s.stack)-1]
		return node
	} else {
		var temp node
		return temp
	}
}

func (s *stack) IsEmpty() bool {
	if s.count == 0 {
		return true
	}
	return false
}

func (q *queue) Enqueue(val node) {
	q.queue = append(q.queue, val)
	q.count += 1
}

func (q *queue) Dequeue() node {
	if q.count > 0 {
		node := q.queue[0]
		q.queue = q.queue[1:]
		q.count -= 1
		return node
	} else {
		var temp node
		return temp
	}
}

func (q *queue) Clear() {
	q.queue = make([]node, 0)
	q.count = 0
}

func (q *queue) IsEmpty() bool {
	if q.count == 0 {
		return true
	}
	return false
}

func (t token) String() string {
	return tokens[t]
}

func (t token) int() int {
	return int(t)
}
