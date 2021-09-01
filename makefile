current_dir = $(shell pwd)


windowsapp:
	@GOOS=windows CGO_ENABLED=1 GOARCH="amd64" CC="x86_64-w64-mingw32-gcc" go build ${current_dir}/main.go

macapp:
	@GOOS=darwin go build ${current_dir}/main.go

linuxapp:
	@GOOS=linux CGO_ENABLED=1 GOARCH="amd64" go build ${current_dir}/main.go

.PHONY:windowsapp,macapp,linuxapp