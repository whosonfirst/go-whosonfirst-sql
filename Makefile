GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -tags sqlite3 -mod $(GOMOD) -ldflags="$(LDFLAGS)" \
		-o bin/wof-sql-index cmd/wof-sql-index/main.go
