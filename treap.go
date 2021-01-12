package safe_treap

import "errors"

// treap structure to define the root node
//
// The root is set by user
type Treap struct {
	handle  *Handle
	root    *Node
}

// node is the recursive data structure that defines a persistent treap
//
// The zero value is ready to use
type Node struct {
	Weight int
	Key, Item  interface{}
	Left, Right *Node
}


// Handle performs purely functional transformations on a treap.
type Handle struct {
	CompareWeights, CompareKeys Comparator
}

func NewTreap(h *Handle) (*Treap, error) {
	if h == nil {
		return nil, errors.New("comparator is nil")
	}
	treap :=  &Treap{handle: h, root: nil}

	return treap, nil
}

// Get an element by key.  Returns nil if the key is not in the treap.
// O(log n) if the treap is balanced (i.e. has uniformly distributed weights).
func (t *Treap) Get(n *Node, key interface{}) (v interface{}, found bool) {
	if n, found = t.GetNode(n, key); found {
		v = n.Item
	}
	return
}

// GetNode returns the subtree whose root has the specified key.  This is equivalent to
// Get, but returns a full node.
func (t *Treap) GetNode(n *Node, key interface{}) (*Node, bool) {
	if n == nil {
		return nil, false
	}

	switch comp := t.handle.CompareKeys(key, n.Key); {
	case comp < 0:
		return t.GetNode(n.Left, key)
	case comp > 0:
		return t.GetNode(n.Right, key)
	default:
		return n, true
	}
}

func (t *Treap) Min() interface{} {
	n := t.root
	if n == nil {
		return nil
	}
	for n.Left != nil {
		n = n.Left
	}
	return n.Item
}

func (t *Treap) Max() interface{} {
	n := t.root
	if n == nil {
		return nil
	}
	for n.Right != nil {
		n = n.Right
	}
	return n.Item
}

// Insert an element into the treap, returning false if the element is already present.
//
// O(log n) if the treap is balanced (see Get).
func (t *Treap) Insert(n *Node, key, val interface{}, weight int) (new *Node, ok bool) {
	return t.upsert(n, key, val, weight, true, false, nil)
}

func (t *Treap) upsert(n *Node, k, v interface{}, w int, create, update bool, fn func(*Node) bool) (res *Node, created bool) {
	if n == nil {
		if create {
			created = true
			res = &Node{Weight: w, Key: k, Item: v}
		}

		return
	}

	switch t.handle.CompareKeys(k, n.Key) {
	case -1:
		// use res as temp variable to avoid extra allocation
		if res, created = t.upsert(n.Left, k, v, w, create, update, fn); res == nil {
			return
		}

		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Item:   n.Item,
			Left:   res,
			Right:  n.Right,
		}
	case 1:
		// use res as temp variable to avoid extra allocation
		if res, created = t.upsert(n.Right, k, v, w, create, update, fn); res == nil {
			return
		}

		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Item:   n.Item,
			Left:   n.Left,
			Right:  res,
		}
	default:
		if !update { // insert only (no upsert)
			return
		}

		if fn != nil && !fn(n) { // InsertIf decided to ignore
			res = n
			return
		}

		res = &Node{
			Weight: w,
			Key:    n.Key,
			Item:   n.Item,
			Left:   n.Left,
			Right:  n.Right,
		}

		if create { // not SetWeight
			res.Item = v // upsert; set new value.
		}
	}

	if res.Left != nil && t.handle.CompareWeights(res.Left.Weight, res.Weight) < 0 {
		res = t.leftRotation(res)
	} else if res.Right != nil && t.handle.CompareWeights(res.Right.Weight, res.Weight) < 0 {
		res = t.rightRotation(res)
	}

	return
}

func (t *Treap) leftRotation(n *Node) *Node {
	return &Node{
		Weight: n.Left.Weight,
		Key:    n.Left.Key,
		Item:   n.Left.Item,
		Left:   n.Left.Left,
		Right: &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Item:  n.Item,
			Left:   n.Left.Right,
			Right:  n.Right,
		},
	}
}

func (t *Treap) rightRotation(n *Node) *Node {
	return &Node{
		Weight: n.Right.Weight,
		Key:    n.Right.Key,
		Item:   n.Right.Item,
		Left: &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Item:   n.Item,
			Left:   n.Left,
			Right:  n.Right.Left,
		},
		Right: n.Right.Right,
	}
}