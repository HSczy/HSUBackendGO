current_dir = $(shell pwd)


windowsapp:
	@GOOS=windows go build ${current_dir}/main.go

macapp:
	@GOOS=darwin go build ${current_dir}/main.go

linuxapp:
	@GOOS=linux go build ${current_dir}/main.go

.PHONY:windowsapp,macapp,linuxapp