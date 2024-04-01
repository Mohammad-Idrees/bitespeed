genwire:
	rm -f wire_gen.go; wire;

build:
	go build main.go wire_gen.go

start:
	go run main.go wire_gen.go; 