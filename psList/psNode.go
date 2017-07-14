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

import (
	"errors"
	"strings"
	"time"
)

type Ps struct {
	Comm  string
	min   []int
	hr    []int
	day   []int
	mon   []int
	wkday []int
	Next  *Ps
}

func isRight(time int, list []int) bool {
	if len(list) != 0 {
		for _, x := range list {
			if x == time {
				return true
			}
		}

		return false
	}
	return true
}

func (process Ps) IsTime(time time.Time) bool {
	// Luckily enough, Go manages time in the exact way the crontab file
	// recommends. For convenience, they're listed below

	// Minutes (0-59)
	// Hours (0-24)
	// Day (1-31)
	// Month (actually an enum, but indexed at 1)
	// Weekday (actually an enum, but indexed at 1)

	return (isRight(time.Minute(), process.min) &&
		isRight(time.Hour(), process.hr) &&
		isRight(time.Day(), process.day) &&
		isRight(int(time.Month()), process.mon) &&
		isRight(int(time.Weekday()), process.wkday))
}

func (process Ps) Add(node *Ps) *Ps {
	if node == nil {
		return nil
	}
	node.Next = process.Next
	process.Next = node
	return node
}

func workBit(bit string, id string) ([]int, error) {
	var slice []int
	var max int
	var min int
	var err error = nil

	if bit == "*" {
		return slice, err
	}

	switch id {
	case "min":
		max = 59
		min = 0
	case "hr":
		max = 23
		min = 0
	case "day":
		max = 31
		min = 1
	case "mon":
		max = 12
		min = 1
	case "wkday":
		max = 7
		min = 0
	}

	for i := 0; i < len(bit); i++ {

		if bit[i] < '0' || bit[i] > '9' {
			err = errors.New("Parse Error 1")
			return slice, err
		}

		var temp int = 0
		for ; i < len(bit) && bit[i] >= '0' && bit[i] <= '9'; i++ {
			temp *= 10
			temp += int(bit[i] - '0')
		}

		slice = append(slice, temp)

		if i == len(bit) {
			break
		}

		switch bit[i] {
		case '-':

			i++

			if bit[i] < '0' || bit[i] > '9' {
				err = errors.New("Parse Error 2")
				return slice, err
			}

			temp = 0

			for ; i < len(bit) && bit[i] >= '0' && bit[i] <= '9'; i++ {
				temp *= 10
				temp += int(bit[i] - '0')
			}

			if slice[len(slice)-1] >= temp {
				err = errors.New("Parse Error 3")
				return slice, err
			}

			for x := slice[len(slice)-1] + 1; x <= temp; x++ {
				slice = append(slice, x)
			}
		case ',':
		default:
			err = errors.New("Parse Error 4")
			return slice, err
		}
	}

	for i := 0; i < len(slice); i++ {

		if slice[i] < min || slice[i] > max {
			err = errors.New("Parse Error 5")
			return slice, err
		}

		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				err = errors.New("Parse Error 6")
				return slice, err
			}
		}
	}

	return slice, err
}

func ParseLine(line string) (*Ps, error) {
	ret := new(Ps)
	var err error = nil
	chunks := strings.Split(line, " ")

	if len(chunks) < 6 {
		err = errors.New("Parse Error 7")
		return nil, err
	}

	ret.min, err = workBit(chunks[0], "min")
	if err != nil {
		return nil, err
	}

	ret.hr, err = workBit(chunks[1], "hr")
	if err != nil {
		return nil, err
	}

	ret.day, err = workBit(chunks[2], "day")
	if err != nil {
		return nil, err
	}

	ret.mon, err = workBit(chunks[3], "mon")
	if err != nil {
		return nil, err
	}

	ret.wkday, err = workBit(chunks[4], "wkday")
	if err != nil {
		return nil, err
	}
	ret.Comm = strings.Join(chunks[5:], " ")

	return ret, err
}
