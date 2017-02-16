package download

const (
	IMGFORMAT_SVG     = 0
	IMGFORMAT_PNG     = 1
	IMGFORMAT_EPS     = 2
	IMGFORMAT_PDF     = 3
	IMGFORMAT_UNKNOWN = 4
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
	default:
		return IMGFORMAT_UNKNOWN
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
	default:
		return "unknown"
	}
}
