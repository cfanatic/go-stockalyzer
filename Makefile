build:
	go build -o bin/stockalyzer cmd/stockalyzer/main.go

deploy:
	@go build -o bin/stockalyzer cmd/stockalyzer/main.go
	@./bin/stockalyzer

clean:
	rm bin/stockalyzer

all: deploy
