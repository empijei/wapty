# Wapty (Web Application Penetration Testing Suite)y

This project is meant to become an open source replacement for Burp Suite

# DISCLAIMER
VERY EARLY STAGE, Don't use this. 

Really, don't

But if you really do want to use this

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
for currently available features and the roadmap please go to [ROADMAP.md](https://github.com/empijei/wapty/blob/master/ROADMAP.md)

# Contribute
Pull Requests are welcome, or if you'd like to actively contribute or become part of the team feel free to shoot me a short email of presentation at [empijei@gmail.com](mailto:empijei@gmail.com)
