# Git P2

Git P2 (Git parallel push) is a simple Golang script that allows you to git push to multiple remotes at the same time and in parallel as opposed to needing to wait for each git push to happen sequentially as it does when you put many remotes under one origin.

## Ignoring remotes
If you don't want gitp2 to push to all your remotes, you can create a file by the name of `.gitp2ignore` in your git root directory. Each line in that file should be the name of a remote to ignore when using Gitp2.

## Installation
Coming soon.

## Usage
Run `gitp2` in your root git directory.

## Development

- Git clone this repo into your Gopath.
- Make changes
- Run `go build`

## Support

- Mac OS
- Linux (untested)