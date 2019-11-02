# lintservemux
An analyzer to detect uses of the default http.ServeMux in Go programs

[![CircleCI](https://circleci.com/gh/reillywatson/lintservemux.svg?style=svg)](https://circleci.com/gh/reillywatson/lintservemux)
[![codecov](https://codecov.io/gh/reillywatson/lintservemux/branch/master/graph/badge.svg)](https://codecov.io/gh/reillywatson/lintservemux)

Serving traffic with the default http.ServeMux has some risks! If anything you import calls http.Handle(), it will get served on the same port as your application.

In particular, importing net/http/pprof registers things on the default ServeMux.

This linter tries to help avoid that by flagging uses of http.Handle() and http.HandleFunc(), calling ListenAndServe() with a nil Handler argument, or creating an http.Server with a nil Handler argument.
