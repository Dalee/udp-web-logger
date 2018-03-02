# UDP Web Logger

[![Build Status](https://travis-ci.org/Dalee/udp-web-logger.svg?branch=master)](https://travis-ci.org/Dalee/udp-web-logger)
[![codecov](https://codecov.io/gh/Dalee/udp-web-logger/branch/master/graph/badge.svg)](https://codecov.io/gh/Dalee/udp-web-logger)

## Requirements

- Go 1.10

## Why

The very common use case is when your applications use `syslog` as a transport for the logs.
In staging / production you can set up `ELK` stack with `logstash` parsing pattern like:

```
<%{POSINT:syslog_pri}>%{SYSLOGTIMESTAMP:syslog_timestamp} %{SYSLOGHOST:syslog_hostname} %{DATA:syslog_program}(?:\[%{POSINT:syslog_pid}\])?: %{GREEDYDATA:syslog_message}
```

But during developing the whole `ELK` stack is not needed mostly. That's where this logger comes in. Cool thing is that
you can embed this daemon into a `Vagrant` box.

## Usage

```
$ UDP_LISTEN="127.0.0.1:10110" WEB_LISTEN="127.0.0.1:10100" ./udp-web-logger --help

Usage: udp-web-logger [options]

Options:
--udp-read-buffer-size   size of buffer to read incoming UDP packet into. Default: 4096.
--max-messages           maximum amount of messages to keep. Default: 50.
--help                   prints this message.

Env:
UDP_LISTEN - address to listen UDP on. Default: 127.0.0.1:9010.
WEB_LISTEN - address to listen HTTP on. Default: 127.0.0.1:9000.
```

## Links

https://github.com/Dalee/elk-playground

https://github.com/Dalee/node-logger

https://github.com/Dalee/monolog-syslog3164
