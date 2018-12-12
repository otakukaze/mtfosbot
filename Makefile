
.PHONY: clean build

build:
	go-bindata -pkg schema -ignore .git -o module/schema/static.go schema/
	GOOS=linux go build -o mtfosbot -ldflags "-s -w" .

clean:
	rm -rf mtfosbot && go clean
