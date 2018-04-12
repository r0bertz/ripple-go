# Ripple

Ripple related go libraries.

[![Build Status](https://travis-ci.org/r0bertz/ripple-go.svg?branch=master)](https://travis-ci.org/r0bertz/ripple-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/r0bertz/ripple-go?style=flat-square)](https://goreportcard.com/report/github.com/r0bertz/ripple-go)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/r0bertz/ripple-go)
[![Release](https://img.shields.io/github/release/r0bertz/ripple-go.svg?style=flat-square)](https://github.com/r0bertz/ripple-go/releases/latest)

## Getting started

### Prerequisites

Install [Bazel](https://docs.bazel.build/versions/master/install-ubuntu.html)

### Generate BUILD files

```
bazel run //:gazelle
```

### Build

```
bazel build //cmd/tx
```

### Test

```
bazel test ...
```
