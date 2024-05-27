package main

const (
	PROGRAM = iota

	EOF
	ILLEGAL
	IDENT
	INT
	FLOAT
	FUNC
	SEMI
	OPENBRACKET
	CLOSEBRACKET
	OPENCBRACKET
	CLOSECBRACKET
	RETURN

	// operators
	ADD
	SUB
	IMUL
	IDIV

	ASSIGN
	POINT
	COMMA
	PIPEIN

	INCLUDE
)

const (
	NOREGISTER = iota
	EAX
	EBX
	ECX
	EDX
	XMM0
	XMM1
	XMM2
	XMM3
)

const (
	VOID = iota
	INTEGER
	FLOATING
)

var types = map[string]int{
	"void":  VOID,
	"int":   INTEGER,
	"float": FLOATING,
}

var registers = map[string]string{
	"IR1": "EAX",
	"IR2": "EBX",
	"IR3": "ECX",
	"IR4": "EDX",
	"FR1": "XMM0",
	"FR2": "XMM1",
	"FR3": "XMM2",
	"FR4": "XMM3",
	"SP":  "ESP",
}

var tokens = []string{
	PROGRAM: "PROGRAM",

	EOF:           "EOF",
	ILLEGAL:       "ILLEGAL",
	IDENT:         "IDENT",
	INT:           "INT",
	FLOAT:         "FLOAT",
	FUNC:          "FUNC",
	SEMI:          ";",
	OPENBRACKET:   "(",
	CLOSEBRACKET:  ")",
	OPENCBRACKET:  "{",
	CLOSECBRACKET: "}",
	RETURN:        "RETURN",

	ADD:  "+",
	SUB:  "-",
	IMUL: "*",
	IDIV: "/",

	ASSIGN: "=",
	POINT:  ".",
	COMMA:  ",",
	PIPEIN: "<<<",

	INCLUDE: "INCLUDE",
}

var precedence = map[int]int{
	OPENCBRACKET:  -2,
	OPENBRACKET:   -2,
	CLOSECBRACKET: -1,
	CLOSEBRACKET:  -1,
	RETURN:        0,
	ASSIGN:        0,
	PIPEIN:        1,
	ADD:           2,
	SUB:           2,
	IMUL:          3,
	IDIV:          3,
}
