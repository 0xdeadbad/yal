package lexer

// Lexer struct responsible to extract all tokens from given string source
type Lexer struct {
	source   string
	current  uint64
	start    uint64
	line     uint64
	column   uint64
	keywords map[string]TokenType
	tokenCh  chan Token
	state    stateFn
}

// Returns a new Lexer with the given source string and predefined keywords
func NewLexer(source string) *Lexer {
	keywords := make(map[string]TokenType)

	keywords["if"] = If
	keywords["else"] = Else
	keywords["for"] = For
	keywords["while"] = While
	keywords["let"] = Let
	keywords["return"] = Return
	keywords["fn"] = Fn
	keywords["true"] = True
	keywords["false"] = False
	keywords["NULL"] = Null
	keywords["switch"] = Switch
	keywords["goto"] = Goto

	return &Lexer{
		source:   source,
		current:  0,
		start:    0,
		line:     1,
		column:   0,
		keywords: keywords,
		tokenCh:  make(chan Token, 2),
		state:    stateMatch,
	}
}

func (l *Lexer) advance() byte {
	if l.isEof() {
		return 0
	}
	c := l.source[l.current]
	l.current++
	l.column++

	return c
}

// Scans the tokens from the given source and return a list of scanned tokens
func (l *Lexer) Scan() ([]Token, error) {
	tokens := []Token{}

	for token := l.NextToken(); token.Type != Eof; token = l.NextToken() {
		tokens = append(tokens[:], token)
	}

	tokens = append(tokens[:], l.newToken(Eof))

	close(l.tokenCh)

	return tokens, nil
}

func (l *Lexer) advanceEmit(token_type TokenType) {
	l.advance()
	l.emit(token_type)
}

func (l *Lexer) emit(token_type TokenType) {
	token := l.newToken(token_type)
	l.tokenCh <- token
}

func (l *Lexer) newToken(token_type TokenType) Token {
	token := Token{
		Type:   token_type,
		Lexeme: l.source[l.start:l.current],
		Line:   l.line,
		Column: l.column,
	}
	l.start = l.current
	return token
}

func (l *Lexer) peek() byte {
	if l.isEof() {
		return 0
	}

	return l.source[l.current]
}

func (l *Lexer) peekNext() byte {
	if l.isEof() && l.current+1 < uint64(len(l.source)) {
		return 0
	}

	return l.source[l.current+1]
}

func (l *Lexer) peekMatch(c byte) bool {
	if l.isEof() {
		return false
	}

	if l.source[l.current] == c {
		l.current++
		return true
	}

	return false
}

func (l *Lexer) isEof() bool {
	return l.current >= uint64(len(l.source))
}

func (l *Lexer) previous() byte {
	return l.source[l.current-1]
}

func (l *Lexer) ignore() {
	l.start = l.current
}

func (l *Lexer) backup() {
	l.current--
	l.start = l.current
}

// Return the next token available to be consumed
func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.tokenCh:
			return token
		default:
			l.state = l.state(l)
		}
	}
}
