NAME := shunin

SRCS := $(shell find . -type f -name '*.go')

all: bin/$(NAME) bin/mcrcon

mcrcon:
	git clone https://github.com/Tiiffi/mcrcon.git
	cd mcrcon; make
	mkdir -p bin

bin/mcrcon: mcrcon
	cp mcrcon/mcrcon bin/mcrcon

bin/$(NAME): $(SRCS)
	go build -o $@ ./shunin

clean:
	rm -rf bin etc

