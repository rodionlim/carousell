# Overview

Carousell is both a library and CLI tool written in Golang that fetches Carousell listings and notifies users.

# Installing

Using Carousell is easy. First, use `go install` to install the latest version of the library. This command will install the `carousell` executable along with the library and its dependencies:

```
go install -u github.com/rodionlim/carousell
```

# Usage

There are two commands, `get` and `notify`. Flags can be used to modify the search behaviour, e.g. `-r` flag will query for only recent listings.

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

# License

Carousell is released under the Apache 2.0 license. See [LICENSE](https://github.com/rodionlim/carousell/blob/master/LICENSE.txt)
