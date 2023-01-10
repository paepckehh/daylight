# OVERVIEW

[paepche.de/daylight](https://paepcke.de/daylight)

Need to collaborate in an worldwide distributed team?
You (or your chatbot/lib) need to know local daylight 
(availibility) times? 

100% pure golang, lib/api has no external dependencies 
Example app has one. Use it as APP or api (see api.go).

Backend is a minimized, boiled down and heavy adapted static
fork of [github.com/sj14/astral](https://github.com/sj14/astral), 
who is afork of [github.com/sffjunkie/astral](https://github.com/sffjunkie/astral).
see pkg sun/sun.go for details (***ALL CREDITS GOES TO THE AUTHOR(S)***)

# INSTALL

```
go install paepcke.de/daylight/cmd/daylight@latest
```

# SHOWTIME (APP)

## Set location via gps coordinates.
```Shell
GSP_LAT=53.564432 GPS_LONG=9.95118 daylight 
Sunrise: 04:00:39 || Sunset 18:37:09 || Noon: 11:18:23 || Daylight: 14h36m30s

```

## Set location via nearest 3 letter Airport code.
```Shell
IATA=TXL daylight 
Sunrise: 04:00:39 || Sunset 18:37:09 || Noon: 11:18:23 || Daylight: 14h36m30s
```

## Ask if we have daylight @ Berlin 
```
IATA=BER daylight ask
true
```

## Ask if we have daylight @ Perth 
```
IATA=PER daylight ask
false
```

## Set Shell env variablesv via 3 letter Airport code
```
IATA=PER daylight unix 
#!/bin/sh
export GPS_LAT="-31.94"
export GPS_LONG="115.97"
export GPS_ELEVATION="0"
export GPS_SUN_RISE="23:04:38"
export GPS_SUN_SET="09:21:26"
export GPS_SUN_NOON="04:13:07"
export GPS_SUN_DAYLIGHT="10h16m48s"
```

# DOCS

[pkg.go.dev/paepcke.de/daylight](https://pkg.go.dev/paepcke.de/daylight)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
