package poly2tri

import "github.com/gorustyt/go-assimp/common/logger"

type Node struct {
	point    *Point
	triangle *Triangle

	next *Node
	prev *Node

	value float64
}

func NewNode(p *Point) *Node {
	n := &Node{point: p, value: p.x}
	return n
}

func NewNodeWithTriangle(p *Point, t *Triangle) *Node {
	n := &Node{point: p, value: p.x, triangle: t}
	return n
}

type AdvancingFront struct {
	head_, tail_, search_node_ *Node
}

func NewAdvancingFront(head, tail *Node) *AdvancingFront {
	return &AdvancingFront{head_: head, tail_: tail}
}
func (ad *AdvancingFront) Head() *Node {
	return ad.head_
}
func (ad *AdvancingFront) Set_head(node *Node) {
	ad.head_ = node
}

func (ad *AdvancingFront) Tail() *Node {
	return ad.tail_
}
func (ad *AdvancingFront) Set_tail(node *Node) {
	ad.tail_ = node
}

func (ad *AdvancingFront) Search() *Node {
	return ad.search_node_
}

func (ad *AdvancingFront) Set_search(node *Node) {
	ad.search_node_ = node
}

func (ad *AdvancingFront) AdvancingFront(head, tail *Node) {
	ad.head_ = head
	ad.tail_ = tail
	ad.search_node_ = head
}

func (ad *AdvancingFront) LocateNode(x float64) *Node {
	node := ad.search_node_
	if x < node.value {
		node = node.prev
		for node != nil {
			if x >= node.value {
				ad.search_node_ = node
				return node
			}
			node = node.prev
		}
	} else {
		for node != nil {
			if x < node.value {
				ad.search_node_ = node.prev
				return node.prev
			}
			node = node.next
		}
	}
	return nil
}

func (ad *AdvancingFront) FindSearchNode(x float64) *Node {
	// suppress compiler warnings "unused parameter 'x'"
	// TODO: implement BST index
	return ad.search_node_
}

func (ad *AdvancingFront) LocatePoint(point *Point) *Node {
	px := point.x
	node := ad.FindSearchNode(px)
	nx := node.point.x

	if px == nx {
		if point != node.point {
			// We might have two nodes with same x value for a short time
			if point == node.prev.point {
				node = node.prev
			} else if point == node.next.point {
				node = node.next
			} else {
				logger.FatalF("LocatePoint error")
			}
		}
	} else if px < nx {
		node = node.prev
		for node != nil {
			if point == node.point {
				break
			}
			node = node.prev
		}
	} else {
		node = node.next
		for node != nil {
			if point == node.point {
				break
			}
			node = node.next
		}
	}
	if node != nil {
		ad.search_node_ = node
	}
	return node
}
