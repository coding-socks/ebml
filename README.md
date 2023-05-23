# EBML

An EBML parser written in Go.

- [Introduction](#introduction)
- [Production readiness](#production-readiness)
- [Documents](#documents)
- [Similar libraries](#similar-libraries)

## Introduction

> Extensible Binary Meta Language (EBML) is a generalized file format for any kind of data, aiming to be a binary equivalent to XML. It provides a basic framework for storing data in XML-like tags. It was originally developed for the Matroska audio/video container format.

Source: https://en.wikipedia.org/wiki/Extensible_Binary_Meta_Language

This library is based on the July 2020 version of [RFC 8794][rfc8794] (with additions from [github.com/ietf-wg-cellar/ebml-specification][ebml-specification]). This document did not reach ["Internet Standard"](https://tools.ietf.org/html/rfc2026#section-4.1.3) status yet. RFC 8794 is in a ["Proposed Standard"](https://tools.ietf.org/html/rfc2026#section-4.1.1) status.

The goal of this project is to create an implementation based on the document and during the implementation provide feedback.

## Production readiness

**This project is still in alpha phase.** In this stage the public API can change between days.

Beta version will be considered when the feature set covers the documents the implementation is based on, and the public API is reached a mature state.

Stable version will be considered only if enough positive feedback is gathered to lock the public API and all document the implementation is based on became ["Internet Standard"](https://tools.ietf.org/html/rfc2026#section-4.1.3).

## Documents

### Official sites

- [libEBML](http://matroska-org.github.io/libebml/)
- [EBML Specification](https://matroska-org.github.io/libebml/specs.html)
- [Matroska](https://www.matroska.org/index.html)
- [Matroska Element Specification](https://matroska.org/technical/elements.html)
- [WebM](https://www.webmproject.org/)
- [WebM Container Guidelines](https://www.webmproject.org/docs/container/)

Huge thanks to the [Matroska.org](https://www.matroska.org/) for their work.

### IETF Documents

- [RFC 8794: Extensible Binary Meta Language][rfc8794]

Huge thanks to the [IETF CELLAR Working Group](https://datatracker.ietf.org/wg/cellar/charter/) for their work.

## Inspiration

Inspiration for the implementation comes from the following places:

- https://pkg.go.dev/database/sql#Drivers
- https://pkg.go.dev/database/sql#Register
- https://pkg.go.dev/encoding/json#Decoder
- https://pkg.go.dev/golang.org/x/image/vp8#Decoder

## Similar libraries

Last updated: 2023-05-22

| URL                                                              | Status                      |
|------------------------------------------------------------------|-----------------------------|
| https://github.com/at-wat/ebml-go                                | In active development       |
| https://github.com/ebml-go/ebml + https://github.com/ebml-go/webm | Last updated on 17 Nov 2022 |
| https://github.com/ehmry/go-ebml                                 | Deleted                     |
| https://github.com/jacereda/ebml                                 | Last updated on 10 Jan 2016 |
| https://github.com/mediocregopher/ebmlstream                     | Last updated on 15 Dec 2014 |
| https://github.com/pankrator/ebml-parser                         | Last updated on 24 Jun 2020 |
| https://github.com/pixelbender/go-matroska                       | Last updated on 29 Oct 2018 |
| https://github.com/pubblic/ebml                                  | Last updated on 12 Dec 2018 |
| https://github.com/quadrifoglio/go-mkv                           | Last updated on 20 Jun 2018 |
| https://github.com/rrerolle/ebml-go                              | Last updated on 1 Dec 2012  |
| https://github.com/remko/go-mkvparse                             | Last updated on 19 May 2022 |
| https://github.com/tpjg/ebml-go                                  | Last updated on 1 Dec 2012  |

[rfc8794]: https://tools.ietf.org/html/rfc8794
[ebml-specification]: https://github.com/ietf-wg-cellar/ebml-specification
