GIT_REV := $(shell git rev-parse --short HEAD)

all:
	sed -i.bak -e "s/HEAD/$(GIT_REV)/" revision.go
	gox -os="linux darwin" -arch="amd64 386" -output "pkg/{{.OS}}_{{.Arch}}/{{.Dir}}"
	git checkout revision.go
	rm revision.go.bak
