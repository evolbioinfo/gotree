# GoTree
GoTree is a set of command line tools to manipulate phylogenetic trees. It is implemented in [Go](https://golang.org/) language.

## Installation
### Binaries
You can download already compiled binaries for the latest release in the [release](https://github.com/fredericlemoine/gotree/releases) section.
Binaries are available for MacOS, Linux, and Windows (32 and 64 bits).

Once downloaded, you can just run the executable without any other downloads.

### From sources
In order to compile gotree, you must first [download](https://golang.org/dl/) and [install](https://golang.org/doc/install) Go on your system.

Then you just have to type :
```
go get github.com/fredericlemoine/gotree/
```
This will download GoTree sources from github, and all its dependencies.

You can then build it with:
```
cd $GOPATH/src/github.com/fredericlemoine/gotree/
make
```
The `gotree` executable should be located in the `$GOPATH/bin` folder.

## Usage
