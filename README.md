# Wappalyzer CLI

Cli app based upon [projectdiscovery/wappalyzergo](https://github.com/projectdiscovery/wappalyzergo).


## Installation
```bash
go install -v github.com/webklex/wappalyzer@main
```

## Usage
```bash
wappalyzer --target https://www.google.com/ --disable-ssl --output output.txt --json
```
Example output:
```json
{
  "Google Web Server":{},
  "HSTS":{},
  "HTTP/3":{}
}
```

### Available arguments
```bash
Usage of wappalyzer:
  -target string  Target to analyze
  -output string  Output file
  -method string  Request method (default "GET")
  -header value   Set additional request headers
  -disable-ssl    Don't verify the site's SSL certificate
  -json           Json output format
  -no-color       Disable color output
  -silent         Don't display any output
  -version        Show version and exit
```


## Build
```bash
git clone https://github.com/webklex/wappalyzer
cd wappalyzer
go build -a -ldflags "-w -s -X main.buildNumber=1 -X main.buildVersion=custom" -o wappalyzer
```


## Security
If you discover any security related issues, please email github@webklex.com instead of using the issue tracker.


## Credits
- [Webklex][link-author]
- [projectdiscovery/wappalyzergo](https://github.com/projectdiscovery/wappalyzergo)
- [All Contributors][link-contributors]


## License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.


[link-author]: https://github.com/webklex
[link-contributors]: https://github.com/webklex/wappalyzer/graphs/contributors