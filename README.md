# webTest.go

A small tool to test that ones IP, TLS and DNS settings for ones little API/Web server is up and running properly. 

Edit config.json with full path to your TLS private and public keys. Set the 'Local' flag to false if server is supposed to be exposed to the interwebs.

Build and run. Then point browser towards your DNS for this server using `https://<your_IP/DNS>/`. If all succeeds, you should get a return similar to this:

```
TLS Certs loaded - running over https
Inbound from     : <some IP>:<some Port>
Response from    : <some IP>:443
```

Should TLS certificates not be found (as in wrong path/name given in `tlsName.json`), this little test-webserver will run using port 80. Assuming port is forwarded, test by pointing browser to `http://<your_IP/DNS>/`. If port 80 is not forwarded, test on same machine as server is running on by pointing browser to localhost. You should get a return similar to this:

```
No TLS Certs loaded - running over http
Inbound from     : <some IP>:<some Port>
Response from    : <some IP>:80
```

This app must be run as SUDO or root as hogging ports for web-server listening purposes demand it if using port 80 or 443.

## Build

Only imports from standard library in this app. No need to get any 3rd party libraries or frameworks.

First make sure you have Go installed and configured correctly, then enter `go build webTest.go` whilst in your local directory of this repo. Finally run using `sudo ./webTest` and follow on-screen instructions.

**Contact:**

location   | name/handle
-----------|---------
github:    | rDybing
Linked In: | Roy Dybing
MeWe:      | Roy Dybing

---

## Releases

- Version format: [major release].[new feature(s)].[bugfix patch-version]
- Date format: yyyy-mm-dd

#### v.1.1.1: 2019-07-06
- Cleaned up a bit.

#### v.1.1.0: 2019-07-05
- Removed manual setting of private IP and port, webserver will figure it out.

#### v.1.0.2: 2019-07-05
- Fixed termination on not finding TLS files.
- Added some test coverage of http-server.

#### v.1.0.1: 2019-07-05
- Now actually working as advertised.
- Changed method of obtaining connecting client IP.

#### v.1.0.0: 2019-07-05
- Working as advertised. No magic or extra bells and whistles added.

---

## License: MIT

**Copyright © 2019 Roy Dybing** 

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

ʕ◔ϖ◔ʔ
