
gonic2:  main_server.go
	go build .

gonic2.exe:  main_server.go
	go build -o $@  .

upd:
	go get -u -v
	go mod tidy

clean:
	@rm gonic2.exe
	rm gonic2

