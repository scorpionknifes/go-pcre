# go-pcre

[![GoDoc](https://pkg.go.dev/github.com/scorpionknifes/go-pcre?status.svg)](https://pkg.go.dev/github.com/scorpionknifes/go-pcre)

This is a Go language package providing support for
Perl Compatible Regular Expressions (PCRE).

Current version is pcre 8.45

## Features:

- Support Windows, Linux & MacOS

## Installation

Install the package with the following:

    go get github.com/scorpionknifes/go-pcre

## Usage

Go programs that depend on this package should import
this package as follows to allow automatic downloading:

    import "github.com/scorpionknifes/go-pcre"

## History

This is a clone of [go-pcre](https://github.com/rubrikinc/go-pcre)

This was a clone of
[golang-pkg-pcre](http://git.enyo.de/fw/debian/golang-pkg-pcre.git)
by Florian Weimer, which has been placed on Github by Glenn Brown,
so it can be fetched automatically by Go's package installer.

Glenn Brown added `FindIndex()` and `ReplaceAll()`
to mimic functions in Go's default regexp package.

Mathieu Payeur Levallois added `Matcher.ExtractString()`.

Malte Nuhn added `GroupIndices()` to retrieve positions of a matching group.

Chandra Sekar S added `Index()` and stopped invoking `Match()` twice in `FindIndex()`.

Misakwa added support for `pkg-config` to locate `libpcre`.

Yann Ramin added `ReplaceAllString()` and changed `Compile()` return type to `error`.

Nikolay Sivko modified `name2index()` to return error instead of panic.

Harry Waye exposed raw `pcre_exec`.

Hazzadous added partial match support.

Pavel Gryaznov added support for JIT compilation.
