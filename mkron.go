package main

import (
	"github.com/SoWhich/mkron/psList"
	"os/exec"
	"os"
	"time"
	"bufio"
	"errors"
	"log"
	"fmt"
)

func runNEmpty(q psList.PsQueue) {
	var cur *psList.Ps
	for !q.IsEmpty() {
		cur = q.Dequeue()
		exec.Command("/bin/sh", cur.Comm)
	}
}

func main() {
	// todo, allow user to specifiy special cron file
	fname := "/etc/crontab"
	tab, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}

	// read the crontab file and load the lines into a slice
	scanner := bufio.NewScanner(tab)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var list psList.PsList
	var queue psList.PsQueue

	for x := range lines {
		pcess, err := psList.ParseLine(lines[x])
		if err != nil {
			fmt.Println(err)
			continue
		}

		list.Add(pcess)
	}

	if list.Head == nil {
		log.Fatal(errors.New("empty or imparsible crontab"))
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
				pcess, err := psList.ParseLine(lines[x])
				if err != nil {
					fmt.Println(err)
					continue
				}

				list.Add(pcess)
			}

			if list.Head == nil {
				log.Fatal(errors.New("empty or imparsible crontab"))
			}
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}
	}
}
