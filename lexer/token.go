package lexer

import (
	"encoding/json"
	"fmt"
)

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

	TypeAnn
)

// String formating for the Tokens enums
func (t TokenType) String() string {
	switch t {
	case Eof:
		return "EOF"
	case LeftParen:
		return "left parenthesis"
	case RightParen:
		return "right parenthesis"
	case LeftBrace:
		return "left brace"
	case RightBrace:
		return "right brace"
	case LeftBracket:
		return "left bracket"
	case RightBracket:
		return "right bracket"

	case Comma:
		return "comma"
	case Colon:
		return "colon"
	case Semicolon:
		return "semicolon"
	case Dot:
		return "dot"

	case Star:
		return "star"
	case Minus:
		return "minus"
	case Plus:
		return "plus"
	case Slash:
		return "slash"
	case Inc:
		return "increment"
	case Dec:
		return "decrement"

	case Bang:
		return "bang"
	case BangEqual:
		return "different"
	case Equal:
		return "assignment"
	case EqualEqual:
		return "equallity"
	case Lesser:
		return "lesser"
	case LesserEqual:
		return "lesser equal"
	case Greater:
		return "greater"
	case GreaterEqual:
		return "greater equal"
	case Ampersand:
		return "ampersand"
	case DoubleAmpersand:
		return "double ampersand"
	case Pipe:
		return "pipe"
	case DoublePipe:
		return "double pipe"
	case PlusEqual:
		return "plus assign"
	case MinusEqual:
		return "minus minus assign"
	case StarEqual:
		return "star assign"
	case SlashEqual:
		return "slash assign"
	case Xor:
		return "xor"
	case XorEqual:
		return "xor assign"
	case Rem:
		return "rem"
	case Shl:
		return "shift left"
	case Shr:
		return "shift right"

	case LeftArrow:
		return "left arrow"
	case RightArrow:
		return "right arrow"
	case FuncArrow:
		return "function arrow"

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

func (t TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
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
