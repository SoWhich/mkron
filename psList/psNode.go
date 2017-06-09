package psList

import (
	"time"
	"strings"
	"errors"
)

type Ps struct {
	Comm string
	min []int
	hr []int
	day []int
	mon []int
	wkday []int
	Next *Ps
}

func isRight(time int, list []int)  bool {
	if len(list) != 0 {
		var x int
		for _ , x := range list {
			if x == time {
				break
			}
		}

		if x != time {
			return false }
	}

	return true
}

func (process Ps) IsTime(time time.Time) bool {
	// Luckily enough, Go manages time in the exact way the crontab
	// file recommends. For convenience, they're listed below
	// Minutes (0-59)
	// Hours (0-24)
	// Day (1-31)
	// Month (actually an enum, but indexed at 1)
	// Weekday (actually an enum, but indexed at 1)

	// for the workflow I designed (see outline.md) the 'percieved' minute
	// must be one more than the actual
	if isRight(time.Minute() + 1, process.min) &&
	   isRight(time.Hour(), process.hr) &&
	   isRight(time.Day(), process.day) &&
	   isRight(int(time.Month()), process.mon) &&
	   isRight(int(time.Weekday()), process.wkday) {
		return true
	} else {
		return false
	}
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
	var err error

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

	if bit[0] < '0' || bit[0] > '9' {
		err = errors.New("Parse Error")
	}

	for i := 0; i < len(bit); i++ {
		var temp int
		for ; bit[i] > '0' && bit[i] < '9'; i++ {
			temp *= 10
			temp += int(bit[i] -'0')
		}

		slice = append(slice, temp)

		switch bit[i] {
		case '-':
			if bit[i - 1] < '0' || bit[i - 1] > '9' {
				err = errors.New("Parse Error")
			}

			i++

			for ; bit[i] > '0' && bit[i] < '9'; i++ {
				temp *= 10
				temp += int(bit[i] -'0')
			}

			if slice[len(slice) -1] >= temp {
				err = errors.New("Parse Error")
			}

			for x := slice[len(slice) -1] + 1; x <= temp; x++ {
				slice = append(slice, x)
			}
		case ',':
		default:
			err = errors.New("Parse Error")
		}
	}

	for i := 0; i < len(slice); i++ {

		if slice[i] < min || slice[i] > max {
			err = errors.New("Parse Error")
		}

		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				err = errors.New("Parse Error")
			}
		}
	}

	return slice, err
}

func ParseLine(line string) (* Ps, error) {
	ret := new(Ps)
	var err error
	chunks := strings.Split(line, " ")

	if len(chunks) < 6 {
		err = errors.New("Parse Error")
	}

	ret.min, err = workBit(chunks[0], "min")
	ret.hr, err = workBit(chunks[1], "hr")
	ret.day, err = workBit(chunks[2], "day")
	ret.mon, err = workBit(chunks[3], "mon")
	ret.wkday, err = workBit(chunks[4], "wkday")
	ret.Comm = strings.Join(chunks[5:], " ")

	if err != nil {
		return nil, err
	}

	return ret, err
}
