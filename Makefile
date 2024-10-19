# Makefile for building and installing a Go program

# Set the name of your Go program
PROGRAM = delete-redis-keys

# Build the program
build:
	go build -o $(PROGRAM) -ldflags="-s" main.go

# Install the program to /usr/local/bin
install: build
	sudo cp $(PROGRAM) /usr/local/bin

# Clean up temporary files
clean:
	rm -f $(PROGRAM)
