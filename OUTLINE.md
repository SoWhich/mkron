MKron
=====
-------------------------------------------------------------------------------

>A cron daemon written by **M**atthew **K**unjummen

MKron is written in go and designed to be simple to write and (roughly)
POSIX compatible.

## Overview ##

### Software Level ###

1.	MKron reads the crontab file from /etc/crontab

2.	The lines in the file are parsed, and each command is put into a list
	type data structure

3.	Every minute (the least gradation in the file) it goes through a queue
	of commands that are slated to be started this minute and then executes
	them, and moves them back to the normal list

4.	Then, it goes through the list and adds to the queue the events that
	need to be started the _next_ minute

5.	After this, it checks if there's been a sigevent and responds
	accordingly

6.	Then checks if there's been a modification of the crontab file

7.	And sleeps till the next minute

### Human Level ###

#### Crontab File ####

	* * * * * command (unquoted) to be run
	| | | | |
	| | | | + Day of the week (0-6 (7), where 0 (or 7) is Sunday)
	| | | +-- Month of the year (1-12)
	| | +---- Day of the Month (1-31)
	| +------ Hour (0-23)
	+-------- Minute (0-59)

MKron accepts lists and durations of the values.

* `0 0 1 1,5 * comm`
runs 'comm' at midnight of the first of every month in January and May

* `0 0 * * 1-5 comm`
runs 'comm' at midnight of every weekday

* `* * 1 1 * comm`
runs 'comm' at every minute of January First
