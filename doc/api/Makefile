REDOC := $(shell command -v redoc-cli 2> /dev/null)

all: index.html

index.html: api.yaml
ifndef REDOC
    $(error "redoc-cli is not available, please install via npm install -g redoc-cli (run as root)")
endif
	redoc-cli bundle -o villaswebapi.html api.yaml

clean:
	rm -rf index.html

.PHONY: docs index.html clean
