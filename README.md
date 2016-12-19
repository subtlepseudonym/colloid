### Colloid Server Sandbox ###

This repo is for a motley collection of web code and project files for screwing around on my home server.

#### Requirements ####

+ Go
	- <a href="https://golang.org/doc/install">Install Golang</a>


#### Project Files ####

+ main/main.go
	- On startup, prints the public ip and port that it's listening on


#### TODO ####

+ Organization
	- Separate into packages
	- Static IPs
	- Separate repo for front end
	- Docker?
+ Transition page serving to separate front end
	- Angular 2!
	- Stop serving files, start serving REST endpoints
+ Do something with go-cron
	- That shit's cool
	- Twilio?
+ Implement log writing
	- Postgresql
	- Log viewer form
	- Actual log viewer
	- Simplified log submission interface
	- <em>Need</em> a key of some kind for submissions (rsa?)
+ Stretch
	- New box / Spanner / RAID
	- Talk to atlas (very stretch)
	- Simple user accts (not so stretch)
	- IRC-like stuff
