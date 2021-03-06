[![Build Status](https://travis-ci.org/msgehard/go-exercism.png?branch=master)](https://travis-ci.org/msgehard/go-exercism)

Goals
===========

Provide developers an easy way to work with [exercism.io](http://exercism.io) that doesn't require a 
Ruby environment.

Development
===========
1. Install Go ```brew install go --cross-compile-common``` or the command appropriate for your platform. If that throws an 
error, try ```brew install go --crosscompile-commone --with-llvm```.
1. Fork and clone.
1. Run ```git submodule update --init --recursive```
1. Write a test.
1. Run ``` bin/test ``` and watch test fail.
1. Make test pass.
1. Submit a pull request.

Building
========
1. Run ```bin/build``` and the binary for your platform will be built into the out directory.
1. Run ```bin/build-all``` and the binaries for OSX, Linux and Windows will be built into the release directory.
