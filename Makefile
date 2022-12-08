All: build run

build: main.go
	go build -o main.exe main.go


run: main.exe
	main.exe
