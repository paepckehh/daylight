# OVERVIEW

[paepche.de/daylight](https://paepcke.de/daylight)

Set location, get information about sunrise, sunset, noon, daylight time. 

# INSTALL

```
go install paepcke.de/daylight/cmd/daylight@latest
```

# SHOWTIME

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

## Ask if we have daylight @ location
```
IATA=XFW daylight ask
true
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

# CONTRIBUTION

Yes, Please! PRs Welcome! 
