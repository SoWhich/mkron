package main

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

import (
	"bufio"
	"errors"
	"flag"
	"github.com/SoWhich/mkron/psList"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	fname := flag.String("f", "/etc/crontab", "the crontab file location")
	flag.Parse()
	tab, err := os.Open(*fname)

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
				cur := ps.Next
				psStack.Add(allPs.Remove(ps))
				ps = cur
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
