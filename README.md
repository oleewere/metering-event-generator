# meteringp

[![GoDoc Widget](https://godoc.org/github.com/oleewere/meteringp/producer?status.svg)](https://godoc.org/github.com/oleewere/meteringp/producer)
[![Build Status](https://travis-ci.org/oleewere/meteringp.svg?branch=master)](https://travis-ci.org/oleewere/meteringp)
[![Go Report Card](https://goreportcard.com/badge/github.com/oleewere/meteringp)](https://goreportcard.com/report/github.com/oleewere/meteringp)
![license](http://img.shields.io/badge/license-Apache%20v2-blue.svg)

Tool for generating metering JSON events to stdout or files. These outputs can be processed by other tools like fluentd in order to parse and send to an external service.

## Installation 

### Installation on Mac OSX
```bash
brew tap oleewere/repo
brew install meteringp
```

### Installation on Linux

Using wget:
```bash
METERINGP_VERSION=0.2.0
wget -qO- "https://github.com/oleewere/meteringp/releases/download/v${METERINGP_VERSION}/meteringp_${METERINGP_VERSION}_linux_64-bit.tar.gz" | tar -C /usr/bin -zxv meteringp
```

Using curl:
```bash
METERINGP_VERSION=0.2.0
curl -L -s "https://github.com/oleewere/meteringp/releases/download/v${METERINGP_VERSION}/meteringp_${METERINGP_VERSION}_linux_64-bit.tar.gz" | tar -C /usr/bin -xzv meteringp
```

## Usage

## Examples

```bash
meteringp --config sample/meteringp.conf
```

Using with docker:

```bash
docker build -t oleewere/meteringp .
docker run --rm oleewere/meteringp --config /sample/meteringp.conf
```

See sample folder to check the configuration options.
