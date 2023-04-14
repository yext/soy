# soy

This is a fork of https://github.com/robfig/soy.

[![GoDoc](http://godoc.org/github.com/yext/soy?status.png)](http://godoc.org/github.com/yext/soy)
[![Build Status](https://travis-ci.org/yext/soy.png?branch=master)](https://travis-ci.org/yext/soy)
[![Go Report Card](https://goreportcard.com/badge/yext/soy)](https://goreportcard.com/report/yext/soy)

Go implementation for Soy templates aka [Google Closure
Templates](https://github.com/google/closure-templates).  See
[godoc](http://godoc.org/github.com/yext/soy) for more details and usage
examples.

This project requires Go 1.12 or higher due to one of the transitive
dependencies requires it as a minimum version; otherwise, Go 1.11 would
suffice for `go mod` support.

Be sure to set the env var `GO111MODULE=on` to use the `go mod` dependency
versioning when building and testing this project.
