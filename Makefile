.PHONY: test install clean

test:
	go test ./...

install:
	go install github.com/fujiwara/go-zabbix-get

dist/:
	goreleaser build --snapshot --rm-dist

clean:
	rm -fr dist/
