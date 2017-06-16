# Wapty (Web Application Penetration Testing Suite)y

[![License](https://img.shields.io/github/license/empijei/wapty.svg?maxAge=360000)](https://raw.githubusercontent.com/empijei/wapty/master/LICENSE)

[![Build Status](https://travis-ci.org/empijei/wapty.svg?branch=master)](https://travis-ci.org/empijei/wapty) 
[![Go Report Card](https://goreportcard.com/badge/github.com/empijei/wapty)](https://goreportcard.com/report/github.com/empijei/wapty)
[![Changelog](https://img.shields.io/github/release/empijei/wapty.svg?maxAge=2592000)](https://github.com/empijei/wapty/releases)
[![Telegram Chat](https://img.shields.io/badge/discuss-telegram-179cde.svg?maxAge=3600)](https://telegram.me/waptysuite)
[![Funding](https://img.shields.io/badge/funding-not yet-lightgrey.svg?maxAge=2592000)]()


This project is meant to become an open source replacement for Burp Suite, but we are [not there yet](https://github.com/empijei/wapty/blob/master/ROADMAP.md).
# Usage
```sh
go get empijei/wapty
cd $GOPATH/src/github.com/empijei/wapty
make installdeps
make
make install
```

on your browser set localhost:8080 as proxy

**_BEWARE_**: ignore proxy for localhost

Then visit localhost:8081

Root certs are in `$HOME/.wapty`

# Status
for currently available features and what will be supported please refer to [our roadmap](https://github.com/empijei/wapty/blob/master/ROADMAP.md).

# Contribute
Pull Requests are welcome, or if you'd like to actively contribute or become part of the team feel free to shoot me a short email of presentation at [empijei@gmail.com](mailto:empijei@gmail.com)
