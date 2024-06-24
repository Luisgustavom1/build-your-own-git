build:
	go build -o playground/mygit cmd/mygit/main.go

init: build
	./playground/mygit init $(dir)

cat-file: build
	./playground/mygit cat-file -$(flag) $(hash)
