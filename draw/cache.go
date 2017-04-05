package draw

import (
	"math"
)

/* Cache for lines and points to draw the tree */
type layoutCache struct {
	tipLabelPoints []*layoutPoint
	branchPaths    []*layoutLine
	nodePoints     []*layoutPoint
	curvePaths     []*layoutCurve
}

type layoutPoint struct {
	x       float64
	y       float64
	brAngle float64 // Angle of the incoming branch
	name    string
}

type layoutLine struct {
	p1 *layoutPoint
	p2 *layoutPoint
}

type layoutCurve struct {
	center      *layoutPoint // center of the circle
	middlepoint *layoutPoint // point on the circle at the middle of the curve
	radius      float64      // radius of the circle
	startAngle  float64
	endAngle    float64
}

func newLayoutCache() *layoutCache {
	return &layoutCache{
		make([]*layoutPoint, 0),
		make([]*layoutLine, 0),
		make([]*layoutPoint, 0),
		make([]*layoutCurve, 0),
	}
}

func (cache *layoutCache) borders() (xmin, ymin, xmax, ymax float64) {
	xmin, ymin = 100000.0, 100000.0
	xmax, ymax = 0.0, 0.0
	for _, line := range cache.branchPaths {
		xmin = math.Min(math.Min(xmin, line.p1.x), line.p2.x)
		ymin = math.Min(math.Min(ymin, line.p1.y), line.p2.y)
		xmax = math.Max(math.Max(xmax, line.p1.x), line.p2.x)
		ymax = math.Max(math.Max(ymax, line.p1.y), line.p2.y)
	}
	return
}
