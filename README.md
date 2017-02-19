# Don't use this.
## Really, don't
### But if you really do want to use this
no dependencies, no need to install, just
```
go run wapty.go
```
on your browser set localhost:8080 as proxy

**_BEWARE_**: ignore proxy for localhost

Then visit localhost:8081

Root certs are in `$HOME/.wapty`

### Currently available fetaures:

#### Proxy > Intercept
ATM it is not possible to change the host on the fronted, Working on it.
![Intercept Tab](/pics/intercept.png "Intercept")

#### Proxy > HTTP History
Clicking on a row shows only the original request/response on the frontend. Working on it.
![Hist Tab](/pics/history.png "History")

(Yes, I know the page is slightly bigger than the viewport and the layout is not resizable working on that too)

### Want to contribute?
If you do please [email me](mailto:empijei@gmail.com), I can provide some insights and planning for what's coming.

=)
