package psList

/*
	MKron Copyright (C) 2017 Matthew Kunjummen

	This file is part of MKron.

	MKron is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	MKron is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with MKron.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
	psList is a simple list struct designed for crontab lines

	Each Node contains slice regarding times when it should run, a string
	with the associated command, and a pointer to the next node. it also
	as an associated function to add a node after itself

	the List itself only has a head pointer, a function to add a node at
	the head, and a function to remove a node at any point in the list
	and a simple method to check if it is empty or not

	Additionally, the package has the functions associated with parsing the
	files to put them into the nodes and to verify if the current time is
	the appropriate time to run the command.
*/

type PsList struct {
	Head *Ps
}

func (top *PsList) Remove(node *Ps) *Ps {

	if node == nil {

	} else if node == top.Head {
		top.Head = top.Head.Next

	} else {
		var cur *Ps
		for cur = top.Head; cur.Next == node; cur = cur.Next {
			if cur.Next == nil {
				return nil
			}
		}

		cur.Next = node.Next
		node.Next = nil
	}

	return node
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
