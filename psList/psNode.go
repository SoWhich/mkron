package psList

import (
	"time"
	"strings"
	"errors"
	"fmt"
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
		for _ , x := range list {
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

		if bit[0] < '0' || bit[0] > '9' {
			err = errors.New("Parse Error 1")
			return slice, err
		}

		var temp int = 0
		for ; i < len(bit) && bit[i] >= '0' && bit[i] <= '9'; i++ {
			temp *= 10
			temp += int(bit[i] -'0')
		}

		slice = append(slice, temp)
		fmt.Println(slice)

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

			for ;i < len(bit) && bit[i] >= '0' && bit[i] <= '9'; i++ {
				temp *= 10
				temp += int(bit[i] -'0')
			}

			if slice[len(slice) - 1] >= temp {
				err = errors.New("Parse Error 3")
				return slice, err
			}

			for x := slice[len(slice) - 1] + 1; x <= temp; x++ {
				slice = append(slice, x)
			}
			fmt.Println(slice)
		case ',':
		default:
			err = errors.New("Parse Error 4")
			fmt.Println(slice)
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

func ParseLine(line string) (* Ps, error) {
	ret := new(Ps)
	var err error = nil
	chunks := strings.Split(line, " ")

	if len(chunks) < 6 {
		err = errors.New("Parse Error 7")
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
