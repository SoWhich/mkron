package main

import (
	"github.com/sowhich/mkron/psList"
	"os/exec"
	"os"
	"time"
	"bufio"
)

func runNEmpty(q psList.PsQueue) {
	var cur *psList.Ps
	for !q.IsEmpty() {
		cur = q.Dequeue()
		exec.Command("/bin/sh", cur.Comm)
	}
}

func main() {
	// todo, allow user to specifiy special config file
	fname := "/etc/crontab"
	tab, _ := os.Open(fname)

	scanner := bufio.NewScanner(tab)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var list psList.PsList
	var queue psList.PsQueue

	for x := range lines {
		list.Add(psList.ParseLine(lines[x]))
	}

	if list.Head == nil {
		// Exit failure?
	}

	for /*no SIGHUP/TERM/KILL*/ {
		runNEmpty(queue)
		time := time.Now()
		for /* freaking make an iterator method MATT*/ i := list.Head; i != nil; i = i.Next {
			if i.IsTime(time) {
				queue.Enqueue(i)
			}
		}

		if /*timestampcheck*/ false {

			list.Head = nil
			queue.Front = nil
			queue.Back = nil

			scanner := bufio.NewScanner(tab)

			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			for x := range lines {
				list.Add(psList.ParseLine(lines[x]))
			}
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}
	}
}
