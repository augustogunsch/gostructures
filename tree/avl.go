package tree

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Node[T constraints.Ordered] struct {
	Value  T
	Left   *Node[T]
	Right  *Node[T]
	height int
}

func Create[T constraints.Ordered](values ...T) *Node[T] {
	var head *Node[T] = nil

	for _, value := range values {
		head = head.Insert(value)
	}

	return head
}

func (node *Node[T]) getHeight() int {
	if node == nil {
		return -1
	}
	return node.height
}

func (node *Node[T]) updateHeight() {
	node.height = max(node.Left.getHeight(), node.Right.getHeight()) + 1
}

func (node *Node[T]) getBalance() int {
	return node.Left.getHeight() - node.Right.getHeight()
}

func (node *Node[T]) Insert(value T) *Node[T] {
	if node == nil {
		return &Node[T]{Value: value}
	}

	if value <= node.Value {
		if node.Left == nil {
			node.Left = &Node[T]{Value: value}
		} else {
			node.Left = node.Left.Insert(value)
		}
	} else {
		if node.Right == nil {
			node.Right = &Node[T]{Value: value}
		} else {
			node.Right = node.Right.Insert(value)
		}
	}

	node = node.rebalance()

	return node
}

func (node *Node[T]) Find(value T) (*Node[T], error) {
	if node == nil {
		return nil, fmt.Errorf("no such node with value: %v", value)
	}

	if value < node.Value {
		return node.Left.Find(value)
	}

	if value > node.Value {
		return node.Right.Find(value)
	}

	// value == node.Value
	return node, nil
}

func (node *Node[T]) popLeftBiggest() *Node[T] {
	prev := node
	node = node.Left

	if node == nil {
		return nil
	}

	if node.Right == nil {
		prev.Left = node.Left
		node.Left = nil
		return node
	}

	for node.Right != nil {
		prev = node
		node = node.Right
	}

	prev.Right = node.Left
	node.Left = nil

	return node
}

func (node *Node[T]) popRightSmallest() *Node[T] {
	prev := node
	node = node.Right

	if node == nil {
		return nil
	}

	if node.Left == nil {
		prev.Right = node.Right
		node.Right = nil
		return node
	}

	for node.Left != nil {
		prev = node
		node = node.Left
	}

	prev.Left = node.Right
	node.Right = nil

	return node
}

// Returns the new subtree.
func (node *Node[T]) Remove() *Node[T] {
	balance := node.getBalance()
	var subs *Node[T]

	switch {
	case balance > 0:
		subs = node.popLeftBiggest()
	case balance < 0:
		subs = node.popRightSmallest()
	case node.Left != nil:
		subs = node.popLeftBiggest()
	case node.Right != nil:
		subs = node.popRightSmallest()
	}

	if subs != nil {
		if subs != node.Left {
			subs.Left = node.Left
		}
		if subs != node.Right {
			subs.Right = node.Right
		}
		node.Left, node.Right = nil, nil
		subs = subs.rebalance()
	}

	return subs
}

// Returns the new subtree.
func (node *Node[T]) RemoveWith(value T) (*Node[T], error) {
	if node == nil {
		return nil, fmt.Errorf("no such node with value: %v", value)
	}

	if value < node.Value {
		left, err := node.Left.RemoveWith(value)
		node.Left = left
		node = node.rebalance()
		return node, err
	}

	if value > node.Value {
		right, err := node.Right.RemoveWith(value)
		node.Right = right
		node = node.rebalance()
		return node, err
	}

	// value == node.Value
	return node.Remove(), nil
}

func (node *Node[T]) rebalance() *Node[T] {
	node.updateHeight()
	balance := node.getBalance()

	switch {
	// Heavy to the left
	case balance > 1:
		if node.Left.getBalance() <= -1 {
			node.Left.rotateLeft()
		}
		node = node.rotateRight()

	// Heavy to the right
	case balance < -1:
		if node.Right.getBalance() >= 1 {
			node.Right.rotateRight()
		}
		node = node.rotateLeft()
	}

	return node
}

func (node *Node[T]) rotateLeft() *Node[T] {
	right := node.Right

	node.Right = right.Left
	node.updateHeight()

	right.Left = node
	right.updateHeight()

	return right
}

func (node *Node[T]) rotateRight() *Node[T] {
	left := node.Left

	node.Left = left.Right
	node.updateHeight()

	left.Right = node
	left.updateHeight()

	return left
}

func (node Node[T]) String() string {
	switch {
	case node.Left != nil && node.Right != nil:
		return fmt.Sprintf("(%s %v %s)",
			node.Left.String(),
			node.Value,
			node.Right.String(),
		)
	case node.Left != nil && node.Right == nil:
		return fmt.Sprintf("(%s %v)",
			node.Left.String(),
			node.Value,
		)
	case node.Left == nil && node.Right != nil:
		return fmt.Sprintf("(%v %s)",
			node.Value,
			node.Right.String(),
		)
	default:
		return fmt.Sprintf("(%v)", node.Value)
	}
}
