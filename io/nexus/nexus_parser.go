package nexus

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/evolbioinfo/goalign/align"
	treeio "github.com/evolbioinfo/gotree/io"
	"github.com/evolbioinfo/gotree/io/newick"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
	translationTable map[string]string // For taxa name translation
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

// scanIgnoreWhitespace scans the next non-whitespace and end of line token.
func (p *Parser) scanIgnoreWhitespaceAndEOL() (tok Token, lit string) {
	tok, lit = p.scan()
	for tok == WS || tok == ENDOFLINE {
		tok, lit = p.scan()
	}
	return
}

// Parses Nexus content from the reader
func (p *Parser) Parse() (*Nexus, error) {
	var nchar, ntax, taxantax int64
	datatype := "dna"
	missing := '*'
	gap := '-'
	var taxlabels map[string]bool = nil
	var names, treestrings, treenames []string
	var sequences map[string]string

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

		if tok == OPENBRACK {
			var err error
			if tok, lit, err = p.consumeComment(tok, lit); err != nil {
				return nil, err
			}
			tok, lit = p.scanIgnoreWhitespace()
		}

		// Beginning of a block
		if tok == BEGIN {
			// Next token should be the name of the block
			tok2, lit2 := p.scanIgnoreWhitespace()
			// Then a ;
			tok3, lit3 := p.scanIgnoreWhitespace()
			if tok3 != ENDOFCOMMAND {
				return nil, fmt.Errorf("found %q, expected ;", lit3)
			}
			var err error
			switch tok2 {
			case TAXA:
				// TAXA BLOCK
				taxantax, taxlabels, err = p.parseTaxa()
			case TREES:
				// TREES BLOCK
				treenames, treestrings, err = p.parseTrees()
			case DATA:
				// DATA/CHARACTERS BLOCK
				names, sequences, nchar, ntax, datatype, missing, gap, err = p.parseData()
			default:
				// If an unsupported block is seen, we just skip it
				treeio.LogWarning(fmt.Errorf("Unsupported block %q, skipping", lit2))
				err = p.parseUnsupportedBlock()
			}

			if err != nil {
				return nil, err
			}
		}
	}

	if int(taxantax) != -1 && int(taxantax) != len(taxlabels) {
		return nil, fmt.Errorf("Number of defined taxa in TAXLABELS/DIMENSIONS (%d) is different from length of taxa list (%d)", taxantax, len(taxlabels))
	}

	if gap != '-' || missing != '*' {
		return nil, fmt.Errorf("We only accept - gaps (not %c) && * missing (not %c) so far", gap, missing)
	}

	// We initialize alignment structure using goalign structure
	if names != nil && sequences != nil {
		al := align.NewAlign(align.AlphabetFromString(datatype))
		if al.Alphabet() == align.UNKNOWN {
			return nil, fmt.Errorf("Unknown datatype: %q", datatype)
		}
		if len(names) != int(ntax) && ntax != -1 {
			return nil, fmt.Errorf("Number of taxa in alignment (%d)  does not correspond to definition %d", len(names), ntax)
		}
		for i, name := range names {
			seq, _ := sequences[name]
			if len(seq) != int(nchar) && nchar != -1 {
				return nil, fmt.Errorf("Number of character in sequence #%d (%d) does not correspond to definition %d", i, len(seq), nchar)
			}
			if err := al.AddSequence(name, seq, ""); err != nil {
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
	// We initialize tree structures using gotree structure
	if treenames != nil && treestrings != nil {
		for i, treestr := range treestrings {
			t, err := newick.NewParser(strings.NewReader(treestr + ";")).Parse()
			if err != nil {
				return nil, err
			}
			// We translate taxa labels if needed
			if p.translationTable != nil {
				if err2 := t.Rename(p.translationTable); err2 != nil {
					return nil, err2
				}
			}
			// We check that tax labels are the same as tree taxa
			if taxlabels != nil {
				tips := t.Tips()
				for _, tip := range tips {
					if _, ok := taxlabels[tip.Name()]; !ok {
						return nil, fmt.Errorf("Taxa name %s in the tree %d is not defined in the TAXLABELS block", tip.Name(), i)
					}
				}
				if len(tips) != len(taxlabels) {
					return nil, fmt.Errorf("Some tax names defined in TAXLABELS are not present in the tree %d", i)
				}
			}
			//t.ReinitIndexes()
			nexus.AddTree(treenames[i], t)
		}
	}
	return nexus, nil
}

// Parse taxa block
func (p *Parser) parseTaxa() (int64, map[string]bool, error) {
	taxlabels := make(map[string]bool)
	var err error
	stoptaxa := false
	var ntax int64 = -1
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
		case DIMENSIONS:
			// Dimensions of the data: ntax
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
				default:
					if err = p.parseUnsupportedKey(lit2); err != nil {
						stopdimensions = true
					}
					treeio.LogWarning(fmt.Errorf("Unsupported key %q in %q command, skipping", lit2, lit))
				}
				if err != nil {
					stoptaxa = true
				}
			}
		case TAXLABELS:
			stoplabels := false
			for !stoplabels {
				tok2, lit2 := p.scanIgnoreWhitespace()
				switch tok2 {
				case ENDOFLINE:
				case ENDOFCOMMAND:
					stoplabels = true
				case IDENT, NUMERIC:
					taxlabels[lit2] = true
				default:
					err = fmt.Errorf("Unknown token %q in taxlabel list", lit2)
					stoplabels = true
				}
			}
			if err != nil {
				stoptaxa = true
			}
		case OPENBRACK:
			if tok, lit, err = p.consumeComment(tok, lit); err != nil {
				stoptaxa = true
			}

		default:
			err = p.parseUnsupportedCommand()
			treeio.LogWarning(fmt.Errorf("Unsupported command %q in block TAXA, skipping", lit))
			if err != nil {
				stoptaxa = true
			}
		}
	}
	return ntax, taxlabels, err
}

// Parse TREES block
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
		case TRANSLATE:
			// TRANSLATE BLOCK
			if p.translationTable, err = p.parseTranslationTable(); err != nil {
				stoptrees = true
				break
			}
		case TREE:
			// A new tree is seen
			tok2, lit2 := p.scanIgnoreWhitespace()
			if tok2 != IDENT && tok2 != NUMERIC {
				err = fmt.Errorf("Expecting a tree name after TREE, got %q", lit2)
				stoptrees = true
				break
			}
			tok3, lit3 := p.scanIgnoreWhitespace()
			if tok3 != EQUAL {
				err = fmt.Errorf("Expecting '=' after tree name, got %q", lit3)
				stoptrees = true
				break
			}
			tok4, lit4 := p.scanIgnoreWhitespace()
			if tok4 == OPENBRACK {
				if tok4, lit4, err = p.consumeComment(tok4, lit4); err != nil {
					stoptrees = true
					break
				}
				tok4, lit4 = p.scanIgnoreWhitespaceAndEOL()
			}
			tree := ""
			// We remove whitespaces in the tree string if any,
			// and keep comments in brackets as part of the newick string
			for tok4 != ENDOFCOMMAND {
				if tok4 != IDENT && tok4 != OPENBRACK && tok4 != CLOSEBRACK && tok4 != COMMA && tok4 != EQUAL && tok4 != NUMERIC {
					err = fmt.Errorf("Expecting a tree after 'TREE name =', got  %q", lit4)
					stoptrees = true
					break
				}
				tree += lit4
				tok4, lit4 = p.scanIgnoreWhitespace()
			}
			if tok4 != ENDOFCOMMAND {
				err = fmt.Errorf("Expecting ';' after 'TREE name = tree', got %q", lit4)
				stoptrees = true
				break
			}
			treenames = append(treenames, lit2)
			treestrings = append(treestrings, tree)
		case OPENBRACK:
			if tok, lit, err = p.consumeComment(tok, lit); err != nil {
				stoptrees = true
			}
		default:
			err = p.parseUnsupportedCommand()
			treeio.LogWarning(fmt.Errorf("Unsupported command %q in block TREES, skipping", lit))
			if err != nil {
				stoptrees = true
			}
		}
	}
	return
}

// Parse TREES block
func (p *Parser) parseTranslationTable() (translationTable map[string]string, err error) {
	translationTable = make(map[string]string)
	stop := false
	for !stop {
		tok, key := p.scanIgnoreWhitespace()
		switch tok {
		case IDENT, NUMERIC:
			if tok2, value := p.scanIgnoreWhitespace(); tok2 != IDENT && tok2 != NUMERIC {
				err = fmt.Errorf("TRANSLATE block: Expecting value name here, and found '%q'", key)
				stop = true
			} else {
				if tok3, end := p.scanIgnoreWhitespace(); tok3 != COMMA && tok3 != ENDOFCOMMAND && tok3 != ENDOFLINE {
					err = fmt.Errorf("TRANSLATE block: Expecting , or ; after key value, and found '%q'", end)
					stop = true
				} else {
					translationTable[key] = value
					stop = (end == ";")
				}
			}
		case ENDOFLINE:
			continue
		case COMMA:
			continue
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", key)
			stop = true
		case EOF:
			err = fmt.Errorf("End of file within a TRANSLATE block (no END;)")
			stop = true
		case ENDOFCOMMAND:
			stop = true
		case OPENBRACK:
			if tok, key, err = p.consumeComment(tok, key); err != nil {
				stop = true
			}
		default:
			err = fmt.Errorf("Unsupported token %q in TRANSLATE command, skipping", key)
			stop = true
		}
	}
	return
}

// DATA / Characters BLOCK
func (p *Parser) parseData() (names []string, sequences map[string]string, nchar, ntax int64, datatype string, missing, gap rune, err error) {
	datatype = "dna"
	missing = '*'
	gap = '-'
	stopdata := false
	sequences = make(map[string]string)
	names = make([]string, 0)
	nchar = -1
	ntax = -1
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
			// Dimensions of the data: nchar , ntax
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
					if err = p.parseUnsupportedKey(lit2); err != nil {
						stopdimensions = true
					}
					treeio.LogWarning(fmt.Errorf("Unsupported key %q in %q command, skipping", lit2, lit))
				}
				if err != nil {
					stopdata = true
				}
			}
		case FORMAT:
			// Format of the data bock: datatype, missing, gap
			stopformat := false
			for !stopformat {
				tok2, lit2 := p.scanIgnoreWhitespace()

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
					if err = p.parseUnsupportedKey(lit2); err != nil {
						stopformat = true
					}
					treeio.LogWarning(fmt.Errorf("Unsupported key %q in %q command, skipping", lit2, lit))
				}
				if err != nil {
					stopdata = true
				}
			}
		case MATRIX:
			// Character matrix (Alignmemnt)
			// So far: Does not handle interleave case...
			stopmatrix := false
			for !stopmatrix {
				tok2, lit2 := p.scanIgnoreWhitespace()
				switch tok2 {
				case IDENT:
					// We remove whitespaces in sequences if any
					// and take into account possibly interleaved
					// sequences
					stopseq := false
					name := lit2
					sequence := ""
					for !stopseq {
						tok3, lit3 := p.scanIgnoreWhitespace()
						switch tok3 {
						case IDENT:
							sequence = sequence + lit3
						case ENDOFLINE:
							stopseq = true
						default:
							err = fmt.Errorf("Expecting sequence after sequence identifier (%q) in Matrix block, got %q", lit2, lit3)
							stopseq = true
						}
					}
					if err != nil {
						stopmatrix = true
					} else {
						addseq(sequences, &names, sequence, name)
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
		case OPENBRACK:
			if tok, lit, err = p.consumeComment(tok, lit); err != nil {
				stopdata = true
			}
		default:
			err = p.parseUnsupportedCommand()
			treeio.LogWarning(fmt.Errorf("Unsupported command %q in block DATA, skipping", lit))
			if err != nil {
				stopdata = true
			}
		}
	}
	return
}

// Just skip the current command
func (p *Parser) parseUnsupportedCommand() (err error) {
	// Unsupported data command
	stopunsupported := false
	for !stopunsupported {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", lit)
			stopunsupported = true
		case EOF:
			err = fmt.Errorf("End of file within a command (no;)")
			stopunsupported = true
		case ENDOFCOMMAND:
			stopunsupported = true
		}
	}
	return
}

// Just skip the current key
func (p *Parser) parseUnsupportedKey(key string) (err error) {
	// Unsupported token
	tok, lit := p.scanIgnoreWhitespace()
	if tok != EQUAL {
		err = fmt.Errorf("Expecting '=' after %s, got %q", key, lit)
	} else {
		tok2, lit2 := p.scanIgnoreWhitespace()
		if tok2 != IDENT && tok2 != NUMERIC {
			err = fmt.Errorf("Expecting an identifier after '%s=', got %q", key, lit2)
		}
	}
	return
}

// Just skip the current block
func (p *Parser) parseUnsupportedBlock() error {
	var err error
	stopunsupported := false
	for !stopunsupported {
		tok, lit := p.scanIgnoreWhitespace()
		switch tok {
		case ILLEGAL:
			err = fmt.Errorf("found illegal token %q", lit)
			stopunsupported = true
		case EOF:
			err = fmt.Errorf("End of file within a block (no END;)")
			stopunsupported = true
		case END:
			tok2, _ := p.scanIgnoreWhitespace()
			if tok2 != ENDOFCOMMAND {
				err = fmt.Errorf("End token without ;")
			}
			stopunsupported = true
		}
	}
	return err
}

// Append a sequence to the hashmap map[name]seq.
// If the name does not exists in the map, adds the name to names and the sequence to the map (to keep order of the input nexus matrix
// Otherwise, append the sequence to the already insert sequence in the map
func addseq(sequences map[string]string, names *[]string, sequence string, name string) {
	seq, ok := sequences[name]
	if !ok {
		*names = append(*names, name)
		sequences[name] = sequence
	} else {
		sequences[name] = seq + sequence
	}
}

// Consumes comment inside brakets [comment] if the given current token is a [.
// At the end returns the matching ] token and lit.
// If the given token is not a [, then returns the input token and lit
func (p *Parser) consumeComment(curtoken Token, curlit string) (outtoken Token, outlit string, err error) {
	outtoken = curtoken
	curlit = curlit
	if curtoken == OPENBRACK {
		for outtoken != CLOSEBRACK {
			outtoken, outlit = p.scanIgnoreWhitespace()
			if outtoken == EOF || outtoken == ILLEGAL {
				err = fmt.Errorf("Unmatched bracket")
			}
		}
	}
	return
}
