package core

// operators
var EXPECTED_KEYWORDS = []int{
	ADD,
	SUB,
	DIV,
	MUL,
	MOD,
	PUT,
	COLON,
	IF,
	EQUAL,
}

type Token struct {
	Pos  int
	Line int
	Type int
	Raw  string
}

const (
	UNKNOWN     = iota + 1
	FLOAT       // 0.0
	STRING      // "text"
	ADD         // +
	SUB         // -
	DIV         // /
	MUL         // *
	PUT         // .
	MOD         // %
	COLON       // :
	LEFT_BRACE  // (
	RIGHT_BRACE // )
	IF          // ?
	EQUAL       // =
	IDENT       // ([a-z]|_)+
	BOOL        // true | false
	EOF
)

var TOKEN_NAME_MAP = map[int]string{
	UNKNOWN:     "UNKNOWN",
	FLOAT:       "FLOAT",
	STRING:      "STRING",
	ADD:         "+",
	SUB:         "-",
	DIV:         "/",
	MUL:         "*",
	PUT:         ".",
	MOD:         "%",
	COLON:       ":",
	IF:          "?",
	EQUAL:       "=",
	IDENT:       "IDENT",
	LEFT_BRACE:  "(",
	RIGHT_BRACE: ")",
	BOOL:        "BOOL",
	EOF:         "EOF",
}
