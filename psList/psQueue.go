package psList

type PsQueue struct {
	Front *Ps
	Back *Ps
}

func (q PsQueue) Enqueue(node *Ps) {
	node.Next = q.Back
	q.Back = node.Next
	if q.Front == nil {
		q.Back = q.Front
	}
}

func (q PsQueue) Dequeue() *Ps {
	var cur *Ps
	var holder *Ps
	for cur = q.Back; cur.Next != q.Front; cur = cur.Next {}
	holder = q.Front
	q.Front = cur.Next
	cur.Next = nil
	return holder
}

func (q PsQueue) IsEmpty() bool {
	if q.Back == nil {
		return true
	} else {
		return false
	}
}
