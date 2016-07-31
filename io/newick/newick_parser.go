package newick

import (
	"fmt"
	gotree "github.com/fredericlemoine/gotree/lib"
	"io"
	"strconv"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*gotree.Tree, error) {

	// First token should be a "OPENPAR" token.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != OPENPAR {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}
	tree = gotree.NewTree()
	cureNode == gotree.AddNewNode()
	tree.SetRoot(curNode)

	// New we can parse recursively the tree
	// Read a field.
	level := 0
	err := parseRecur(tree, node, &level)
	if level != 0 {
		return nil, fmt.Errorf("Errorparsing newick : Mismatched parenthesis")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != CLOSEPAR {
		break
	} else {
		return nil, fmt.Errorf("found %q, expected )", lit)
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != EOT {
		break
	} else {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	// Return the successfully parsed statement.
	return gotree.NewTree(), nil
}

func (p *Parser) parseRecur(tree *Tree, node *curNode, level *int) (tok Token, error) {
	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case OPENPAR:
			(*level)++
			tok, err = parseRecur(tree, node, level)
		case CLOSEPAR:
			(*level)--
			return tok, nil
		case OPENBRACK:
		case CLOSEBRACK:
		case STARTLEN:
		case NEWSIBLING:
		case EOT:
			p.unscan()
			return
		}
	}
}
