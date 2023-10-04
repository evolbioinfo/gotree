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

// isIdent checks whether the given rune is part of a identifier, such as:
// -tip, node and branch name
// - comments
// - branch length
// - branch support
// If it corresponds to a newick keyword, then returns false
// If ignore semicolumn is true, then ";" is not considered as
// a newick keyword. (useful for parsing comments [...;...])
func isIdent(ch rune, ignoreSemiColumn bool) bool {
	return ch != '[' && ch != ']' &&
		ch != '(' && ch != ')' &&
		ch != ',' && ch != ':' &&
		(ignoreSemiColumn || ch != ';')
}
