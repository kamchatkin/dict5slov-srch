build:
	cp dict.txt search/dict.txt
	go build .
	rm search/dict.txt

build-linux:
	cp dict.txt search/dict.txt
	env GOOS=linux GOARCH=amd64 go build -o wordle-hack-linux .
	rm search/dict.txt
