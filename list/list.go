package list

type Node[T any] struct {
	item T
	next *Node[T]
	prev *Node[T]
}

type List[T any] struct {
	length int
	head   *Node[T]
	tail   *Node[T]
}

func New[T any](vs ...T) *List[T] {
	ll := &List[T]{}
	for _, v := range vs {
		ll.Put(ll.Len(), v)
	}

	return ll
}

func (ll *List[T]) Len() int {
	return ll.length
}

func (ll *List[T]) iter(rcv func(i int, n *Node[T]) bool) {
	for i, n := 0, ll.head; n.next != nil; i, n = i+1, n.next {
		if !rcv(i, n) {
			return
		}
	}
}

func (ll *List[T]) Iter(rcv func(i int, v T) bool) {
	ll.iter(func(i int, n *Node[T]) bool { return rcv(i, n.item) })
}

func (ll *List[T]) putAtHead(val T) {
	node := Node[T]{item: val}
	node.next = ll.head

	switch ll.head {
	case nil:
		ll.tail = &node
	default:
		ll.head.prev = &node
	}

	ll.head = &node
	ll.length++
}

func (ll *List[T]) putAtTail(val T) {
	node := Node[T]{item: val}
	node.prev = ll.tail

	switch ll.tail {
	case nil:
		ll.head = &node
	default:
		ll.tail.next = &node
	}

	ll.tail = &node
	ll.length++
}

func (ll *List[T]) Put(index int, val T) bool {
	if ll.length < index || index < 0 {
		return false
	}
	if index == 0 || ll.head == nil {
		ll.putAtHead(val)
		return true
	}
	if ll.length == index {
		ll.putAtTail(val)
		return true
	}

	ll.iter(
		func(i int, n *Node[T]) bool {
			if i < index {
				return true
			}

			node := Node[T]{item: val}
			node.prev = n.prev
			node.next = n

			n.prev = &node
			if node.prev != nil {
				node.prev.next = &node
			}

			return false
		},
	)

	ll.length++
	return true
}

func (ll *List[T]) delHead() (r T, ok bool) {
	if ll.head == nil {
		return r, ok
	}

	oldHead := ll.head
	ll.head = oldHead.next

	ll.length--
	return oldHead.item, true
}

func (ll *List[T]) delTail() (r T, ok bool) {
	if ll.tail == nil {
		return r, ok
	}

	oldTail := ll.tail
	ll.tail = oldTail.prev

	ll.length--
	return oldTail.item, true
}

func (ll *List[T]) Del(index int) (r T, ok bool) {
	if ll.length < index {
		return r, ok
	}
	if index == 0 || ll.head == nil {
		return ll.delHead()
	}
	if index == ll.length {
		return ll.delTail()
	}

	ll.iter(
		func(i int, n *Node[T]) bool {
			if i < index {
				return true
			}

			r = n.item
			ok = true

			if oldPrev := n.prev; oldPrev != nil {
				oldPrev.next = n.next
			}
			if oldNext := n.next; oldNext != nil {
				oldNext.prev = n.prev
			}

			return false
		},
	)

	ll.length--
	return r, ok
}

func (ll *List[T]) Get(index int) (r T, ok bool) {
	if index < 0 || index > ll.Len() {
		return r, ok
	}

	ll.Iter(
		func(i int, v T) bool {
			if i < index {
				return true
			}

			r = v
			ok = true
			return false
		},
	)

	return r, ok
}

func (ll *List[T]) IntoSlice() []T {
	return IntoSlice[T](ll)
}

func (ll *List[T]) IntoChan() <-chan T {
	return IntoChan(ll)
}

func (ll *List[T]) IntoGenerator() func() (T, bool) {
	return IntoGen(ll)
}
