package main

const (
	PROGRAM = iota

	EOF
	ILLEGAL
	IDENT
	INT
	FLOAT
	BOOL
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
	AND
	OR

	// comparisons
	ISEQUAL
	ISNOTEQUAL
	GREATER
	GREATEROREQUAL
	LESS
	LESSOREQUAL

	ASSIGN
	POINT
	COMMA
	PIPEIN

	INCLUDE

	IF
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
	BOOLEAN
)

var types = map[string]int{
	"void":  VOID,
	"int":   INTEGER,
	"float": FLOATING,
	"bool":  BOOLEAN,
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
	BOOL:          "BOOL",
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
	AND:  "&",
	OR:   "|",

	ISEQUAL:        "==",
	ISNOTEQUAL:     "!=",
	GREATER:        ">",
	GREATEROREQUAL: ">=",
	LESS:           "<",
	LESSOREQUAL:    "<=",

	ASSIGN: "=",
	POINT:  ".",
	COMMA:  ",",
	PIPEIN: "<<<",

	INCLUDE: "INCLUDE",

	IF: "IF",
}

var precedence = map[int]int{
	OPENCBRACKET:   -2,
	OPENBRACKET:    -2,
	CLOSECBRACKET:  -1,
	CLOSEBRACKET:   -1,
	RETURN:         0,
	ASSIGN:         0,
	OR:             1,
	AND:            2,
	ISEQUAL:        3,
	ISNOTEQUAL:     3,
	GREATER:        3,
	GREATEROREQUAL: 3,
	LESS:           3,
	LESSOREQUAL:    3,
	ADD:            4,
	SUB:            4,
	IMUL:           4,
	IDIV:           5,
	PIPEIN:         6,
}
