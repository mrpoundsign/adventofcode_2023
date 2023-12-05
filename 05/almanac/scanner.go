package almanac

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	WS
	NL

	// Identifiers and literals
	IDENT

	// Misc Characters
	COLON // :

	// Variables
	NUMBER

	FEOF = 0
)

func (t Token) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case WS:
		return "WS"
	case NL:
		return "NL"
	case IDENT:
		return "IDENT"
	case COLON:
		return "COLON"
	case NUMBER:
		return "NUMBER"
	default:
		return "ILLEGAL"
	}
}

type Scanner struct {
	r    *bufio.Reader
	Line int
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return FEOF
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isEOF(ch) {
		return EOF, ""
	}

	if isNewline(ch) {
		s.unread()
		s.Line++
		return s.scanNewline()
	}

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	}

	if isNumber(ch) {
		s.unread()
		return s.scanNumber()
	}

	if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	if isNumber(ch) {
		s.unread()
		return s.scanNumber()
	}

	if ch == ':' {
		return COLON, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	ch := s.read()
	for isWhitespace(ch) {
		ch = s.read()
	}
	s.unread()
	return WS, string(ch)
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	ch := s.read()
	for isLetter(ch) {
		buf.WriteRune(ch)
		ch = s.read()
	}
	s.unread()

	// not unreading because we're at the end of the identifier
	return IDENT, buf.String()
}

func (s *Scanner) scanNumber() (tok Token, lit string) {
	var buf bytes.Buffer
	ch := s.read()
	for isNumber(ch) {
		buf.WriteRune(ch)
		ch = s.read()
	}
	s.unread()
	return NUMBER, buf.String()
}

func (s *Scanner) scanNewline() (tok Token, lit string) {
	ch := s.read()
	for isNewline(ch) {
		ch = s.read()
	}
	s.unread()
	return NL, string(ch)
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '-'
}

func isNewline(ch rune) bool {
	return ch == '\n'
}

func isNumber(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isEOF(ch rune) bool {
	return ch == FEOF
}
