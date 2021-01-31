# EBML

An EBML parser written in Go.

- [Introduction](#introduction)
- [Documents](#documents)
- [Similar libraries](#similar-libraries)

## Introduction

> Extensible Binary Meta Language (EBML) is a generalized file format for any kind of data, aiming to be a binary equivalent to XML. It provides a basic framework for storing data in XML-like tags. It was originally developed for the Matroska audio/video container format.

Source: https://en.wikipedia.org/wiki/Extensible_Binary_Meta_Language

This implementation is based on the July 2020 version of [RFC 8794][rfc8794]. This RFC which is in a ["PROPOSED STANDARD"](https://www.rfc-editor.org/rfc/rfc2026.html#section-4.1.1) status.

```
[...]
 
    Implementors should treat Proposed Standards as immature  
    specifications.  It is desirable to implement them in order to gain  
    experience and to validate, test, and clarify the specification.  
    However, since the content of Proposed Standards may be changed if  
    problems are found or better solutions are identified, deploying  
    implementations of such standards into a disruption-sensitive  
    environment is not recommended.
 
[...]
```

Source: https://www.rfc-editor.org/rfc/rfc2026.html#section-4.1.1

## Documents

### Official sites

- [libEBML](http://matroska-org.github.io/libebml/)
- [EBML Specification](https://matroska-org.github.io/libebml/specs.html)
- [Matroska](https://www.matroska.org/index.html)
- [Matroska Element Specification](https://matroska.org/technical/elements.html)
- [WebM](https://www.webmproject.org/)
- [WebM Container Guidelines](https://www.webmproject.org/docs/container/)

### IETF Documents

- [RFC 8794: Extensible Binary Meta Language][rfc8794]
- [draft-ietf-cellar-matroska-06: Matroska Media Container Format Specifications](https://www.ietf.org/archive/id/draft-ietf-cellar-matroska-06.html)

Huge thanks to the [IETF CELLAR Working Group](https://datatracker.ietf.org/wg/cellar/charter/) for their work.

## Similar libraries

- https://github.com/at-wat/ebml-go
- https://github.com/pankrator/ebml-parser
- https://github.com/remko/go-mkvparse
- https://github.com/pixelbender/go-matroska
- https://github.com/quadrifoglio/go-mkv
- https://github.com/ebml-go/webm

[rfc8794]: https://tools.ietf.org/html/rfc8794
