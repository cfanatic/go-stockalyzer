build:
	go build -o bin/stockalyzer cmd/stockalyzer/main.go

deploy:
	@go build -o bin/stockalyzer cmd/stockalyzer/main.go
	@./bin/stockalyzer

clean:
	rm -f bin/stockalyzer
	rm -f bin/output.png

all: deploy
