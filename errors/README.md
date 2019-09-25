# errors

Recently [Golang 1.13](https://golang.org/doc/go1.13) comes out with some [ideas](https://golang.org/doc/go1.13#error_wrapping) from the [error handling proposal](https://go.googlesource.com/proposal/+/master/design/29934-error-values.md).

I just write something based on them to see what I can do.

## cerrors - common errors

This package contains some sample errors that commonly used to show how to use `detail.Detail`.

## detail

The `detail.Detail` is an error with a message, an aliasing error, an inner error, and a stack trace.

It could be used as a template to create other errors. See [cerrors](./cerrors).

## wrap

This package contains some functions to create a simple error quickly.
