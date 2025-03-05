all:
	IATA=BER go run main.go

examples:
	IATA=BER go run main.go
	IATA=HAM go run main.go
	IATA=TLS go run main.go
	GPS_LAT=53.564432 GPS_LONG=9.95118 go run main.go 
	GPS_LAT=53.564432 GPS_LONG=9.95118 go run main.go unix
	IATA=TLS go run main.go ask
	IATA=TLS go run main.go unix

check: 
	gofmt -w -s .
	go vet .
	staticcheck
	golangci-lint run
