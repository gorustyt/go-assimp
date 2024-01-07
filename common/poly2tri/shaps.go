package poly2tri

import "github.com/gorustyt/go-assimp/common/logger"

// Represents a simple polygon's edge
type Edge struct {
	p, q *Point
}

func NewEdge(p1, p2 *Point) *Edge {
	e := &Edge{}
	if p1.y > p2.y {
		e.q = p1
		e.p = p2
	} else if p1.y == p2.y {
		if p1.x > p2.x {
			e.q = p1
			e.p = p2
		} else if p1.x == p2.x {
			// Repeat points
			// ASSIMP_CHANGE (aramis_acg)
			logger.FatalF("repeat points")
			//assert(false);
		}
	}
	e.q.edge_list = append(e.q.edge_list, e)
	return e
}

type Triangle struct {
	/// Flags to determine if an edge is a Constrained edge
	constrained_edge [3]bool
	/// Flags to determine if an edge is a Delauney edge
	delaunay_edge [3]bool
	/// Triangle points
	points_ [3]*Point
	/// Neighbor list
	neighbors_ [3]*Triangle

	/// Has this triangle been marked as an interior triangle?
	interior_ bool
}

func NewTriangle(a, b, c *Point) *Triangle {
	tri := &Triangle{}
	tri.points_[0] = a
	tri.points_[1] = b
	tri.points_[2] = c
	tri.neighbors_[0] = nil
	tri.neighbors_[1] = nil
	tri.neighbors_[2] = nil
	tri.constrained_edge[2] = false
	tri.constrained_edge[1] = tri.constrained_edge[2]
	tri.constrained_edge[0] = tri.constrained_edge[1]

	tri.delaunay_edge[2] = false
	tri.delaunay_edge[1] = tri.delaunay_edge[2]
	tri.delaunay_edge[0] = tri.delaunay_edge[1]
	tri.interior_ = false
	return tri
}

func (tri *Triangle) Add(a, b *Point) *Point {
	return NewPoint(a.x+b.x, a.y+b.y)
}

func (tri *Triangle) Sub(a, b *Point) *Point {
	return NewPoint(a.x-b.x, a.y-b.y)
}

func (tri *Triangle) Mul(s float64, a *Point) *Point {
	return NewPoint(s*a.x, s*a.y)
}

// / Peform the dot product on two vectors.
func (tri *Triangle) Dot(a, b *Point) float64 {
	return a.x*b.x + a.y*b.y
}

// / Perform the cross product on two vectors. In 2D this produces a scalar.
func (tri *Triangle) Cross(a, b *Point) float64 {
	return a.x*b.y - a.y*b.x
}

// / Perform the cross product on a point and a scalar. In 2D this produces
// / a point.
func (tri *Triangle) Cross1(a *Point, s float64) *Point {
	return NewPoint(s*a.y, -s*a.x)
}

// / Perform the cross product on a scalar and a point. In 2D this produces
// / a point.
func (tri *Triangle) Cross2(s float64, a *Point) *Point {
	return NewPoint(-s*a.y, s*a.x)
}

func (tri *Triangle) GetPoint(index int) *Point {
	return tri.points_[index]
}

func (tri *Triangle) GetNeighbor(index int) *Triangle {
	return tri.neighbors_[index]
}

func (tri *Triangle) ContainsPoint(p *Point) bool {
	return p == tri.points_[0] || p == tri.points_[1] || p == tri.points_[2]
}

func (tri *Triangle) ContainsEdge(e *Edge) bool {
	return tri.ContainsPoint(e.p) && tri.ContainsPoint(e.q)
}

func (tri *Triangle) ContainsPoints(p *Point, q *Point) bool {
	return tri.ContainsPoint(p) && tri.ContainsPoint(q)
}

func (tri *Triangle) IsInterior() bool {
	return tri.interior_
}

func (tri *Triangle) SetInterior(b bool) {
	tri.interior_ = b
}

// Update neighbor pointers
func (tri *Triangle) MarkNeighborPointers(p1, p2 *Point, t *Triangle) {
	if (p1 == tri.points_[2] && p2 == tri.points_[1]) || (p1 == tri.points_[1] && p2 == tri.points_[2]) {
		tri.neighbors_[0] = t
	} else if (p1 == tri.points_[0] && p2 == tri.points_[2]) || (p1 == tri.points_[2] && p2 == tri.points_[0]) {
		tri.neighbors_[1] = t
	} else if (p1 == tri.points_[0] && p2 == tri.points_[1]) || (p1 == tri.points_[1] && p2 == tri.points_[0]) {
		tri.neighbors_[2] = t
	} else {
		logger.FatalF("Update neighbor pointers error!")
	}

}

// Exhaustive search to update neighbor pointers
func (tri *Triangle) MarkNeighbor(t *Triangle) {
	if t.ContainsPoints(tri.points_[1], tri.points_[2]) {
		tri.neighbors_[0] = t
		t.MarkNeighborPointers(tri.points_[1], tri.points_[2], tri)
	} else if t.ContainsPoints(tri.points_[0], tri.points_[2]) {
		tri.neighbors_[1] = t
		t.MarkNeighborPointers(tri.points_[0], tri.points_[2], tri)
	} else if t.ContainsPoints(tri.points_[0], tri.points_[1]) {
		tri.neighbors_[2] = t
		t.MarkNeighborPointers(tri.points_[0], tri.points_[1], tri)
	}
}

/**
 * Clears all references to all other triangles and points
 */
func (tri *Triangle) Clear() {
	var t *Triangle
	for i := 0; i < 3; i++ {
		t = tri.neighbors_[i]
		if t != nil {
			t.ClearNeighbor(tri)
		}
	}
	tri.ClearNeighbors()
	tri.points_[2] = nil
	tri.points_[1] = tri.points_[2]
	tri.points_[0] = tri.points_[1]
}

func (tri *Triangle) ClearNeighbor(triangle *Triangle) {
	if tri.neighbors_[0] == triangle {
		tri.neighbors_[0] = nil
	} else if tri.neighbors_[1] == triangle {
		tri.neighbors_[1] = nil
	} else {
		tri.neighbors_[2] = nil
	}
}

func (tri *Triangle) ClearNeighbors() {
	tri.neighbors_[0] = nil
	tri.neighbors_[1] = nil
	tri.neighbors_[2] = nil
}

func (tri *Triangle) ClearDelunayEdges() {
	tri.delaunay_edge[2] = false
	tri.delaunay_edge[1] = tri.delaunay_edge[2]
	tri.delaunay_edge[0] = tri.delaunay_edge[1]
}

func (tri *Triangle) OppositePoint(t *Triangle, p *Point) *Point {
	cw := t.PointCW(p)
	return tri.PointCW(cw)
}

// Legalized triangle by rotating clockwise around point(0)
func (tri *Triangle) Legalize(point *Point) {
	tri.points_[1] = tri.points_[0]
	tri.points_[0] = tri.points_[2]
	tri.points_[2] = point
}

// Legalize triagnle by rotating clockwise around oPoint
func (tri *Triangle) Legalize1(opoint *Point, npoint *Point) {
	if opoint == tri.points_[0] {
		tri.points_[1] = tri.points_[0]
		tri.points_[0] = tri.points_[2]
		tri.points_[2] = npoint
	} else if opoint == tri.points_[1] {
		tri.points_[2] = tri.points_[1]
		tri.points_[1] = tri.points_[0]
		tri.points_[0] = npoint
	} else if opoint == tri.points_[2] {
		tri.points_[0] = tri.points_[2]
		tri.points_[2] = tri.points_[1]
		tri.points_[1] = npoint
	} else {
		logger.FatalF("Legalize triagnle by rotating clockwise around oPoint ")
	}
}

func (tri *Triangle) Index(p *Point) int {
	if p == tri.points_[0] {
		return 0
	} else if p == tri.points_[1] {
		return 1
	} else if p == tri.points_[2] {
		return 2
	}
	logger.FatalF("find error index ")
	return -1
}

func (tri *Triangle) EdgeIndex(p1, p2 *Point) int {
	if tri.points_[0] == p1 {
		if tri.points_[1] == p2 {
			return 2
		} else if tri.points_[2] == p2 {
			return 1
		}
	} else if tri.points_[1] == p1 {
		if tri.points_[2] == p2 {
			return 0
		} else if tri.points_[0] == p2 {
			return 2
		}
	} else if tri.points_[2] == p1 {
		if tri.points_[0] == p2 {
			return 1
		} else if tri.points_[1] == p2 {
			return 0
		}
	}
	return -1
}
func (tri *Triangle) MarkConstrainedEdgeIndex(index int) {
	tri.constrained_edge[index] = true
}

func (tri *Triangle) MarkConstrainedEdgeEdge(edge *Edge) {
	tri.MarkConstrainedEdge(edge.p, edge.q)
}

// Mark edge as constrained
func (tri *Triangle) MarkConstrainedEdge(p, q *Point) {
	if (q == tri.points_[0] && p == tri.points_[1]) || (q == tri.points_[1] && p == tri.points_[0]) {
		tri.constrained_edge[2] = true
	} else if (q == tri.points_[0] && p == tri.points_[2]) || (q == tri.points_[2] && p == tri.points_[0]) {
		tri.constrained_edge[1] = true
	} else if (q == tri.points_[1] && p == tri.points_[2]) || (q == tri.points_[2] && p == tri.points_[1]) {
		tri.constrained_edge[0] = true
	}
}

// The point counter-clockwise to given point
func (tri *Triangle) PointCW(point *Point) *Point {
	if point == tri.points_[0] {
		return tri.points_[2]
	} else if point == tri.points_[1] {
		return tri.points_[0]
	} else if point == tri.points_[2] {
		return tri.points_[1]
	}
	logger.FatalF("The point counter-clockwise to given point")
	return nil
}

// The point counter-clockwise to given point
func (tri *Triangle) PointCCW(point *Point) *Point {
	if point == tri.points_[0] {
		return tri.points_[1]
	} else if point == tri.points_[1] {
		return tri.points_[2]
	} else if point == tri.points_[2] {
		return tri.points_[0]
	}
	logger.FatalF("The point counter-clockwise to given point")
	return nil
}

// The neighbor clockwise to given point
func (tri *Triangle) NeighborCW(point *Point) *Triangle {
	if point == tri.points_[0] {
		return tri.neighbors_[1]
	} else if point == tri.points_[1] {
		return tri.neighbors_[2]
	}
	return tri.neighbors_[0]
}

// The neighbor counter-clockwise to given point
func (tri *Triangle) NeighborCCW(point *Point) *Triangle {
	if point == tri.points_[0] {
		return tri.neighbors_[2]
	} else if point == tri.points_[1] {
		return tri.neighbors_[0]
	}
	return tri.neighbors_[1]
}

func (tri *Triangle) GetConstrainedEdgeCCW(p *Point) bool {
	if p == tri.points_[0] {
		return tri.constrained_edge[2]
	} else if p == tri.points_[1] {
		return tri.constrained_edge[0]
	}
	return tri.constrained_edge[1]
}

func (tri *Triangle) GetConstrainedEdgeCW(p *Point) bool {
	if p == tri.points_[0] {
		return tri.constrained_edge[1]
	} else if p == tri.points_[1] {
		return tri.constrained_edge[2]
	}
	return tri.constrained_edge[0]
}

func (tri *Triangle) SetConstrainedEdgeCCW(p *Point, ce bool) {
	if p == tri.points_[0] {
		tri.constrained_edge[2] = ce
	} else if p == tri.points_[1] {
		tri.constrained_edge[0] = ce
	} else {
		tri.constrained_edge[1] = ce
	}
}

func (tri *Triangle) SetConstrainedEdgeCW(p *Point, ce bool) {
	if p == tri.points_[0] {
		tri.constrained_edge[1] = ce
	} else if p == tri.points_[1] {
		tri.constrained_edge[2] = ce
	} else {
		tri.constrained_edge[0] = ce
	}
}

func (tri *Triangle) GetDelunayEdgeCCW(p *Point) bool {
	if p == tri.points_[0] {
		return tri.delaunay_edge[2]
	} else if p == tri.points_[1] {
		return tri.delaunay_edge[0]
	}
	return tri.delaunay_edge[1]
}

func (tri *Triangle) GetDelunayEdgeCW(p *Point) bool {
	if p == tri.points_[0] {
		return tri.delaunay_edge[1]
	} else if p == tri.points_[1] {
		return tri.delaunay_edge[2]
	}
	return tri.delaunay_edge[0]
}

func (tri *Triangle) SetDelunayEdgeCCW(p *Point, e bool) {
	if p == tri.points_[0] {
		tri.delaunay_edge[2] = e
	} else if p == tri.points_[1] {
		tri.delaunay_edge[0] = e
	} else {
		tri.delaunay_edge[1] = e
	}
}

func (tri *Triangle) SetDelunayEdgeCW(p *Point, e bool) {
	if p == tri.points_[0] {
		tri.delaunay_edge[1] = e
	} else if p == tri.points_[1] {
		tri.delaunay_edge[2] = e
	} else {
		tri.delaunay_edge[0] = e
	}
}

// The neighbor across to given point
func (tri *Triangle) NeighborAcross(opoint *Point) *Triangle {
	if opoint == tri.points_[0] {
		return tri.neighbors_[0]
	} else if opoint == tri.points_[1] {
		return tri.neighbors_[1]
	}
	return tri.neighbors_[2]
}

func (tri *Triangle) DebugPrint() {
	logger.InfoF("x :%v ,y:%v", tri.points_[0].x, tri.points_[0].y)
	logger.InfoF("x :%v ,y:%v", tri.points_[1].x, tri.points_[1].y)
	logger.InfoF("x :%v ,y:%v", tri.points_[2].x, tri.points_[2].y)
}
