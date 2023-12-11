package poly2tri

import (
	"assimp/common/logger"
	"math"
)

type Sweep struct {
	nodes_ []*Node
}

func NewSweep() *Sweep {
	return &Sweep{}
}

// Triangulate simple polygon with holes
func (s *Sweep) Triangulate(tcx *SweepContext) {
	tcx.InitTriangulation()
	tcx.CreateAdvancingFront(s.nodes_)
	// Sweep points; build mesh
	s.SweepPoints(tcx)
	// Clean up
	s.FinalizationPolygon(tcx)
}

func (s *Sweep) SweepPoints(tcx *SweepContext) {
	for i := 1; i < tcx.point_count(); i++ {
		point := tcx.GetPoint(i)
		node := s.PointEvent(tcx, point)
		for ii := 0; ii < len(point.edge_list); ii++ {
			s.EdgeEvent(tcx, point.edge_list[ii], node)
		}
	}
}

func (s *Sweep) FinalizationPolygon(tcx *SweepContext) {
	// Get an Internal triangle to start with
	t := tcx.Front().Head().next.triangle
	p := tcx.Front().Head().next.point
	for t.GetConstrainedEdgeCW(p) {
		t = t.NeighborCCW(p)
	}

	// Collect interior triangles constrained by edges
	tcx.MeshClean(t)
}

func (s *Sweep) PointEvent(tcx *SweepContext, point *Point) *Node {
	node := tcx.LocateNode(point)
	new_node := s.NewFrontTriangle(tcx, point, node)

	// Only need to check +epsilon since point never have smaller
	// x value than node due to how we fetch nodes from the front
	if point.x <= node.point.x+EPSILON {
		s.Fill(tcx, node)
	}

	//tcx.AddNode(new_node);

	s.FillAdvancingFront(tcx, new_node)
	return new_node
}

func (s *Sweep) EdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	tcx.edge_event.constrained_edge = edge
	tcx.edge_event.right = (edge.p.x > edge.q.x)

	if s.IsEdgeSideOfTriangle(node.triangle, edge.p, edge.q) {
		return
	}

	// For now we will do all needed filling
	// TODO: integrate with flip process might give some better performance
	//       but for now this avoid the issue with cases that needs both flips and fills
	s.FillEdgeEvent(tcx, edge, node)
	s.EdgeEvent1(tcx, edge.p, edge.q, node.triangle, edge.q)
}
func (s *Sweep) EdgeEvent1(tcx *SweepContext, ep, eq *Point, triangle *Triangle, point *Point) {
	if s.IsEdgeSideOfTriangle(triangle, ep, eq) {
		return
	}

	p1 := triangle.PointCCW(point)
	o1 := Orient2d(eq, p1, ep)
	if o1 == COLLINEAR {

		if triangle.ContainsPoints(eq, p1) {
			triangle.MarkConstrainedEdge(eq, p1)
			// We are modifying the constraint maybe it would be better to
			// not change the given constraint and just keep a variable for the new constraint
			tcx.edge_event.constrained_edge.q = p1
			triangle = triangle.NeighborAcross(point)
			s.EdgeEvent1(tcx, ep, p1, triangle, p1)
		} else {
			// ASSIMP_CHANGE (aramis_acg)
			logger.Fatal("EdgeEvent - collinear points not supported")
		}
		return
	}

	p2 := triangle.PointCW(point)
	o2 := Orient2d(eq, p2, ep)
	if o2 == COLLINEAR {

		if triangle.ContainsPoints(eq, p2) {
			triangle.MarkConstrainedEdge(eq, p2)
			// We are modifying the constraint maybe it would be better to
			// not change the given constraint and just keep a variable for the new constraint
			tcx.edge_event.constrained_edge.q = p2
			triangle = triangle.NeighborAcross(point)
			s.EdgeEvent1(tcx, ep, p2, triangle, p2)
		} else {
			// ASSIMP_CHANGE (aramis_acg)
			logger.Fatal("EdgeEvent - collinear points not supported")
		}
		return
	}

	if o1 == o2 {
		// Need to decide if we are rotating CW or CCW to get to a triangle
		// that will cross edge
		if o1 == CW {
			triangle = triangle.NeighborCCW(point)
		} else {
			triangle = triangle.NeighborCW(point)
		}
		s.EdgeEvent1(tcx, ep, eq, triangle, point)
	} else {
		// This triangle crosses constraint so lets flippin start!
		s.FlipEdgeEvent(tcx, ep, eq, triangle, point)
	}
}

func (s *Sweep) IsEdgeSideOfTriangle(triangle *Triangle, ep, eq *Point) bool {
	index := triangle.EdgeIndex(ep, eq)

	if index != -1 {
		triangle.MarkConstrainedEdgeIndex(index)
		t := triangle.GetNeighbor(index)
		if t != nil {
			t.MarkConstrainedEdge(ep, eq)
		}
		return true
	}
	return false
}

func (s *Sweep) NewFrontTriangle(tcx *SweepContext, point *Point, node *Node) *Node {
	triangle := NewTriangle(point, node.point, node.next.point)

	triangle.MarkNeighbor(node.triangle)
	tcx.AddToMap(triangle)

	new_node := NewNode(point)
	s.nodes_ = append(s.nodes_, new_node)

	new_node.next = node.next
	new_node.prev = node
	node.next.prev = new_node
	node.next = new_node

	if !s.Legalize(tcx, triangle) {
		tcx.MapTriangleToNodes(triangle)
	}

	return new_node
}

func (s *Sweep) Fill(tcx *SweepContext, node *Node) {
	triangle := NewTriangle(node.prev.point, node.point, node.next.point)

	// TODO: should copy the constrained_edge value from neighbor triangles
	//       for now constrained_edge values are copied during the legalize
	triangle.MarkNeighbor(node.prev.triangle)
	triangle.MarkNeighbor(node.triangle)

	tcx.AddToMap(triangle)

	// Update the advancing front
	node.prev.next = node.next
	node.next.prev = node.prev

	// If it was legalized the triangle has already been mapped
	if !s.Legalize(tcx, triangle) {
		tcx.MapTriangleToNodes(triangle)
	}

}
func (s *Sweep) FillAdvancingFront(tcx *SweepContext, n *Node) {

	// Fill right holes
	node := n.next

	for node.next != nil {
		// if HoleAngle exceeds 90 degrees then break.
		if s.LargeHole_DontFill(node) {
			break
		}
		s.Fill(tcx, node)
		node = node.next
	}

	// Fill left holes
	node = n.prev

	for node.prev != nil {
		// if HoleAngle exceeds 90 degrees then break.
		if s.LargeHole_DontFill(node) {
			break
		}
		s.Fill(tcx, node)
		node = node.prev
	}

	// Fill right basins
	if n.next != nil && n.next.next != nil {
		angle := s.BasinAngle(n)
		if angle < PI_3div4 {
			s.FillBasin(tcx, n)
		}
	}
}

// True if HoleAngle exceeds 90 degrees.
func (s *Sweep) LargeHole_DontFill(node *Node) bool {

	nextNode := node.next
	prevNode := node.prev
	if !s.AngleExceeds90Degrees(node.point, nextNode.point, prevNode.point) {
		return false
	}

	// Check additional points on front.
	next2Node := nextNode.next
	// "..Plus.." because only want angles on same side as point being added.
	if (next2Node != nil) && !s.AngleExceedsPlus90DegreesOrIsNegative(node.point, next2Node.point, prevNode.point) {
		return false
	}

	prev2Node := prevNode.prev
	// "..Plus.." because only want angles on same side as point being added.
	if (prev2Node != nil) && !s.AngleExceedsPlus90DegreesOrIsNegative(node.point, nextNode.point, prev2Node.point) {
		return false
	}

	return true
}

func (s *Sweep) AngleExceeds90Degrees(origin, pa, pb *Point) bool {
	angle := s.Angle(origin, pa, pb)
	return ((angle > PI_div2) || (angle < -PI_div2))
}

func (s *Sweep) AngleExceedsPlus90DegreesOrIsNegative(origin, pa, pb *Point) bool {
	angle := s.Angle(origin, pa, pb)
	return (angle > PI_div2) || (angle < 0)
}

func (s *Sweep) Angle(origin, pa, pb *Point) float64 {
	/* Complex plane
	 * ab = cosA +i*sinA
	 * ab = (ax + ay*i)(bx + by*i) = (ax*bx + ay*by) + i(ax*by-ay*bx)
	 * atan2(y,x) computes the principal value of the argument function
	 * applied to the complex number x+iy
	 * Where x = ax*bx + ay*by
	 *       y = ax*by - ay*bx
	 */
	px := origin.x
	py := origin.y
	ax := pa.x - px
	ay := pa.y - py
	bx := pb.x - px
	by := pb.y - py
	x := ax*by - ay*bx
	y := ax*bx + ay*by
	return math.Atan2(x, y)
}

func (s *Sweep) BasinAngle(node *Node) float64 {
	ax := node.point.x - node.next.next.point.x
	ay := node.point.y - node.next.next.point.y
	return math.Atan2(ay, ax)
}

func (s *Sweep) HoleAngle(node *Node) float64 {
	/* Complex plane
	 * ab = cosA +i*sinA
	 * ab = (ax + ay*i)(bx + by*i) = (ax*bx + ay*by) + i(ax*by-ay*bx)
	 * atan2(y,x) computes the principal value of the argument function
	 * applied to the complex number x+iy
	 * Where x = ax*bx + ay*by
	 *       y = ax*by - ay*bx
	 */
	ax := node.next.point.x - node.point.x
	ay := node.next.point.y - node.point.y
	bx := node.prev.point.x - node.point.x
	by := node.prev.point.y - node.point.y
	return math.Atan2(ax*by-ay*bx, ax*bx+ay*by)
}

func (s *Sweep) Legalize(tcx *SweepContext, t *Triangle) bool {
	// To legalize a triangle we start by finding if any of the three edges
	// violate the Delaunay condition
	for i := 0; i < 3; i++ {
		if t.delaunay_edge[i] {
			continue
		}

		ot := t.GetNeighbor(i)

		if ot != nil {
			p := t.GetPoint(i)
			op := ot.OppositePoint(t, p)
			oi := ot.Index(op)

			// If this is a Constrained Edge or a Delaunay Edge(only during recursive legalization)
			// then we should not try to legalize
			if ot.constrained_edge[oi] || ot.delaunay_edge[oi] {
				t.constrained_edge[i] = ot.constrained_edge[oi]
				continue
			}

			inside := s.Incircle(p, t.PointCCW(p), t.PointCW(p), op)

			if inside {
				// Lets mark this shared edge as Delaunay
				t.delaunay_edge[i] = true
				ot.delaunay_edge[oi] = true

				// Lets rotate shared edge one vertex CW to legalize it
				s.RotateTrianglePair(t, p, ot, op)

				// We now got one valid Delaunay Edge shared by two triangles
				// This gives us 4 new edges to check for Delaunay

				// Make sure that triangle to node mapping is done only one time for a specific triangle
				not_legalized := !s.Legalize(tcx, t)
				if not_legalized {
					tcx.MapTriangleToNodes(t)
				}

				not_legalized = !s.Legalize(tcx, ot)
				if not_legalized {
					tcx.MapTriangleToNodes(ot)
				}

				// Reset the Delaunay edges, since they only are valid Delaunay edges
				// until we add a new triangle or point.
				// XXX: need to think about this. Can these edges be tried after we
				//      return to previous recursive level?
				t.delaunay_edge[i] = false
				ot.delaunay_edge[oi] = false

				// If triangle have been legalized no need to check the other edges since
				// the recursive legalization will handles those so we can end here.
				return true
			}
		}
	}
	return false
}

func (s *Sweep) Incircle(pa, pb, pc, pd *Point) bool {
	adx := pa.x - pd.x
	ady := pa.y - pd.y
	bdx := pb.x - pd.x
	bdy := pb.y - pd.y

	adxbdy := adx * bdy
	bdxady := bdx * ady
	oabd := adxbdy - bdxady

	if oabd <= 0 {
		return false
	}

	cdx := pc.x - pd.x
	cdy := pc.y - pd.y

	cdxady := cdx * ady
	adxcdy := adx * cdy
	ocad := cdxady - adxcdy

	if ocad <= 0 {
		return false
	}

	bdxcdy := bdx * cdy
	cdxbdy := cdx * bdy

	alift := adx*adx + ady*ady
	blift := bdx*bdx + bdy*bdy
	clift := cdx*cdx + cdy*cdy

	det := alift*(bdxcdy-cdxbdy) + blift*ocad + clift*oabd

	return det > 0
}

func (s *Sweep) RotateTrianglePair(t *Triangle, p *Point, ot *Triangle, op *Point) {

	n1 := t.NeighborCCW(p)
	n2 := t.NeighborCW(p)
	n3 := ot.NeighborCCW(op)
	n4 := ot.NeighborCW(op)

	ce1 := t.GetConstrainedEdgeCCW(p)
	ce2 := t.GetConstrainedEdgeCW(p)
	ce3 := ot.GetConstrainedEdgeCCW(op)
	ce4 := ot.GetConstrainedEdgeCW(op)

	de1 := t.GetDelunayEdgeCCW(p)
	de2 := t.GetDelunayEdgeCW(p)
	de3 := ot.GetDelunayEdgeCCW(op)
	de4 := ot.GetDelunayEdgeCW(op)

	t.Legalize1(p, op)
	ot.Legalize1(op, p)

	// Remap delaunay_edge
	ot.SetDelunayEdgeCCW(p, de1)
	t.SetDelunayEdgeCW(p, de2)
	t.SetDelunayEdgeCCW(op, de3)
	ot.SetDelunayEdgeCW(op, de4)

	// Remap constrained_edge
	ot.SetConstrainedEdgeCCW(p, ce1)
	t.SetConstrainedEdgeCW(p, ce2)
	t.SetConstrainedEdgeCCW(op, ce3)
	ot.SetConstrainedEdgeCW(op, ce4)

	// Remap neighbors
	// XXX: might optimize the markNeighbor by keeping track of
	//      what side should be assigned to what neighbor after the
	//      rotation. Now mark neighbor does lots of testing to find
	//      the right side.
	t.ClearNeighbors()
	ot.ClearNeighbors()
	if n1 != nil {
		ot.MarkNeighbor(n1)
	}
	if n2 != nil {
		t.MarkNeighbor(n2)
	}
	if n3 != nil {
		t.MarkNeighbor(n3)
	}
	if n4 != nil {
		ot.MarkNeighbor(n4)
	}
	t.MarkNeighbor(ot)
}

func (s *Sweep) FillBasin(tcx *SweepContext, node *Node) {
	if Orient2d(node.point, node.next.point, node.next.next.point) == CCW {
		tcx.basin.left_node = node.next.next
	} else {
		tcx.basin.left_node = node.next
	}

	// Find the bottom and right node
	tcx.basin.bottom_node = tcx.basin.left_node
	for tcx.basin.bottom_node.next != nil && tcx.basin.bottom_node.point.y >= tcx.basin.bottom_node.next.point.y {
		tcx.basin.bottom_node = tcx.basin.bottom_node.next
	}
	if tcx.basin.bottom_node == tcx.basin.left_node {
		// No valid basin
		return
	}

	tcx.basin.right_node = tcx.basin.bottom_node
	for tcx.basin.right_node.next != nil && tcx.basin.right_node.point.y < tcx.basin.right_node.next.point.y {
		tcx.basin.right_node = tcx.basin.right_node.next
	}
	if tcx.basin.right_node == tcx.basin.bottom_node {
		// No valid basins
		return
	}

	tcx.basin.width = tcx.basin.right_node.point.x - tcx.basin.left_node.point.x
	tcx.basin.left_highest = tcx.basin.left_node.point.y > tcx.basin.right_node.point.y

	s.FillBasinReq(tcx, tcx.basin.bottom_node)
}

func (s *Sweep) FillBasinReq(tcx *SweepContext, node *Node) {
	// if shallow stop filling
	if s.IsShallow(tcx, node) {
		return
	}

	s.Fill(tcx, node)

	if node.prev == tcx.basin.left_node && node.next == tcx.basin.right_node {
		return
	} else if node.prev == tcx.basin.left_node {
		o := Orient2d(node.point, node.next.point, node.next.next.point)
		if o == CW {
			return
		}
		node = node.next
	} else if node.next == tcx.basin.right_node {
		o := Orient2d(node.point, node.prev.point, node.prev.prev.point)
		if o == CCW {
			return
		}
		node = node.prev
	} else {
		// Continue with the neighbor node with lowest Y value
		if node.prev.point.y < node.next.point.y {
			node = node.prev
		} else {
			node = node.next
		}
	}

	s.FillBasinReq(tcx, node)
}

func (s *Sweep) IsShallow(tcx *SweepContext, node *Node) bool {
	var height float64

	if tcx.basin.left_highest {
		height = tcx.basin.left_node.point.y - node.point.y
	} else {
		height = tcx.basin.right_node.point.y - node.point.y
	}

	// if shallow stop filling
	if tcx.basin.width > height {
		return true
	}
	return false
}

func (s *Sweep) FillEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	if tcx.edge_event.right {
		s.FillRightAboveEdgeEvent(tcx, edge, node)
	} else {
		s.FillLeftAboveEdgeEvent(tcx, edge, node)
	}
}

func (s *Sweep) FillRightAboveEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	for node.next.point.x < edge.p.x {
		// Check if next node is below the edge
		if Orient2d(edge.q, node.next.point, edge.p) == CCW {
			s.FillRightBelowEdgeEvent(tcx, edge, node)
		} else {
			node = node.next
		}
	}
}

func (s *Sweep) FillRightBelowEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	if node.point.x < edge.p.x {
		if Orient2d(node.point, node.next.point, node.next.next.point) == CCW {
			// Concave
			s.FillRightConcaveEdgeEvent(tcx, edge, node)
		} else {
			// Convex
			s.FillRightConvexEdgeEvent(tcx, edge, node)
			// Retry this one
			s.FillRightBelowEdgeEvent(tcx, edge, node)
		}
	}
}

func (s *Sweep) FillRightConcaveEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	s.Fill(tcx, node.next)
	if node.next.point != edge.p {
		// Next above or below edge?
		if Orient2d(edge.q, node.next.point, edge.p) == CCW {
			// Below
			if Orient2d(node.point, node.next.point, node.next.next.point) == CCW {
				// Next is concave
				s.FillRightConcaveEdgeEvent(tcx, edge, node)
			} else {
				// Next is convex
			}
		}
	}

}

func (s *Sweep) FillRightConvexEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	// Next concave or convex?
	if Orient2d(node.next.point, node.next.next.point, node.next.next.next.point) == CCW {
		// Concave
		s.FillRightConcaveEdgeEvent(tcx, edge, node.next)
	} else {
		// Convex
		// Next above or below edge?
		if Orient2d(edge.q, node.next.next.point, edge.p) == CCW {
			// Below
			s.FillRightConvexEdgeEvent(tcx, edge, node.next)
		} else {
			// Above
		}
	}
}

func (s *Sweep) FillLeftAboveEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	for node.prev.point.x > edge.p.x {
		// Check if next node is below the edge
		if Orient2d(edge.q, node.prev.point, edge.p) == CW {
			s.FillLeftBelowEdgeEvent(tcx, edge, node)
		} else {
			node = node.prev
		}
	}
}

func (s *Sweep) FillLeftBelowEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	if node.point.x > edge.p.x {
		if Orient2d(node.point, node.prev.point, node.prev.prev.point) == CW {
			// Concave
			s.FillLeftConcaveEdgeEvent(tcx, edge, node)
		} else {
			// Convex
			s.FillLeftConvexEdgeEvent(tcx, edge, node)
			// Retry this one
			s.FillLeftBelowEdgeEvent(tcx, edge, node)
		}
	}
}

func (s *Sweep) FillLeftConvexEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	// Next concave or convex?
	if Orient2d(node.prev.point, node.prev.prev.point, node.prev.prev.prev.point) == CW {
		// Concave
		s.FillLeftConcaveEdgeEvent(tcx, edge, node.prev)
	} else {
		// Convex
		// Next above or below edge?
		if Orient2d(edge.q, node.prev.prev.point, edge.p) == CW {
			// Below
			s.FillLeftConvexEdgeEvent(tcx, edge, node.prev)
		} else {
			// Above
		}
	}
}

func (s *Sweep) FillLeftConcaveEdgeEvent(tcx *SweepContext, edge *Edge, node *Node) {
	s.Fill(tcx, node.prev)
	if node.prev.point != edge.p {
		// Next above or below edge?
		if Orient2d(edge.q, node.prev.point, edge.p) == CW {
			// Below
			if Orient2d(node.point, node.prev.point, node.prev.prev.point) == CW {
				// Next is concave
				s.FillLeftConcaveEdgeEvent(tcx, edge, node)
			} else {
				// Next is convex
			}
		}
	}

}

func (s *Sweep) FlipEdgeEvent(tcx *SweepContext, ep, eq *Point, t *Triangle, p *Point) {
	ot := t.NeighborAcross(p)
	op := ot.OppositePoint(t, p)

	if InScanArea(p, t.PointCCW(p), t.PointCW(p), op) {
		// Lets rotate shared edge one vertex CW
		s.RotateTrianglePair(t, p, ot, op)
		tcx.MapTriangleToNodes(t)
		tcx.MapTriangleToNodes(ot)

		if p == eq && op == ep {
			if eq == tcx.edge_event.constrained_edge.q && ep == tcx.edge_event.constrained_edge.p {
				t.MarkConstrainedEdge(ep, eq)
				ot.MarkConstrainedEdge(ep, eq)
				s.Legalize(tcx, t)
				s.Legalize(tcx, ot)
			} else {
				// XXX: I think one of the triangles should be legalized here?
			}
		} else {
			o := Orient2d(eq, op, ep)
			t = s.NextFlipTriangle(tcx, o, t, ot, p, op)
			s.FlipEdgeEvent(tcx, ep, eq, t, p)
		}
	} else {
		newP := s.NextFlipPoint(ep, eq, ot, op)
		s.FlipScanEdgeEvent(tcx, ep, eq, t, ot, newP)
		s.EdgeEvent1(tcx, ep, eq, t, p)
	}
}

func (s *Sweep) NextFlipTriangle(tcx *SweepContext, o Orientation, t, ot *Triangle, p, op *Point) *Triangle {
	if o == CCW {
		// ot is not crossing edge after flip
		edge_index := ot.EdgeIndex(p, op)
		ot.delaunay_edge[edge_index] = true
		s.Legalize(tcx, ot)
		ot.ClearDelunayEdges()
		return t
	}

	// t is not crossing edge after flip
	edge_index := t.EdgeIndex(p, op)

	t.delaunay_edge[edge_index] = true
	s.Legalize(tcx, t)
	t.ClearDelunayEdges()
	return ot
}

func (s *Sweep) NextFlipPoint(ep, eq *Point, ot *Triangle, op *Point) *Point {
	o2d := Orient2d(eq, op, ep)
	if o2d == CW {
		// Right
		return ot.PointCCW(op)
	} else if o2d == CCW {
		// Left
		return ot.PointCW(op)
	}
	logger.Fatal("[Unsupported] Opposing point on constrained edge")
	return nil
}

func (s *Sweep) FlipScanEdgeEvent(tcx *SweepContext, ep, eq *Point, flip_triangle *Triangle,
	t *Triangle, p *Point) {
	ot := t.NeighborAcross(p)
	op := ot.OppositePoint(t, p)

	if InScanArea(eq, flip_triangle.PointCCW(eq), flip_triangle.PointCW(eq), op) {
		// flip with new edge op.eq
		s.FlipEdgeEvent(tcx, eq, op, ot, op)
		// TODO: Actually I just figured out that it should be possible to
		//       improve this by getting the next ot and op before the the above
		//       flip and continue the flipScanEdgeEvent here
		// set new ot and op here and loop back to inScanArea test
		// also need to set a new flip_triangle first
		// Turns out at first glance that this is somewhat complicated
		// so it will have to wait.
	} else {
		newP := s.NextFlipPoint(ep, eq, ot, op)
		s.FlipScanEdgeEvent(tcx, ep, eq, flip_triangle, ot, newP)
	}
}
