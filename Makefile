build:
	go build -o bin/go-stockalyzer

deploy:
	@go build -o bin/go-stockalyzer
	@./bin/go-stockalyzer

clean:
	rm bin/go-stockalyzer

all: deploy
