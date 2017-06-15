package main

import (
	"bufio"
	"errors"
	"github.com/SoWhich/mkron/psList"
	"log"
	"os"
	"os/exec"
	"time"
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
				log.Printf("%s on line %d\n", err, lineNr+1)
				continue
			}
			allPs.Add(pcess)
		}
	}

	if allPs.Head == nil {
		log.Fatal(errors.New("empty/imparsible crontab"))
	}

	psStack := new(psList.PsList)

	for /*no SIGHUP/TERM/KILL*/ {

		now := time.Now().Local()
		for ps := allPs.Head; ps != nil; {
			if ps.IsTime(now) {
				psStack.Add(allPs.Remove(allPs.Head))
				ps = allPs.Head
			} else {
				ps = ps.Next
			}
		}

		for !psStack.IsEmpty() {
			cur := psStack.Remove(psStack.Head)
			ps := exec.Command("sh", "-c", cur.Comm)
			go ps.Start()
			allPs.Add(cur)
		}

		if /*Signal*/ false {
			/*response (in switch case) */
		}

		// sleep till next minute
		time.Sleep(time.Until(now.Truncate(time.Minute).Add(time.Minute)))
	}
}
