package poly2tri

import (
	"container/list"
	"sort"
)

const (
	kAlpha = 0.3
)

type Basin struct {
	left_node    *Node
	bottom_node  *Node
	right_node   *Node
	width        float64
	left_highest bool
}

func NewBasin() *Basin {
	return &Basin{}
}
func (b *Basin) Clear() {
	b.left_node = nil
	b.bottom_node = nil
	b.right_node = nil
	b.width = 0.0
	b.left_highest = false
}

type EdgeEvent struct {
	constrained_edge *Edge
	right            bool
}
type SweepContext struct {
	edge_list []*Edge

	basin      Basin
	edge_event EdgeEvent

	triangles_ []*Triangle
	map_       list.List
	points_    []*Point

	// Advancing front
	front_ *AdvancingFront
	// head point used with advancing front
	head_ *Point
	// tail point used with advancing front
	tail_ *Point

	af_head_, af_middle_, af_tail_ *Node
}

func NewSweepContext(polyline []*Point) *SweepContext {
	s := &SweepContext{}
	s.InitEdges(polyline)
	return s
}
func (s *SweepContext) Front() *AdvancingFront {
	return s.front_
}

func (s *SweepContext) point_count() int {
	return len(s.points_)
}

func (s *SweepContext) set_head(p1 *Point) {
	s.head_ = p1
}

func (s *SweepContext) Head() *Point {
	return s.head_
}

func (s *SweepContext) Set_tail(p1 *Point) {
	s.tail_ = p1
}

func (s *SweepContext) Tail() *Point {
	return s.tail_
}
func (s *SweepContext) InitEdges(polyline []*Point) {
	num_points := len(polyline)
	for i := 0; i < num_points; i++ {
		j := 0
		if i < num_points-1 {
			j = i + 1
		}
		s.edge_list = append(s.edge_list, NewEdge(polyline[i], polyline[j]))
	}
}

func (s *SweepContext) AddHole(polyline []*Point) {
	s.InitEdges(polyline)
	for i := 0; i < len(polyline); i++ {
		s.points_ = append(s.points_, polyline[i])
	}
}

func (s *SweepContext) AddPoint(point *Point) {
	s.points_ = append(s.points_, point)
}

func (s *SweepContext) GetTriangles() []*Triangle {
	return s.triangles_
}

func (s *SweepContext) GetMap() list.List {
	return s.map_
}

func (s *SweepContext) InitTriangulation() {
	xmax := s.points_[0].x
	xmin := s.points_[0].x
	ymax := s.points_[0].y
	ymin := s.points_[0].y

	// Calculate bounds.
	for i := 0; i < len(s.points_); i++ {
		p := s.points_[i]
		if p.x > xmax {
			xmax = p.x
		}

		if p.x < xmin {
			xmin = p.x
		}

		if p.y > ymax {
			ymax = p.y
		}

		if p.y < ymin {
			ymin = p.y
		}

	}

	dx := kAlpha * (xmax - xmin)
	dy := kAlpha * (ymax - ymin)
	s.head_ = NewPoint(xmax+dx, ymin-dy)
	s.tail_ = NewPoint(xmin-dx, ymin-dy)

	// Sort points along y-axis
	sort.Slice(s.points_, func(i, j int) bool {
		return PointCmp(s.points_[i], s.points_[j])
	})
}

func (s *SweepContext) GetPoint(index int) *Point {
	return s.points_[index]
}

func (s *SweepContext) AddToMap(triangle *Triangle) {
	s.map_.PushBack(triangle)
}

func (s *SweepContext) LocateNode(point *Point) *Node {
	// TODO implement search tree
	return s.front_.LocateNode(point.x)
}

func (s *SweepContext) CreateAdvancingFront(nodes []*Node) {
	// Initial triangle
	triangle := NewTriangle(s.points_[0], s.tail_, s.head_)

	s.map_.PushBack(triangle)

	s.af_head_ = NewNodeWithTriangle(triangle.GetPoint(1), triangle)
	s.af_middle_ = NewNodeWithTriangle(triangle.GetPoint(0), triangle)
	s.af_tail_ = NewNode(triangle.GetPoint(2))
	s.front_ = NewAdvancingFront(s.af_head_, s.af_tail_)

	// TODO: More intuitive if head is middles next and not previous?
	//       so swap head and tail
	s.af_head_.next = s.af_middle_
	s.af_middle_.next = s.af_tail_
	s.af_middle_.prev = s.af_head_
	s.af_tail_.prev = s.af_middle_
}

func (s *SweepContext) MapTriangleToNodes(t *Triangle) {
	for i := 0; i < 3; i++ {
		if t.GetNeighbor(i) == nil {
			n := s.front_.LocatePoint(t.PointCW(t.GetPoint(i)))
			if n != nil {
				n.triangle = t
			}
		}
	}
}

func (s *SweepContext) RemoveFromMap(triangle *Triangle) {
	e := s.map_.Front()
	for ; e != nil; e = e.Next() {
		if e.Value.(*Triangle) == triangle {
			s.map_.Remove(e)
		}
	}
}

func (s *SweepContext) MeshClean(triangle *Triangle) {
	var triangles []*Triangle
	triangles = append(triangles, triangle)
	for len(triangles) > 0 {
		t := triangles[len(triangles)-1]
		triangles = triangles[:len(triangles)-1]
		if t != nil && !t.IsInterior() {
			t.SetInterior(true)
			triangles = append(triangles, t)
			for i := 0; i < 3; i++ {
				if !t.constrained_edge[i] {
					triangles = append(triangles, t.GetNeighbor(i))
				}
			}
		}
	}
}
