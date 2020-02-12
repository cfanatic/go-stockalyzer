build:
	go build -o bin/go-stockalyzer cmd/main.go

deploy:
	@go build -o bin/go-stockalyzer cmd/main.go
	@./bin/go-stockalyzer

clean:
	rm bin/go-stockalyzer

all: deploy
