package poly2tri

import "math"

const (
	PI_3div4 = 3 * math.Pi / 4
	PI_div2  = 1.57079632679489661923
	EPSILON  = 1e-12
)

type Orientation int

const (
	CW Orientation = iota
	CCW
	COLLINEAR
)

/**
 * Forumla to calculate signed area<br>
 * Positive if CCW<br>
 * Negative if CW<br>
 * 0 if collinear<br>
 * <pre>
 * A[P1,P2,P3]  =  (x1*y2 - y1*x2) + (x2*y3 - y2*x3) + (x3*y1 - y3*x1)
 *              =  (x1-x3)*(y2-y3) - (y1-y3)*(x2-x3)
 * </pre>
 */
func Orient2d(pa, pb, pc *Point) Orientation {
	detleft := (pa.x - pc.x) * (pb.y - pc.y)
	detright := (pa.y - pc.y) * (pb.x - pc.x)
	val := detleft - detright
	if val > -EPSILON && val < EPSILON {
		return COLLINEAR
	} else if val > 0 {
		return CCW
	}
	return CW
}

func InScanArea(pa, pb, pc, pd *Point) bool {
	oadb := (pa.x-pb.x)*(pd.y-pb.y) - (pd.x-pb.x)*(pa.y-pb.y)
	if oadb >= -EPSILON {
		return false
	}

	oadc := (pa.x-pc.x)*(pd.y-pc.y) - (pd.x-pc.x)*(pa.y-pc.y)
	if oadc <= EPSILON {
		return false
	}
	return true
}
