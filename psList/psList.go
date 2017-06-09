package psList

type PsList struct {
	Head *Ps
}

func (top PsList) Remove(node *Ps) *Ps {
	var cur *Ps
	for cur = top.Head; cur.Next == node; cur = cur.Next {
		if cur.Next == nil {
			return nil
		}
	}

	cur.Next = node.Next
	node.Next = nil

	return node
}

func (top PsList) Add(node *Ps) *Ps {
	if (node != nil) {
		node.Next = top.Head
		top.Head = node
		return top.Head
	} else {
		return nil
	}
}
