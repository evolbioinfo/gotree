package newick

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/evolbioinfo/gotree/tree"
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
// ignoreSemiColumn allows to parse identifiers that contain ";"
// such as comments [...;...]
func (p *Parser) scan(ignoreSemiColumn bool) (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan(ignoreSemiColumn)

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan(false)
	if tok == WS {
		tok, lit = p.scan(false)
	}
	return
}

// Parses a Newick String.
func (p *Parser) Parse() (newtree *tree.Tree, err error) {
	// May have information inside [] before the tree
	tok, lit := p.scanIgnoreWhitespace()
	if tok == OPENBRACK || tok == LABEL {
		if _, err = p.consumeComment(tok, lit); err != nil {
			return
		}
		// Next token should be a "OPENPAR" token.
		tok, lit = p.scanIgnoreWhitespace()
	}

	//Next token should be a "OPENPAR" token.
	if tok != OPENPAR {
		err = fmt.Errorf("found %q, expected (", lit)
		return
	}
	p.unscan()
	newtree = tree.NewTree()

	// Now we can parse recursively the tree
	// Read a field.
	level := 0
	if _, err = p.parseIter(newtree, &level); err != nil {
		return
	}
	if level != 0 {
		err = fmt.Errorf("newick error : mismatched parenthesis after parsing")
		return
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != EOT {
		err = fmt.Errorf("found %q, expected ;", lit)
		return
	}
	/* Remove spaces before and after tip names */
	for _, tip := range newtree.Tips() {
		tip.SetName(strings.TrimSpace(tip.Name()))
	}

	//newtree.ReinitIndexes()
	//tree.UpdateTipIndex()
	// err = tree.ClearBitSets()
	// if err != nil {
	// 	return nil, err
	// }
	//tree.UpdateBitSet()
	// Not necessary at the parsing step...
	// may be too long to do each time
	//tree.ComputeDepths()
	// Return the successfully parsed statement.
	return
}

func (p *Parser) parseIter(t *tree.Tree, level *int) (prevTok Token, err error) {
	var nodeStack *NodeStack = NewNodestack()
	var newNode, node *tree.Node = nil, nil
	var edge *tree.Edge = nil
	var length, support, pval float64
	prevTok = -1
	var nedges, nnodes int = 0, 0
	defer nodeStack.Clear()

	for {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case OPENPAR:
			if node == nil {
				if *level > 0 {
					err = errors.New("nil node at depth > 0")
					return
				}
				node = t.NewNode()
				node.SetId(nnodes)
				nnodes++
				nodeStack.Push(node, nil)
				t.SetRoot(node)
			} else {
				if *level == 0 {
					err = errors.New("newick Error: An open parenthesis while the stack is empty... Forgot a ';' at the end of previous tree?")
					return
				}
				newNode = t.NewNode()
				newNode.SetId(nnodes)
				nnodes++
				edge = t.ConnectNodes(node, newNode)
				edge.SetId(nedges)
				nedges++
				node = newNode
				nodeStack.Push(node, edge)
			}
			(*level)++
			prevTok = tok
		case CLOSEPAR:
			prevTok = tok
			(*level)--
			if _, _, err = nodeStack.Pop(); err != nil {
				err = errors.New("newick Error: Closing parenthesis while the stack is already empty")
				return
			}
			node, edge, _ = nodeStack.Head()
		case OPENBRACK, LABEL:
			var comment string
			//if prevTok == OPENPAR || prevTok == NEWSIBLING || prevTok == -1 {
			if comment, err = p.consumeComment(tok, lit); err != nil {
				return
			}
			// Add comment to edge if comment located after branch length
			if prevTok == STARTLEN && edge != nil {
				edge.AddComment(comment)
			} else if prevTok == STARTLEN && edge == nil && node != nil {
				node.AddComment(comment)
			} else if (prevTok == CLOSEPAR || prevTok == IDENT || prevTok == NUMERIC || prevTok == CLOSEBRACK) && node != nil {
				// Else we add comment to node
				node.AddComment(comment)
			} else {
				err = errors.New("newick error: comment should not be located here: " + lit)
				return
			}
			prevTok = CLOSEBRACK
		case CLOSEBRACK:
			// Error here should not have
			err = errors.New("newick error: mismatched ] here")
			return
		case STARTLEN:
			if tok, lit = p.scanIgnoreWhitespace(); tok != NUMERIC {
				err = errors.New("newick error: no numeric value after ':'")
				return
			}
			// We skip length if the length is assigned to the root node
			if node != nil && *level != 0 {
				if edge == nil {
					err = errors.New("Newick Error: Edge length should not be located here: " + lit)
					return
				}
				if edge.Length() != tree.NIL_LENGTH {
					err = errors.New("Newick Error: More than one length is given :" + lit)
					return
				}
				if length, err = strconv.ParseFloat(lit, 64); err != nil {
					err = errors.New("Newick Error: Length is not a float value : " + lit)
					return
				}
				edge.SetLength(length)
			} else if *level == 0 {
				log.Print("Newick : Branch lengths attached to root node are ignored")
			} else {
				// For root node, level==0, we just ignore it
				err = errors.New("Newick Error: Cannot assign length to nil node :" + lit)
				return
			}
			prevTok = STARTLEN
		case NEWSIBLING:
			if _, _, err = nodeStack.Pop(); err != nil {
				err = errors.New("Newick Error: Stack is empty, a coma should not be located here: " + lit)
				return
			}
			node, edge, _ = nodeStack.Head()
			prevTok = NEWSIBLING
		case IDENT, NUMERIC:
			// Here we have a node name or a bootstrap value
			if prevTok == CLOSEPAR {
				// Bootstrap support value (numeric)
				if tok == NUMERIC {
					if *level == 0 || edge == nil {
						log.Print("Newick : Support values attached to root node are ignored")
						//return -1, errors.New("Newick Error: We do not accept support value on root")
					} else {
						if support, err = strconv.ParseFloat(lit, 64); err != nil {
							return
						}
						edge.SetSupport(support)
					}
				} else {
					// If of the form numeric/numeric => then Support value/pvalue
					vals := strings.Split(lit, "/")
					hasname := true
					if len(vals) == 2 && edge != nil {
						if support, err = strconv.ParseFloat(vals[0], 64); err == nil {
							if pval, err = strconv.ParseFloat(vals[1], 64); err == nil {
								edge.SetSupport(support)
								edge.SetPValue(pval)
								hasname = false
							}
						}
					}
					if hasname {
						// Node name
						if node == nil {
							err = errors.New("Newick Error: Cannot assign node name to nil node: " + lit)
							return
						}
						node.SetName(lit)
					}
				}
			} else {
				// Else we have a new tip
				if prevTok != OPENPAR && prevTok != NEWSIBLING {
					err = errors.New("Newick Error: There should not be a tip name in this context: " + lit)
					return
				}
				if node == nil {
					err = errors.New("Cannot create a new tip with no parent: " + lit)
					return
				}
				newNode = t.NewNode()
				newNode.SetId(nnodes)
				nnodes++
				newNode.SetName(lit)
				edge = t.ConnectNodes(node, newNode)
				edge.SetId(nedges)
				nedges++
				node = newNode
				nodeStack.Push(node, edge)
				prevTok = tok
			}
		case EOT:
			p.unscan()
			if (*level) != 0 {
				err = errors.New("newick Error: Mismatched parenthesis at ;")
				return
			}
			prevTok = tok
			return
		case EOF:
			prevTok = tok
			return
		}
	}
}

// Consumes comment inside brackets [comment] if the given current token is a [.
// At the end returns the matching ] token and lit.
// If the given token is not a [, then returns an error
func (p *Parser) consumeComment(curtoken Token, curlit string) (comment string, err error) {
	if curtoken == OPENBRACK || curtoken == LABEL {
		commenttoken, commentlit := p.scan(true)
		for (curtoken == LABEL && commenttoken != LABEL) || (curtoken == OPENBRACK && commenttoken != CLOSEBRACK) {
			if commenttoken == EOF || commenttoken == ILLEGAL {
				err = fmt.Errorf("unmatched bracket: %s (%s)", comment, commentlit)
				return
			} else {
				comment += commentlit
			}
			commenttoken, commentlit = p.scan(true)
		}
	} else {
		err = fmt.Errorf("a comment must start with [")
	}
	return
}
