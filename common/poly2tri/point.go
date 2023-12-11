package poly2tri

import "math"

type Point struct {
	x, y float64
	/// The edges this point constitutes an upper ending point
	edge_list []*Edge
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}
func (p *Point) SetPoint(x, y float64) {
	p.x = x
	p.y = y
}
func (p *Point) Add(v *Point) {
	p.x += v.x
	p.y += v.y
}

func (p *Point) Sub(v *Point) {
	p.x -= v.x
	p.y -= v.y
}
func (p *Point) Mul(a float64) {
	p.x *= a
	p.y *= a
}
func (p *Point) Negative() *Point {
	return NewPoint(-p.x, -p.y)
}
func (p *Point) SetZero() {
	p.x = 0.0
	p.y = 0.0
}
func (p *Point) Length() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

func (p *Point) Normalize() float64 {
	Length := p.Length()
	p.x /= Length
	p.y /= Length
	return Length
}

func PointCmp(a, b *Point) bool {
	if a.y < b.y {
		return true
	} else if a.y == b.y {
		// Make sure q is point with greater x value
		if a.x < b.x {
			return true
		}
	}
	return false
}
