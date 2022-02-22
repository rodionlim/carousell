# Overview

Carousell is both a library and CLI tool written in Golang that fetches Carousell listings and notifies users.

# Installation

Using Carousell is easy. First, use `go install` to install the latest version of the executable. This command will install the `carousell` executable in `$GOPATH/bin`:

```
go install github.com/rodionlim/carousell@latest
```

To use the library, `go get` will install the libraries and dependencies for your project.

```
go get github.com/rodionlim/carousell
```

Later, to receive updates, run

```
go get -u github.com/rodionlim/carousell
```

# Executable Usage

There are two commands, `get` and `notify`. Flags can be used to modify the search behaviour, e.g. `-r` flag will query for only recent listings, `-v` flag represents verbose mode and add logging to the output

`get` will fetch the listings and output them to the console.

```
carousell get "nintendo switch" -r
```

`notify` will periodically fetch the listings, and notify users on new listings in Slack. For slack to work, the environment variable `SLACK_ACCESS_TOKEN` has to be set and the appropriate permissions granted, e.g. inviting the application to the slack channel

```
carousell notify --slack-channel=CHANNEL_ID "nintendo switch" -r
```

To get help on the available flags, use the `-h` flag.

```
carousell -h
```

# Quick Start

The following Go code generates the latest 20 listing based on key words specified

```
r := carousell.NewReq(crs.WithSearch("nintendo switch"))
r.Get()
```

# License

Carousell is released under the Apache 2.0 license. See [LICENSE](https://github.com/rodionlim/carousell/blob/master/LICENSE.txt)
