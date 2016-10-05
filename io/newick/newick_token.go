package newick

type Token int64

var eof = rune(0)

const (
	ILLEGAL Token = iota
	EOF
	WS
	IDENT      // Name of Node, or comment
	NUMERIC    // Any numerical value
	OPENPAR    // (
	CLOSEPAR   // )
	STARTLEN   // :
	OPENBRACK  // [ : For comment
	CLOSEBRACK // ] : For comment
	NEWSIBLING // ,
	EOT        // ;
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isIdent(ch rune) bool {
	return ch != '[' && ch != ']' &&
		ch != '(' && ch != ')' &&
		ch != ',' && ch != ':' && ch != ';'
}
