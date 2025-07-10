package lexer_test

import (
	"context"
	"testing"
	"yal/lexer"
)

// TODO: write actual tests

func TestComments(t *testing.T) {
	t.Run("Test multiLineCommentState()", func(t *testing.T) {
		src := ` /* aaa
                    bbbb
                    ccc
                */`

		l := lexer.NewLexer(context.Background(), string(src))

		tokens, err := l.Scan()
		if err != nil {
			panic(err)
		}

		if len(tokens) != 1 && tokens[0].TokenType != lexer.Eof {
			t.Errorf("%s got %+v\n", "there should be only one token here: Eof", tokens)
		}
	})

	t.Run("Test oneLineCommentState()", func(t *testing.T) {
		src := `// aaaa bbbb ccc`

		l := lexer.NewLexer(context.Background(), string(src))

		tokens, err := l.Scan()
		if err != nil {
			panic(err)
		}

		if len(tokens) != 1 && tokens[0].TokenType != lexer.Eof {
			t.Errorf("%s got %+v\n", "there should be only one token here: Eof", tokens)
		}
	})
}
