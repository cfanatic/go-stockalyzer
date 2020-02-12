build:
	go build

deploy:
	@go build
	@./go-stockalyzer

clean:
	rm go-stockalyzer

all: deploy
