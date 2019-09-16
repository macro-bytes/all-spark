all: allspark_cli allspark_daemon
allspark_daemon:
	go build -o allspark_daemon allspark_orchestrator

allspark_cli:
	go build -o allspark_cli --tags cli allspark_orchestrator

clean:
	rm -f allspark_cli allspark_daemon
	
