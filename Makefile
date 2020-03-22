build:
	go build -o bin/stockalyzer cmd/stockalyzer/main.go

deploy:
	@go build -o bin/stockalyzer cmd/stockalyzer/main.go
	@./bin/stockalyzer -mode release

clean:
	rm -f bin/stockalyzer
	rm -f -r misc/output

all: deploy
