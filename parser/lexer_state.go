package parser

import "fmt"

type stateFn func(*Lexer) stateFn

func defaultActionState(l *Lexer) stateFn {
	c := l.previous()

	if IsAlpha(c) || c == '_' {
		return identifierState
	} else if IsBase10(c) {
		return numberState
	}

	panic(fmt.Sprintf("error: invalid '%c' character", c))
}

func multiLineCommentState(l *Lexer) stateFn {
	for l.peek() != '*' && l.peekNext() != '/' {
		l.advance()
	}
	l.advance()
	l.advance()
	l.ignore()

	return stateMatch(l)
}

func oneLineCommentState(l *Lexer) stateFn {
	for l.peek() != '\n' && l.peek() != 0 {
		l.advance()
	}
	l.ignore()

	return stateMatch(l)
}

func numberState(l *Lexer) stateFn {
	for IsBase10(l.peek()) {
		l.advance()
	}
	if l.peekMatch('.') {
		for IsBase10(l.peek()) {
			l.advance()
		}
	}
	l.emit(Number10)

	return stateMatch
}

func stringState(l *Lexer) stateFn {
	for l.peek() != '"' {
		l.advance()
	}
	l.advance()

	l.emit(String)

	return stateMatch
}

func identifierState(l *Lexer) stateFn {
	for IsAlpha(l.peek()) || IsBase10(l.peek()) {
		l.advance()
	}
	if id, ok := l.keywords[l.source[l.start:l.current]]; ok {
		l.emit(id)
	} else {
		l.emit(Identifier)
	}

	return stateMatch
}

func simpleTokenState(l *Lexer) stateFn {
	c := l.previous()
	switch c {
	case ')':
		l.emit(RightParen)
	case '(':
		l.emit(LeftParen)
	case '}':
		l.emit(RightBrace)
	case '{':
		l.emit(LeftBrace)
	case ',':
		l.emit(Comma)
	case ':':
		l.emit(Colon)
	case ';':
		l.emit(Semicolon)
	case '.':
		l.emit(Dot)
	case ']':
		l.emit(RightBracket)
	case '[':
		l.emit(LeftBracket)
	}

	return stateMatch
}

func compoundTokenState(l *Lexer) stateFn {
	c := l.previous()
	switch c {
	case '!':
		switch l.peek() {
		case '=':
			l.advanceEmit(BangEqual)
		default:
			l.emit(Bang)
		}

	case '=':
		switch l.peek() {
		case '=':
			l.advanceEmit(EqualEqual)
		case '>':
			l.advanceEmit(FuncArrow)
		default:
			l.emit(Equal)
		}

	case '>':
		switch l.peek() {
		case '=':
			l.advanceEmit(GreaterEqual)
		case '>':
			l.advanceEmit(Shr)
		default:
			l.emit(Greater)
		}

	case '<':
		switch l.peek() {
		case '=':
			l.advanceEmit(LesserEqual)
		case '-':
			l.advanceEmit(RightArrow)
		case '<':
			l.advanceEmit(Shl)
		default:
			l.emit(Lesser)
		}
	case '-':
		switch l.peek() {
		case '>':
			l.advanceEmit(LeftArrow)
		case '-':
			l.advanceEmit(Dec)
		case '=':
			l.advanceEmit(MinusEqual)
		default:
			l.emit(Minus)
		}
	case '*':
		switch l.peek() {
		case '=':
			l.advanceEmit(StarEqual)
		default:
			l.emit(Star)
		}
	case '/':
		switch l.peek() {
		case '=':
			l.advanceEmit(SlashEqual)
		case '/':
			l.advance()
			return oneLineCommentState(l)
		case '*':
			l.advance()
			return multiLineCommentState(l)
		default:
			l.emit(Slash)
		}
	case '+':
		switch l.peek() {
		case '+':
			l.advanceEmit(Inc)
		case '=':
			l.advanceEmit(PlusEqual)
		default:
			l.emit(Plus)
		}
	case '^':
		switch l.peek() {
		case '=':
			l.advanceEmit(XorEqual)
		default:
			l.emit(Xor)
		}
	case '|':
		switch l.peek() {
		case '|':
			l.advanceEmit(DoublePipe)
		default:
			l.emit(Pipe)
		}
	case '&':
		switch l.peek() {
		case '&':
			l.advanceEmit(DoubleAmpersand)
		default:
			l.emit(Ampersand)
		}
	}

	return stateMatch
}

func ignoreState(l *Lexer) stateFn {
	l.start++
	return stateMatch(l)
}

func numberWithBaseState(l *Lexer) stateFn {
	switch l.peek() {
	case 'x':
		readBase16(l)
	case 'b':
		readBase2(l)
	case 'o':
		readBase8(l)
	default:
		readBase8(l)
	}

	return stateMatch
}

func readBase16(l *Lexer) {
	l.backup()
	l.advance()
	l.advance()
	for IsBase16(l.peek()) {
		l.advance()
	}
	if l.peekMatch('.') {
		for IsBase16(l.peek()) {
			l.advance()
		}
	}

	l.emit(Number16)
}

/*func readBase10(l *Lexer) {
	for IsBase10(l.peek()) {
		l.advance()
	}

	l.emit(Number10)
}*/

func readBase8(l *Lexer) {
	for IsBase8(l.peek()) {
		l.advance()
	}
	if l.peekMatch('.') {
		for IsBase8(l.peek()) {
			l.advance()
		}
	}

	l.emit(Number8)
}

func readBase2(l *Lexer) {
	l.backup()
	l.advance()
	l.advance()
	for IsBase2(l.peek()) {
		l.advance()
	}
	if l.peekMatch('.') {
		for IsBase2(l.peek()) {
			l.advance()
		}
	}

	l.emit(Number2)
}

func stateMatch(l *Lexer) stateFn {
	if l.isEof() {
		l.emit(Eof)
		return nil
	}
	c := l.advance()
	switch c {
	case ')', '(', '}', '{', ',', ':', ';', '.', '[', ']':
		return simpleTokenState

	case '!', '=', '>', '<', '-', '*', '/', '+', '|', '&':
		return compoundTokenState

	case ' ', '\t', '\r':
		return ignoreState

	case '\n':
		l.line++
		return ignoreState

	case '"':
		return stringState

	case '0':
		return numberWithBaseState

	case 0:
		l.emit(Eof)
		return nil

	default:
		return defaultActionState
	}
}
