build:
	@echo "Building Project ..."
	go build
run:
	@echo "Running Service ..."
	go run main.go
clean:
	@echo "Cleaning Up ..."
	rm -rvf tweety
build_docker:
	@echo "Building Project ..."
	docker build -it tweety/tweety .
