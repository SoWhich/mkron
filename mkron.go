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

	allPs := new(psList.PsList)
	for lineNr := range lines {
		if lines[lineNr][0] != '#' {
			pcess, err := psList.ParseLine(lines[lineNr])
			if err != nil {
				log.Printf("%s on line %d\n", err, lineNr + 1 )
				continue
			}
			allPs.Add(pcess)
		}
	}

	if allPs.Head == nil {
		log.Fatal(errors.New("empty/imparsible crontab"))
	}

	var queuedPs []*psList.Ps
	for /*no SIGHUP/TERM/KILL*/ {

		if len(queuedPs) > 0 {
			var cur *psList.Ps

			for len(queuedPs) > 1 {
				cur = queuedPs[0]
				ps := exec.Command("sh", "-c", cur.Comm)
				go ps.Start()
				queuedPs = queuedPs[1:]
			}

			cur = queuedPs[0]
			ps := exec.Command("sh", "-c", cur.Comm,)
			go ps.Start()
			queuedPs = []*psList.Ps{}
		}

		now := time.Now().Local()
		for ps := allPs.Head; ps != nil; ps = ps.Next {
			if ps.IsTime(now) {
				queuedPs = append(queuedPs, ps)
			}
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}

		// sleep till next minute
		time.Sleep(time.Until(now.Truncate(time.Minute).Add(time.Minute)))
	}
}
