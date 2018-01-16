package nexus

type Token int64

var eof = rune(0)

const (
	ILLEGAL Token = iota
	EOF
	WS
	IDENT        // Name of Node, or comment, or keyword
	NUMERIC      // Any numerical value
	OPENBRACK    // [ : For comment
	CLOSEBRACK   // ] : For comment
	ENDOFCOMMAND // ; : End of command
	ENDOFLINE    // \r \n
	COMMA        // , : separator for tables (translation command for ex)

	// Keywords
	NEXUS     // #NEXUS : Start of nexus file
	EQUAL     // '=' between keyword and value
	BEGIN     // Begin
	DATA      // Begin data -> Alignment
	TAXA      // Begin taxa -> Definition of taxa
	TAXLABELS // Begin taxa : list of  taxlabels
	TREES     // Begin trees -> Definition of trees
	TREE      // A specific tree in the BEGIN TREES section
	TRANSLATE // Command that defines a translation table for taxa names

	DIMENSIONS // Dimensions
	NTAX       // Dimensions : Number of taxa
	NCHAR      // Dimensions : Length of alignment

	FORMAT   // Format
	DATATYPE // Format datatype=dna
	MISSING  // Format missing=?  missing char
	GAP      // Format gap=- gap character

	MATRIX // Matrix
	END    // End
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isIdent(ch rune) bool {
	return ch != '[' && ch != ']' && ch != ';' && ch != '=' && ch != '\r' && ch != '\n' && ch != ',' && !isWhitespace(ch)
}

func isEndOfLine(ch rune) bool {
	return ch == '\n' || ch == '\r'
}

func isCR(ch rune) bool {
	return ch == '\r'
}

func isNL(ch rune) bool {
	return ch == '\n'
}
