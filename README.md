# Git P2 (Beta)

Git P2 (Git parallel push) is a command line utility written in Go that allows you to git push to multiple remotes at the same time and in parallel as opposed to needing to wait for each git push to happen sequentially as it does when you put many urls under one remote.

This is extremely useful if you manage many versions of the same app for different clients on Heroku and/or Dokku.

## Ignoring remotes
If you don't want gitp2 to push to all your remotes, you can create a file by the name of `.gitp2ignore` in your git root directory. Each line in that file should be the name of a remote to ignore when using Gitp2.

## Installation and Upgrading

`wget -O /usr/local/bin/gitp2 https://github.com/chiedolabs/gitp2/raw/master/gitp2 && chmod +x /usr/local/bin/gitp2
`

## Usage
Run `gitp2` in your root git directory.

## Development

- Git clone this repo into your Gopath.
- Make changes
- Run `go build`

#### Testing

- Run `go build`
- Try to push your changes using `./gitp2`

## Support

- Mac OS
- Linux (untested)
