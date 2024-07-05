package submitter

type OrderedSetNode struct {
	Key  interface{}
	Prev *OrderedSetNode
	Next *OrderedSetNode
}

type OrderedSet struct {
	end *OrderedSetNode
	set map[interface{}]*OrderedSetNode
}

func NewOrderedSet() *OrderedSet {
	end := OrderedSetNode{}
	end.Prev = &end
	end.Next = &end
	return &OrderedSet{
		end: &end,
		set: make(map[interface{}]*OrderedSetNode),
	}
}

func (o *OrderedSet) Add(key interface{}) {
	if _, ok := o.set[key]; ok {
		return
	}
	end := o.end
	curr := end.Prev
	curr.Next = &OrderedSetNode{Key: key, Prev: curr, Next: end}
	end.Prev = curr.Next
	o.set[key] = curr.Next
}

func (o *OrderedSet) Pop(last bool) interface{} {
	if o.end == nil {
		return nil
	}
	var key interface{}
	if last {
		if o.end.Prev != nil {
			key = o.end.Prev.Key
			delete(o.set, o.end.Prev.Key)
			o.end.Prev = o.end.Prev.Prev
			o.end.Prev.Next = o.end
		}
	} else {
		if o.end.Next != nil {
			key = o.end.Next.Key
			delete(o.set, o.end.Next.Key)
			o.end.Next = o.end.Next.Next
			o.end.Next.Prev = o.end
		}
	}
	return key
}

func (o *OrderedSet) Len() int {
	return len(o.set)
}
