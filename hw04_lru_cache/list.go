package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Root ListItem
	len  int
}

func (l *list) Len() int { return l.len }

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.Root.Next
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.Root.Prev
}

func (l *list) insert(e, at *ListItem) *ListItem {
	l.Root.Next.Prev = &l.Root
	l.Root.Prev.Next = &l.Root
	e.Prev = at
	e.Next = at.Next
	e.Prev.Next = e
	e.Next.Prev = e
	l.Root.Next.Prev = nil
	l.Root.Prev.Next = nil

	l.len++
	return e
}

func (l *list) insertValue(v interface{}, at *ListItem) *ListItem {
	return l.insert(&ListItem{Value: v}, at)
}

func (l *list) PushFront(v interface{}) *ListItem {
	return l.insertValue(v, &l.Root)
}

func (l *list) PushBack(v interface{}) *ListItem {
	return l.insertValue(v, l.Root.Prev)
}

func (l *list) Remove(e *ListItem) {
	l.Root.Next.Prev = &l.Root
	l.Root.Prev.Next = &l.Root
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev
	e.Next = nil
	e.Prev = nil
	l.Root.Next.Prev = nil
	l.Root.Prev.Next = nil
	l.len--
}

func (l *list) move(e, at *ListItem) {
	if e == at {
		return
	}

	l.Root.Next.Prev = &l.Root
	l.Root.Prev.Next = &l.Root
	e.Prev.Next = e.Next
	e.Next.Prev = e.Prev

	e.Prev = at
	e.Next = at.Next
	e.Prev.Next = e
	e.Next.Prev = e
	l.Root.Next.Prev = nil
	l.Root.Prev.Next = nil
}

func (l *list) MoveToFront(e *ListItem) {
	if l.Root.Next == e {
		return
	}
	l.move(e, &l.Root)
}

func NewList() List {
	l := new(list)
	l.Root.Prev = &l.Root
	l.Root.Next = &l.Root
	return l
}
