package poly2tri

import "container/list"

type CDT struct {
	sweep_context_ *SweepContext
	sweep_         *Sweep
}

func NewCDT(polyline []*Point) *CDT {
	cdt := &CDT{}
	cdt.sweep_context_ = NewSweepContext(polyline)
	cdt.sweep_ = NewSweep()
	return cdt
}
func (cdt *CDT) AddHole(polyline []*Point) {
	cdt.sweep_context_.AddHole(polyline)
}

func (cdt *CDT) AddPoint(point *Point) {
	cdt.sweep_context_.AddPoint(point)
}

func (cdt *CDT) Triangulate() {
	cdt.sweep_.Triangulate(cdt.sweep_context_)
}

func (cdt *CDT) GetTriangles() []*Triangle {
	return cdt.sweep_context_.GetTriangles()
}

func (cdt *CDT) GetMap() list.List {
	return cdt.sweep_context_.GetMap()
}
