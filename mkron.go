package main

import (
	"os/exec"
	"os"
	"time"
	"bufio"
	"errors"
	"log"
	"fmt"
	"github.com/SoWhich/mkron/psList"
)


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

	list := new(psList.PsList)
	var queue []*psList.Ps

	for x := range lines {
		pcess, err := psList.ParseLine(lines[x])
		if err != nil {
			fmt.Printf("%s on line %d\n", err, x + 1 )
			continue
		}
		list.Add(pcess)
	}

	if list.Head == nil {
		log.Fatal(errors.New("empty/imparsible crontab"))
	}

	for /*no SIGHUP/TERM/KILL*/ {

		if len(queue) > 0 {
			var cur *psList.Ps

			for len(queue) > 1 {
				cur = queue[0]
				fmt.Println(cur.Comm)
				q := exec.Command("/bin/sh", "-c", "\"", cur.Comm, "\"")
				fmt.Println(q)
				q.Run()
				queue = queue[1:]
			}

			cur = queue[0]
			fmt.Println(cur.Comm)
			q := exec.Command("/bin/sh", "-c", "\"", cur.Comm, "\"")
			fmt.Println(q)
			q.Run()
			queue = []*psList.Ps{}
		}

		now := time.Now().Local()
		fmt.Println(now)

		for  i := list.Head; i != nil; i = i.Next {
			if i.IsTime(now) {
				queue = append(queue, i)
			}
		}

		if /*timestampcheck*/ false {

			list.Head = nil
			queue = []*psList.Ps{}

			scanner := bufio.NewScanner(tab)

			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			for x := range lines {
				pcess, err := psList.ParseLine(lines[x])
				if err != nil {
					fmt.Printf("%s on line %d\n", err, x+1)
					continue
				}

				list.Add(pcess)
			}

			if list.Head == nil {
				log.Fatal(errors.New("empty/imparsible " +
					"crontab"))
			}
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}

		// sleep till next minute
		time.Sleep(time.Until(now.Truncate(time.Minute).Add(time.Minute)))
	}
}
