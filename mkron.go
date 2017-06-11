package main

import (
	"os/exec"
	"os"
	"time"
	"bufio"
	"errors"
	"log"
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

	err = tab.Close()

	if err != nil {
		log.Println(err)
	}

	list := new(psList.PsList)
	var queue []*psList.Ps

	for x := range lines {
		if lines[x][0] != '#' {
			pcess, err := psList.ParseLine(lines[x])
			if err != nil {
				log.Printf("%s on line %d\n", err, x + 1 )
				continue
			}
			list.Add(pcess)
		}
	}

	if list.Head == nil {
		log.Fatal(errors.New("empty/imparsible crontab"))
	}

	for /*no SIGHUP/TERM/KILL*/ {

		if len(queue) > 0 {
			var cur *psList.Ps

			for len(queue) > 1 {
				cur = queue[0]
				q := exec.Command("sh", "-c", cur.Comm)
				go q.Start()
				queue = queue[1:]
			}

			cur = queue[0]
			q := exec.Command("sh", "-c", cur.Comm,)
			go q.Start()
			queue = []*psList.Ps{}
		}

		now := time.Now().Local()
		for  i := list.Head; i != nil; i = i.Next {
			if i.IsTime(now) {
				queue = append(queue, i)
			}
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}

		// sleep till next minute
		time.Sleep(time.Until(now.Truncate(time.Minute).Add(time.Minute)))
	}
}
