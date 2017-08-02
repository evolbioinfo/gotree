package nexus

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fredericlemoine/goalign/align"
	"github.com/fredericlemoine/gotree/io/newick"
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

// Parses a Newick String.
func (p *Parser) Parse() (*Nexus, error) {
	var nchar, ntax int64
	datatype := "dna"
	missing := '*'
	gap := '-'
	var taxlabels map[string]bool = nil
	var names, sequences, treestrings, treenames []string
	nexus := NewNexus()

	// First token should be a "NEXUS" token.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != NEXUS {
		return nil, fmt.Errorf("found %q, expected #NEXUS", lit)
	}

	// Now we can parse the remaining of the file
	for {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == ILLEGAL {
			return nil, fmt.Errorf("found illegal token %q", lit)
		}
		if tok == EOF {
			break
		}
		if tok == ENDOFLINE {
			continue
		}
		if tok == BEGIN {
			tok2, _ := p.scanIgnoreWhitespace()
			tok3, _ := p.scanIgnoreWhitespace()
			if tok3 != ENDOFCOMMAND {
				return nil, fmt.Errorf("found %q, expected ;", lit)
			}
			var err error
			switch tok2 {
			case TAXA:
				taxlabels, err = p.parseTaxa()
			case TREES:
				treenames, treestrings, err = p.parseTrees()
			case DATA:
				names, sequences, nchar, ntax, datatype, missing, gap, err = p.parseData()
			}
			if err != nil {
				return nil, err
			}
		}
	}
	if gap != '-' || missing != '*' {
		return nil, fmt.Errorf("We only accept - gaps (not %c) && * missing (not %c) so far", gap, missing)
	}

	// We initialize alignment structure
	if names != nil && sequences != nil {
		al := align.NewAlign(align.AlphabetFromString(datatype))
		if al.Alphabet() == align.UNKNOWN {
			return nil, fmt.Errorf("Unknown datatype: %q", datatype)
		}
		if len(names) != int(ntax) {
			return nil, fmt.Errorf("Number of taxa in alignment (%d)  does not corresponds to definition %d", len(names), ntax)
		}
		for i, seq := range sequences {
			if len(seq) != int(nchar) {
				return nil, fmt.Errorf("Number of character in sequence #%d (%d) does not corresponds to definition %d", i, len(seq), nchar)
			}
			if err := al.AddSequence(names[i], seq, ""); err != nil {
				return nil, err
			}
		}
		// We check that tax labels are the same as alignment sequence names
		if taxlabels != nil {
			var err error
			al.Iterate(func(name string, sequence string) {
				if _, ok := taxlabels[name]; !ok {
					err = fmt.Errorf("Sequence name %s in the alignment is not defined in the TAXLABELS block", name)
				}
			})
			if err != nil {
				return nil, err
			}
			if al.NbSequences() != len(taxlabels) {
				return nil, fmt.Errorf("Some taxa names defined in TAXLABELS are not present in the alignment")
			}
		}

		nexus.SetAlignment(al)
	}
	// We initialize tree structures
	if treenames != nil && treestrings != nil {
		for i, treestr := range treestrings {
			t, err := newick.NewParser(strings.NewReader(treestr + ";")).Parse()
			if err != nil {
				return nil, err
			}
			// We check that tax labels are the same as tree taxa
			if taxlabels != nil {
				tips := t.Tips()
				for _, tip := range tips {
					if _, ok := taxlabels[tip.Name()]; !ok {
						return nil, fmt.Errorf("Taxa name %s in the tree %d is not defined in the TAXLABELS block", i, tip.Name())
					}
				}
				if len(tips) != len(taxlabels) {
					return nil, fmt.Errorf("Some tax names defined in TAXLABELS are not present in the tree %d", i)
				}
			}
			nexus.AddTree(treenames[i], t)
		}
	}
	return nexus, nil
}

func (p *Parser) parseTaxa() (map[string]bool, error) {
	taxlabels := make(map[string]bool)
	var err error
	stoptaxa := false
	for !stoptaxa {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case ENDOFLINE:
			continue
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", lit)
			stoptaxa = true
		case EOF:
			err = fmt.Errorf("End of file within a TAXA block (no END;)")
			stoptaxa = true
		case END:
			tok2, _ := p.scanIgnoreWhitespace()
			if tok2 != ENDOFCOMMAND {
				err = fmt.Errorf("End token without ;")
			}
			stoptaxa = true
		case TAXLABELS:
			stoplabels := false
			for !stoplabels {
				tok2, lit2 := p.scanIgnoreWhitespace()
				switch tok2 {
				case ENDOFCOMMAND:
					stoplabels = true
				case IDENT:
					taxlabels[lit2] = true
				default:
					err = fmt.Errorf("Unknown token %q in taxlabel list", lit2)
					stoplabels = true
				}
			}
			if err != nil {
				stoptaxa = true
			}
		default:
			err = fmt.Errorf("Unknown token %q", lit)
			stoptaxa = true
		}
	}
	return taxlabels, err
}

func (p *Parser) parseTrees() (treenames, treestrings []string, err error) {
	treenames = make([]string, 0)
	treestrings = make([]string, 0)
	stoptrees := false
	for !stoptrees {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case ENDOFLINE:
			continue
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", lit)
			stoptrees = true
		case EOF:
			err = fmt.Errorf("End of file within a TREES block (no END;)")
			stoptrees = true
		case END:
			tok2, _ := p.scanIgnoreWhitespace()
			if tok2 != ENDOFCOMMAND {
				err = fmt.Errorf("End token without ;")
			}
			stoptrees = true
		case TREE:
			tok2, lit2 := p.scanIgnoreWhitespace()
			if tok2 != IDENT {
				err = fmt.Errorf("Expecting a tree name after TREE, got %q", lit2)
				stoptrees = true
			}
			tok3, lit3 := p.scanIgnoreWhitespace()
			if tok3 != EQUAL {
				err = fmt.Errorf("Expecting '=' after tree name, got %q", lit3)
				stoptrees = true
			}
			tok4, lit4 := p.scanIgnoreWhitespace()
			if tok4 != IDENT {
				err = fmt.Errorf("Expecting a tree after 'TREE name =', got  %q", lit4)
				stoptrees = true
			}
			tok5, lit5 := p.scanIgnoreWhitespace()
			if tok5 != ENDOFCOMMAND {
				err = fmt.Errorf("Expecting ';' after 'TREE name = tree', got %q", lit5)
				stoptrees = true
			}
			treenames = append(treenames, lit2)
			treestrings = append(treestrings, lit4)
		default:
			err = fmt.Errorf("Unknown token %q", lit)
			stoptrees = true
		}
	}
	return
}

func (p *Parser) parseData() (names, sequences []string, nchar, ntax int64, datatype string, missing, gap rune, err error) {
	datatype = "dna"
	missing = '*'
	gap = '-'
	stopdata := false
	sequences = make([]string, 0)
	names = make([]string, 0)
	for !stopdata {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case ENDOFLINE:
			break
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", lit)
			stopdata = true
		case EOF:
			err = fmt.Errorf("End of file within a TAXA block (no END;)")
			stopdata = true
		case END:
			tok2, _ := p.scanIgnoreWhitespace()
			if tok2 != ENDOFCOMMAND {
				err = fmt.Errorf("End token without ;")
			}
			stopdata = true
		case DIMENSIONS:
			stopdimensions := false
			for !stopdimensions {
				tok2, lit2 := p.scanIgnoreWhitespace()
				switch tok2 {
				case ENDOFCOMMAND:
					stopdimensions = true
				case NTAX:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 != EQUAL {
						err = fmt.Errorf("Expecting '=' after NTAX, got %q", lit3)
						stopdimensions = true
					}
					tok4, lit4 := p.scanIgnoreWhitespace()
					if tok4 != NUMERIC {
						err = fmt.Errorf("Expecting Integer value after 'NTAX=', got %q", lit4)
						stopdimensions = true
					}
					ntax, err = strconv.ParseInt(lit4, 10, 64)
					if err != nil {
						stopdimensions = true
					}
				case NCHAR:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 != EQUAL {
						err = fmt.Errorf("Expecting '=' after NTAX, got %q", lit3)
						stopdimensions = true
					}
					tok4, lit4 := p.scanIgnoreWhitespace()
					if tok4 != NUMERIC {
						err = fmt.Errorf("Expecting Integer value after 'NTAX=', got %q", lit4)
						stopdimensions = true
					}
					nchar, err = strconv.ParseInt(lit4, 10, 64)
					if err != nil {
						stopdimensions = true
					}
				default:
					err = fmt.Errorf("Unknown token %q in taxlabel list", lit2)
				}
				if err != nil {
					stopdata = true
				}
			}
		case FORMAT:
			stopformat := false
			for !stopformat {
				tok2, _ := p.scanIgnoreWhitespace()

				switch tok2 {
				case ENDOFCOMMAND:
					stopformat = true
				case DATATYPE:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 != EQUAL {
						err = fmt.Errorf("Expecting '=' after DATATYPE, got %q", lit3)
						stopformat = true
					}
					tok4, lit4 := p.scanIgnoreWhitespace()
					if tok4 == IDENT {
						datatype = lit4
					} else {
						err = fmt.Errorf("Expecting identifier after 'DATATYPE=', got %q", lit4)
						stopformat = true
					}
				case MISSING:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 != EQUAL {
						err = fmt.Errorf("Expecting '=' after MISSING, got %q", lit3)
						stopformat = true
					}
					tok4, lit4 := p.scanIgnoreWhitespace()
					if tok4 != IDENT {
						err = fmt.Errorf("Expecting Integer value after 'MISSING=', got %q", lit4)
						stopformat = true
					}
					if len(lit4) != 1 {
						err = fmt.Errorf("Expecting a single character after MISSING=', got %q", lit4)
						stopformat = true
					}
					missing = []rune(lit4)[0]
				case GAP:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 != EQUAL {
						err = fmt.Errorf("Expecting '=' after GAP, got %q", lit3)
						stopformat = true
					}
					tok4, lit4 := p.scanIgnoreWhitespace()
					if tok4 != IDENT {
						err = fmt.Errorf("Expecting an identifier after 'GAP=', got %q", lit4)
						stopformat = true
					}
					if len(lit4) != 1 {
						err = fmt.Errorf("Expecting a single character after GAP=', got %q", lit4)
						stopformat = true
					}
					gap = []rune(lit4)[0]
				default:
					err = fmt.Errorf("Unknown token %q in taxlabel list", lit)
				}
				if err != nil {
					stopdata = true
				}
			}
		case MATRIX:
			stopmatrix := false
			for !stopmatrix {
				tok2, lit2 := p.scanIgnoreWhitespace()
				switch tok2 {
				case IDENT:
					tok3, lit3 := p.scanIgnoreWhitespace()
					if tok3 == IDENT {
						names = append(names, lit2)
						sequences = append(sequences, lit3)
					} else {
						err = fmt.Errorf("Expecting sequence after sequence identifier (%q) in Matrix block, got %q", lit2, lit3)
						stopmatrix = true
					}
				case ENDOFLINE:
					break
				case ENDOFCOMMAND:
					stopmatrix = true
				default:
					err = fmt.Errorf("Expecting sequence identifier in Matrix block, got %q", lit2)
					stopmatrix = true
				}
			}
			if err != nil {
				stopdata = true
			}

		default:
			err = fmt.Errorf("Unknown token %q", lit)
			stopdata = true
		}
	}
	return
}
