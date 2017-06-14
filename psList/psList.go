package psList

type PsList struct {
	Head *Ps
}

func (top *PsList) Remove(node *Ps) *Ps {
	var cur *Ps

	if node != nil && node == top.Head {
		cur = top.Head
		top.Head = top.Head.Next
		return cur

	} else {
		for cur = top.Head; cur.Next == node; cur = cur.Next {
			if cur.Next == nil {
				return nil
			}
		}

		cur.Next = node.Next
		node.Next = nil

		return node
	}
}

func (top *PsList) Add(node *Ps) {
	if node != nil {
		node.Next = top.Head
		top.Head = node
	}
}

func (top *PsList) IsEmpty() bool {
	if top.Head == nil {
		return true
	} else {
		return false
	}
}
