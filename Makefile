# build a make file to run the go programa that is in cmd/mygit/main.go
# and create a binary file called mygit
# Usage:
# make build
# make run
# make clean
#

build:
	go build -o playground/mygit cmd/mygit/main.go

init: build
	./playground/mygit init

cat-file: init build
	./playground/mygit cat-file $(hash)
