package common

import "math"

// -------------------------------------------------------------------------------
/** Compute the signed area of a triangle.
 *  The function accepts an unconstrained template parameter for use with
 *  both aiVector3D and aiVector2D, but generally ignores the third coordinate.*/

func GetArea2D(v1, v2, v3 *AiVector2D) float64 {
	return float64(0.5 * (v1.X*(v3.Y-v2.Y) +
		v2.X*(v1.Y-v3.Y) +
		v3.X*(v2.Y-v1.Y)))
}

const AiEpsilon = 1e-6

// -------------------------------------------------------------------------------
/** Test if a given point p2 is on the left side of the line formed by p0-p1.
 *  The function accepts an unconstrained template parameter for use with
 *  both aiVector3D and aiVector2D, but generally ignores the third coordinate.*/

func OnLeftSideOfLine2D(p0, p1, p2 *AiVector2D) int {
	area := GetArea2D(p0, p2, p1)
	if math.Abs(area) < AiEpsilon {
		return 0
	} else if area > 0 {
		return 1
	} else {
		return -1
	}

}

// -------------------------------------------------------------------------------
/** Test if a given point is inside a given triangle in R2.
 * The function accepts an unconstrained template parameter for use with
 *  both aiVector3D and aiVector2D, but generally ignores the third coordinate.*/

func PointInTriangle2D(p0, p1, p2, pp *AiVector2D) bool {
	// pp should be left side of the three triangle side, by ccw arrow
	c1 := OnLeftSideOfLine2D(p0, p1, pp)
	c2 := OnLeftSideOfLine2D(p1, p2, pp)
	c3 := OnLeftSideOfLine2D(p2, p0, pp)
	return (c1 >= 0) && (c2 >= 0) && (c3 >= 0)
}

// -------------------------------------------------------------------------------
/** Check whether the winding order of a given polygon is counter-clockwise.
 *  The function accepts an unconstrained template parameter, but is intended
 *  to be used only with aiVector2D and aiVector3D (z axis is ignored, only
 *  x and y are taken into account).
 * @note Code taken from http://cgm.cs.mcgill.ca/~godfried/teaching/cg-projects/97/Ian/applet1.html and translated to C++
 */
func SliceAiVector3DToAiVector2D(in []*AiVector3D) (out []*AiVector2D) {
	for _, v := range in {
		out = append(out, &AiVector2D{X: v.X, Y: v.Y})
	}
	return out
}

func SliceAiVector3DToFloatArr(in []*AiVector3D) (out []float32) {
	for _, v := range in {
		out = append(out, v.X, v.Y, v.Z)
	}
	return out
}

func IsCCW(in []*AiVector2D) bool {
	npoints := len(in)
	var aa, bb, cc, b, c, theta float64
	var convex_turn float64
	convex_sum := 0.

	if npoints < 3 {
		panic("npoints < 3")
	}

	for i := 0; i < npoints-2; i++ {
		aa = float64((in[i+2].X-in[i].X)*(in[i+2].X-in[i].X)) +
			float64((-in[i+2].Y+in[i].Y)*(-in[i+2].Y+in[i].Y))

		bb = float64((in[i+1].X-in[i].X)*(in[i+1].X-in[i].X)) +
			float64((-in[i+1].Y+in[i].Y)*(-in[i+1].Y+in[i].Y))

		cc = (float64(in[i+2].X-in[i+1].X) *
			float64(in[i+2].X-in[i+1].X)) +
			(float64(-in[i+2].Y+in[i+1].Y) *
				float64(-in[i+2].Y+in[i+1].Y))

		b = math.Sqrt(bb)
		c = math.Sqrt(cc)
		theta = math.Acos((bb + cc - aa) / (2 * b * c))

		if OnLeftSideOfLine2D(in[i], in[i+2], in[i+1]) == 1 {
			//  if (convex(in[i].x, in[i].y,
			//      in[i+1].x, in[i+1].y,
			//      in[i+2].x, in[i+2].y)) {
			convex_turn = math.Pi - theta
			convex_sum += convex_turn
		} else {
			convex_sum -= math.Pi - theta
		}
	}
	aa = (float64(in[1].X-in[npoints-2].X) *
		float64(in[1].X-in[npoints-2].X)) +
		(float64(-in[1].Y+in[npoints-2].Y) *
			float64(-in[1].Y+in[npoints-2].Y))

	bb = (float64(in[0].X-in[npoints-2].X) *
		float64(in[0].X-in[npoints-2].X)) +
		(float64(-in[0].Y+in[npoints-2].Y) *
			float64(-in[0].Y+in[npoints-2].Y))

	cc = (float64(in[1].X-in[0].X) * float64(in[1].X-in[0].X)) +
		(float64(-in[1].Y+in[0].Y) * float64(-in[1].Y+in[0].Y))

	b = math.Sqrt(bb)
	c = math.Sqrt(cc)
	theta = math.Acos((bb + cc - aa) / (2 * b * c))

	//if (convex(in[npoints-2].x, in[npoints-2].y,
	//  in[0].x, in[0].y,
	//  in[1].x, in[1].y)) {
	if OnLeftSideOfLine2D(in[npoints-2], in[1], in[0]) == 1 {
		convex_turn = math.Pi - theta
		convex_sum += convex_turn
	} else {
		convex_sum -= math.Pi - theta
	}

	return convex_sum >= (2 * math.Pi)
}

// -------------------------------------------------------------------------------
/** Compute the normal of an arbitrary polygon in R3.
 *
 *  The code is based on Newell's formula, that is a polygons normal is the ratio
 *  of its area when projected onto the three coordinate axes.
 *
 *  @param out Receives the output normal
 *  @param num Number of input vertices
 *  @param x X data source. x[ofs_x*n] is the n'th element.
 *  @param y Y data source. y[ofs_y*n] is the y'th element
 *  @param z Z data source. z[ofs_z*n] is the z'th element
 *
 *  @note The data arrays must have storage for at least num+2 elements. Using
 *  this method is much faster than the 'other' NewellNormal()
 */

func NewellNormal(num int, x, y, z []float32, ofs_x, ofs_y, ofs_z int) (out *AiVector3D) {
	// Duplicate the first two vertices at the end
	x[(num+0)*ofs_x] = x[0]
	x[(num+1)*ofs_x] = x[ofs_x]

	y[(num+0)*ofs_y] = y[0]
	y[(num+1)*ofs_y] = y[ofs_y]

	z[(num+0)*ofs_z] = z[0]
	z[(num+1)*ofs_z] = z[ofs_z]

	sum_xy := float32(0.0)
	sum_yz := float32(0.0)
	sum_zx := float32(0.0)

	xptr := 0 + ofs_x
	xlow := 0
	xhigh := 0 + ofs_x*2
	yptr := 0 + ofs_y
	ylow := 0
	yhigh := 0 + ofs_y*2
	zptr := 0 + ofs_z
	zlow := 0
	zhigh := 0 + ofs_z*2

	for tmp := 0; tmp < num; tmp++ {
		sum_xy += x[xptr] * (y[yhigh] - y[ylow])
		sum_yz += y[yptr] * (z[zhigh] - z[zlow])
		sum_zx += z[zptr] * (x[xhigh] - x[xlow])

		xptr += ofs_x
		xlow += ofs_x
		xhigh += ofs_x

		yptr += ofs_y
		ylow += ofs_y
		yhigh += ofs_y

		zptr += ofs_z
		zlow += ofs_z
		zhigh += ofs_z
	}
	return NewAiVector3D3(sum_yz, sum_zx, sum_xy)
}
