GIT_VER := $(shell git describe --tags)

all:
	go get github.com/fujiwara/go-zabbix-get

binary:
	gox -os="linux darwin" -arch="amd64" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-X main.Version ${GIT_VER}"
	gox -os="linux" -arch="arm" -output "pkg/{{.Dir}}-${GIT_VER}-{{.OS}}-{{.Arch}}" -ldflags "-X main.Version ${GIT_VER}"
	cd pkg && find . -name "*${GIT_VER}*" -type f -exec zip {}.zip {} \;

test:
	cd zabbix && go test

clean:
	rm -f pkg/*
