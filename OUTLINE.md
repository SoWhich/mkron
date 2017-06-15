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

3.	The list is read, and the events that need executing are put into a
	stack (which really technically is a list, but it's used as a stack)

4.	then the commands are popped off the stack and executed as goroutines

5.	After this, it checks if there's been a sigevent and responds
	accordingly (STILL A TODO)

6.	And sleeps till the next minute

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
