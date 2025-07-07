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
	First *ListItem
	Last  *ListItem
}

func NewList() List {
	return &list{
		First: nil,
		Last:  nil,
	}
}

func (l list) Len() int {
	if l.First == nil {
		return 0
	}
	counter := 1
	index := true
	elem := l.First
	for index {
		if elem.Next != nil {
			counter++
			elem = elem.Next
		} else {
			index = false
		}
	}
	return counter
}

func (l list) Front() *ListItem {
	return l.First
}

func (l list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	added := ListItem{
		Value: v,
		Next:  l.First,
	}
	l.First = &added
	switch {
	case l.Last == nil:
		l.Last = l.First
	case l.First.Next == l.Last:
		l.Last.Prev = l.First
	}
	return l.First
}

func (l *list) PushBack(v interface{}) *ListItem {
	added := ListItem{
		Value: v,
	}
	if l.Last == nil {
		l.Last = &added
		l.Last.Prev = l.First
	} else {
		added.Prev = l.Last
		l.Last = &added
		l.Last.Prev.Next = l.Last
	}

	if l.First == nil {
		l.First = l.Last
	}
	return l.Last
}

func (l *list) Remove(i *ListItem) {
	next := i.Next
	prev := i.Prev
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
	if l.Last == i {
		l.Last = prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.First != i {
		nextToGoal := i.Next
		prevToGoal := i.Prev
		exFirst := l.First
		if i == l.Last {
			l.Last = prevToGoal
		} else {
			nextToGoal.Prev = prevToGoal
		}
		if prevToGoal != nil {
			prevToGoal.Next = nextToGoal
		}
		i.Prev = nil
		i.Next = exFirst
		l.First = i
		l.First.Next.Prev = l.First
	}
}
