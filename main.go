package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Application struct {
	Output     string
	Target     string
	Method     string
	Json       bool
	DisableSSL bool
	Headers    map[string]string
}

type headers []string

var buildNumber string
var buildVersion string
var silent bool

func (i *headers) String() string {
	return ""
}

func (i *headers) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	a := &Application{
		Output:     "",
		Target:     "",
		Method:     "GET",
		Json:       false,
		DisableSSL: false,
	}

	flag.CommandLine.StringVar(&a.Target, "target", a.Target, "Target to analyze")
	flag.CommandLine.StringVar(&a.Output, "output", a.Output, "Output file")
	flag.CommandLine.StringVar(&a.Method, "method", a.Method, "Request method")
	flag.CommandLine.BoolVar(&a.Json, "json", a.Json, "Json output format")
	flag.CommandLine.BoolVar(&a.DisableSSL, "disable-ssl", a.DisableSSL, "Don't verify the site's SSL certificate")

	h := headers{
		"user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}
	flag.Var(&h, "header", "Set additional request headers")

	sv := flag.Bool("version", false, "Show version and exit")
	nc := flag.Bool("no-color", false, "Disable color output")
	s := flag.Bool("silent", false, "Don't display any output")
	flag.Parse()

	silent = *s
	if *nc {
		color.NoColor = true // disables colorized output
	}

	if *sv {
		fmt.Printf("version: %s\nbuild number: %s\n", color.CyanString(buildVersion), color.CyanString(buildNumber))
		os.Exit(0)
	}

	if a.Target == "" {
		if silent == false {
			color.Red("[error] no target specified.")
		}
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: a.DisableSSL},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(a.Method, a.Target, nil)
	handleError(err)

	for _, header := range h {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) == 2 {
			req.Header.Set(parts[0], parts[1])
		} else {
			handleError(errors.New("invalid header provided: " + header))
		}
	}

	resp, err := client.Do(req)
	handleError(err)

	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	wappalyzerClient, err := wappalyzer.New()
	handleError(err)

	fingerprints := wappalyzerClient.Fingerprint(resp.Header, data)

	result := ""
	if a.Json {
		d, err := json.Marshal(fingerprints)
		handleError(err)

		result = string(d)
	} else {
		lines := make([]string, 0)
		for name, values := range fingerprints {
			lines = append(lines, fmt.Sprintf("%s: %v", color.CyanString(name), values))
		}
		result = strings.Join(lines, "\n")
	}

	if a.Output != "" {
		err = ioutil.WriteFile(a.Output, []byte(result)[:], 0644)
		handleError(err)
	}

	if silent == false {
		fmt.Printf(result)
	}
}

func handleError(err error) {
	if err != nil {
		if silent == false {
			color.Red("[error] %s", err.Error())
		}
		os.Exit(1)
	}
}
