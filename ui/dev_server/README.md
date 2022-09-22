# Dev Server

Since developing web UI is a visual, iterative process, it's important to have quick feedback when making changes. These shared
UI fragments are imported by other projects via Github, so having to commit,push,pull,import every time a change is made is not
feasible.

Dev Server is a simple HTTP server with a page for each UI fragment that can be used during development for quick feedback.

## Instructions

Run `./copy.sh && go build && ./dev_server`

The `copy` script copies assets from the UI packages into the Dev Server's `assets/local` folder. Add any new files
that are needed for the UI fragments to this script, as well as the wrapper HTML in `server.go`