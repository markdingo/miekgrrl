# miekgrrl - Derive a [rrl](https://github.com/markdingo/rrl) ResponseTuple from a [miekg](https://github.com/miekg/dns) Msg

## Introduction

`miekgrrl` is a tiny specialized adaptor package with a single `Derive()` function which
converts a `dns.Msg` to a `rrl.ResponseTuple` suitable for passing to
the [rrl](https://github.com/markdingo/rrl) `Debit()` function.

By way of background, the [rrl](https://github.com/markdingo/rrl) package is an agnostic
implementation of the [Response Rate Limiting](https://kb.isc.org/docs/aa-01148) algorithm
as originally developed by [ISC](https://www.isc.org). You will need to be familiar with
both the [rrl](https://github.com/markdingo/rrl) and the
[miekg](https://github.com/miekg/dns) package prior to understanding the purpose of this
package.

The sole reason for adaptor packages such as this is to keep the `rrl` package agnostic.

### Project Status

[![Build Status](https://github.com/markdingo/miekgrrl/actions/workflows/go.yml/badge.svg)](https://github.com/markdingo/miekgrrl/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/markdingo/miekgrrl/branch/main/graph/badge.svg?token=211OVOI2AV)](https://codecov.io/gh/markdingo/rrl)
[![CodeQL](https://github.com/markdingo/miekgrrl/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/markdingo/miekgrrl/actions/workflows/codeql-analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/markdingo/miekgrrl)](https://goreportcard.com/report/github.com/markdingo/miekgrrl)
[![Go Reference](https://pkg.go.dev/badge/github.com/markdingo/miekgrrl.svg)](https://pkg.go.dev/github.com/markdingo/miekgrrl)

## Sample Code

    package main

    import (
        "net"

        "github.com/markdingo/miekgrrl"
        "github.com/markdingo/rrl"
        "github.com/miekg/dns"
    )

    func main() {

      R := rrl.NewRRL(rrl.NewConfig())
      response := &dns.Msg{}
      ... populate response ...

      tuple := miekgrrl.Derive(response, "")   // Derive the RRL.ResponseTuple from dns.Msg

      action, _, _ := R.Debit(net.Addr, tuple) // Apply RRL rules to ResponseTuple
      switch action {
          case rrl.Drop:
          ...
          case rrl.Send:
          ..
          case rrl.Slip:
          ..


Alternatively, the `Derive()` function is designed such that it can be used directly as an
argument to the `rrl.Debit()` function. E.g:

      ...
      R := rrl.NewRRL(rrl.NewConfig())
      response := &dns.Msg{}
      ... populate response ...

      action, _, _ := R.Debit(net.Addr, miekgrrl.Derive(response, ""))
      switch action {
      ...

## Installation

`miekgrrl` requires [go](https://golang.org) version 1.19 or later.

Once your application imports `"github.com/markdingo/miekgrrl"`, then `"go build"` or `"go
mod tidy"` should download and compile the `miekgrrl` package automatically.

## Community

If you have any problems using `miekgrrl` or suggestions on how it can do a better job,
don't hesitate to create an [issue](https://github.com/markdingo/miekgrrl/issues) on the
project home page.
This package can only improve with your feedback.

## Copyright and License

`miekgrrl` is Copyright :copyright: 2023 Mark Delany and is licensed under the BSD
2-Clause "Simplified" License.
