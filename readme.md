# tagpath

tagpath is a tool to generate query selectors for an element in an HTML document.

## Installation

	go get -u -v github.com/insomnimus/tagpath/...

If above doesn't work, try this:

	git clone https://www.github.com/insomnimus/tagpath
	cd tagpath
	go mod tidy
	go install

## Library Use

You can also use tagpath/engine in your project, for more information see [the documentation](https://pkg.go.dev/github.com/insomnimus/tagpath/engine).

## Tool Use

	tagpath [options] file|url '<tag>'

Options are:

-	`-a, --all: print all the matches instead of just the first`
-	`-h, --help: show help`