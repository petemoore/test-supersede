# test-supersede

This is a utility that generates 5 docker-worker tasks and publishes a supersedes json file containing the 5 tasks, for testing superseding functionality.

You'll need go (golang) installed and configured.

# How to use

1. Fork my repo!
2. `go get github.com/<your github user>/test-supersede`
3. `cd "${GOPATH}/src/github.com/<your github user>/test-supersede"`
3. `./run.sh`

Depending on what you want to test, you can edit the task definition(s) in `main.go`.

The checked-in version creates a task that is designed to fail, which supersedes the other 4 tasks which are designed to pass.
