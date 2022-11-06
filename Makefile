build:
	cp dict.txt search/dict.txt
	go build .
	rm search/dict.txt
