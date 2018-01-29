package download

const (
	IMGFORMAT_SVG = iota
	IMGFORMAT_PNG
	IMGFORMAT_EPS
	IMGFORMAT_PDF
	TXTFORMAT_NEWICK
	TXTFORMAT_NEXUS
	TXTFORMAT_PHYLOXML
	FORMAT_UNKNOWN
)

func Format(format string) int {
	switch format {
	case "svg":
		return IMGFORMAT_SVG
	case "png":
		return IMGFORMAT_PNG
	case "eps":
		return IMGFORMAT_EPS
	case "pdf":
		return IMGFORMAT_PDF
	case "newick":
		return TXTFORMAT_NEWICK
	case "nexus":
		return TXTFORMAT_NEXUS
	case "phyloxml":
		return TXTFORMAT_PHYLOXML
	default:
		return FORMAT_UNKNOWN
	}
}

func StrFormat(format int) string {
	switch format {
	case IMGFORMAT_SVG:
		return "svg"
	case IMGFORMAT_PNG:
		return "png"
	case IMGFORMAT_EPS:
		return "eps"
	case IMGFORMAT_PDF:
		return "pdf"
	case TXTFORMAT_NEWICK:
		return "newick"
	case TXTFORMAT_NEXUS:
		return "nexus"
	case TXTFORMAT_PHYLOXML:
		return "phyloxml"
	default:
		return "unknown"
	}
}
