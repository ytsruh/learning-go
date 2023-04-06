# Running NodeJS inside Go
Using the os/exec package you can run Node files from Go. This however assumes that the system has access to Go. In a local machine this is not an issues (assuming Node is installed) but on a remote server the Go binary probably doesn't have access to NodeJS.
To overcome this, a Dockerfile is included that allows the Go binary to be built and for NodeJS & NPM to be installed so it can be used & accessed as it can on a local machine.
