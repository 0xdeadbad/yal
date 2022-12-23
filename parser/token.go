package parser

import "fmt"

// Token struct holding the lexeme and its position on the source code
type Token struct {
	Type   TokenType
	Lexeme string
	Line   uint64
	Column uint64
}

// Pretty printing for the Token struct
func (t Token) String() string {
	return fmt.Sprintf("{\n\t\tType: %+v\n\t\tLexeme: %s\n\t\tLine: %d\n\t\tColumn: %d\n\t}", t.Type, t.Lexeme, t.Line, t.Column)
}

type TokenType uint16

// Enum for all possible Tokens
const (
	Eof TokenType = iota

	LeftParen
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Comma
	Colon
	Semicolon
	Dot

	Star
	Minus
	Plus
	Slash
	Inc
	Dec

	Bang
	BangEqual
	Equal
	EqualEqual
	Lesser
	LesserEqual
	Greater
	GreaterEqual
	Ampersand
	DoubleAmpersand
	Pipe
	DoublePipe
	PlusEqual
	MinusEqual
	StarEqual
	SlashEqual
	Xor
	XorEqual
	Rem
	Shl
	Shr

	LeftArrow
	RightArrow
	FuncArrow

	If
	Else
	For
	While
	Let
	Return
	Fn
	True
	False
	Null
	Switch
	Goto

	Identifier
	String
	Number16
	Number10
	Number8
	Number2
)

// String formating for the Tokens enums
func (t TokenType) String() string {
	switch t {
	case Eof:
		return "EOF"
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case LeftBracket:
		return "["
	case RightBracket:
		return "]"

	case Comma:
		return ","
	case Colon:
		return ":"
	case Semicolon:
		return ";"
	case Dot:
		return "."

	case Star:
		return "*"
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Slash:
		return "/"
	case Inc:
		return "++"
	case Dec:
		return "--"

	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case EqualEqual:
		return "=="
	case Lesser:
		return "<"
	case LesserEqual:
		return "<="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Ampersand:
		return "&"
	case DoubleAmpersand:
		return "&&"
	case Pipe:
		return "|"
	case DoublePipe:
		return "||"
	case PlusEqual:
		return "+="
	case MinusEqual:
		return "-="
	case StarEqual:
		return "*="
	case SlashEqual:
		return "/="
	case Xor:
		return "^"
	case XorEqual:
		return "^="
	case Rem:
		return "%"
	case Shl:
		return "<<"
	case Shr:
		return ">>"

	case LeftArrow:
		return "->"
	case RightArrow:
		return "<-"
	case FuncArrow:
		return "=>"

	case If:
		return "if"
	case Else:
		return "else"
	case For:
		return "for"
	case While:
		return "while"
	case Let:
		return "let"
	case Return:
		return "return"
	case Fn:
		return "fn"
	case True:
		return "true"
	case False:
		return "false"
	case Null:
		return "NULL"

	case Identifier:
		return "identifier"
	case String:
		return "string"
	case Number16:
		return "number(hex)"
	case Number10:
		return "number(decimal)"
	case Number8:
		return "number(octal)"
	case Number2:
		return "number(binary)"
	}

	panic("Unknown token")
}

const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

func (op TokenType) Precedence() int {
	switch op {
	case DoublePipe:
		return 1
	case DoubleAmpersand:
		return 2
	case Equal, BangEqual, Lesser, LesserEqual, Greater, GreaterEqual:
		return 3
	case Plus, Minus, Pipe, Xor:
		return 4
	case Star, Slash, Rem, Shl, Shr, Ampersand:
		return 5
	}
	return LowestPrec
}
