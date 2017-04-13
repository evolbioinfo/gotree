package newick

import (
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/tree"
	"io"
	"strconv"
	"strings"
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
func (p *Parser) Parse() (*tree.Tree, error) {

	// First token should be a "OPENPAR" token.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != OPENPAR {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}
	p.unscan()
	tree := tree.NewTree()

	// Now we can parse recursively the tree
	// Read a field.
	level := 0
	_, err := p.parseRecur(tree, nil, &level)
	if err != nil {
		return nil, err
	}
	if level != 0 {
		return nil, fmt.Errorf("Newick Error : Mismatched parenthesis")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != EOT {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	tree.UpdateTipIndex()
	err = tree.ClearBitSets()
	if err != nil {
		return nil, err
	}
	tree.UpdateBitSet()
	// Not necessary at the parsing step...
	// may be too long to do each time
	//tree.ComputeDepths()
	// Return the successfully parsed statement.
	return tree, nil
}

func (p *Parser) parseRecur(t *tree.Tree, node *tree.Node, level *int) (Token, error) {

	var newNode *tree.Node = node
	var prevTok Token = -1
	var err error

	for {
		tok, lit := p.scanIgnoreWhitespace()

		switch tok {
		case OPENPAR:
			newNode = t.NewNode()
			if node == nil {
				if *level > 0 {
					return -1, errors.New("Nil node at depth > 0")
				}
				t.SetRoot(newNode)
				node = newNode
			} else {
				if *level == 0 {
					return -1, errors.New("An open parenthesis at level 0 of recursion... Forgot a ';' at the end of previous tree?")
				}
				t.ConnectNodes(node, newNode)
			}
			(*level)++
			prevTok, err = p.parseRecur(t, newNode, level)
			if err != nil {
				return -1, err
			}
			if prevTok != CLOSEPAR {
				return -1, errors.New("Newick Error: Mismatched parenthesis after parseRecur")
			}
		case CLOSEPAR:
			(*level)--
			return tok, nil
		case OPENBRACK:
			//if prevTok == OPENPAR || prevTok == NEWSIBLING || prevTok == -1 {
			if newNode == nil {
				return -1, errors.New("Newick Error: Comment should not be located here: " + lit)
			}
			tok, lit = p.scanIgnoreWhitespace()
			if tok != IDENT && tok != NUMERIC {
				return -1, errors.New("Newick Error: There should be a comment after [")
			}
			// Add comment to node
			newNode.AddComment(lit)
			tok, lit = p.scanIgnoreWhitespace()
			if tok != CLOSEBRACK {
				return -1, errors.New("Newick Error: There should be a ] after a comment")
			}
			prevTok = CLOSEBRACK
		case CLOSEBRACK:
			// Error here should not have
			return -1, errors.New("Newick Error: Mismatched ] here...")
		case STARTLEN:
			tok, lit = p.scanIgnoreWhitespace()
			if tok != NUMERIC {
				return -1, errors.New("Newick Error: No numeric value after :")
			}
			// We skip length if the length is assigned to the root node
			if newNode != nil && *level != 0 && newNode != node {
				e, err := newNode.ParentEdge()
				if err != nil {
					return -1, err
				}

				if e.Length() != -1 {
					return -1, errors.New("Newick Error: More than one length is given :" + lit)
				}
				length, errf := strconv.ParseFloat(lit, 64)
				if errf != nil {
					return -1, errors.New("Newick Error: Length is not a float value : " + lit)
				}
				e.SetLength(length)
			} else {
				if newNode.Name() == "" && prevTok != ')' && *level != 0 {
					return -1, errors.New("Newick Error: Cannot assign length to nil node :" + lit)
				}
			}
			prevTok = tok
		case NEWSIBLING:
			newNode = nil
			prevTok = NEWSIBLING
		case IDENT, NUMERIC:
			// Here we have a node name or a bootstrap value
			if prevTok == CLOSEPAR {
				// Bootstrap support value (numeric)
				if tok == NUMERIC {
					if *level == 0 {
						return -1, errors.New("Newick Error: We do not accept support value on root")
					}
					e, err := newNode.ParentEdge()
					if err != nil {
						return -1, err
					}
					support, errf := strconv.ParseFloat(lit, 64)
					if errf != nil {
						return -1, err
					}
					e.SetSupport(support)
				} else {
					// Node name
					if newNode == nil {
						return -1, errors.New("Newick Error: Cannot assign node name to nil node: " + lit)
					}
					// If of the form numeric/numeric => then Support value/pvalue
					vals := strings.Split(lit, "/")
					hasname := true
					if len(vals) == 2 {
						if supp, errf := strconv.ParseFloat(vals[0], 64); errf == nil {
							e, err := newNode.ParentEdge()
							if err != nil {
								return -1, err
							}
							if pval, errf := strconv.ParseFloat(vals[1], 64); errf == nil {
								e.SetSupport(supp)
								e.SetPValue(pval)
								hasname = false
							}
						}
					}
					if hasname {
						newNode.SetName(lit)
					}
				}
			} else {
				// Else we have a new tip
				if prevTok != -1 && prevTok != NEWSIBLING {
					return -1, errors.New("Newick Error: There should not be a tip name in this context: " + lit)
				}
				if node == nil {
					return -1, errors.New("Cannot create a new tip with no parent: " + lit)
				}
				newNode = t.NewNode()
				newNode.SetName(lit)
				t.ConnectNodes(node, newNode)
				prevTok = tok
			}
		case EOT:
			p.unscan()
			if (*level) != 0 {
				return -1, errors.New("Newick Error: Mismatched parenthesis")
			}
			return tok, nil
		case EOF:
			return tok, nil
		}
	}
}
