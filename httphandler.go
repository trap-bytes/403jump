package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/trap-bytes/403jump/module"
	"github.com/trap-bytes/403jump/utils"
)

func CreateHTTPClientWProxy(proxy string) (*http.Client, error) {
	parts := strings.Split(proxy, ":")
	proxyIP := parts[0]
	proxyPortStr := parts[1]
	proxyPort, err := strconv.Atoi(proxyPortStr)
	if err != nil {
		fmt.Printf("Error converting proxy port to integer: %v\n", err)
		proxyPort = 8080
	}

	client := &http.Client{}
	if proxyIP != "" && proxyPort != 0 {
		proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%d", proxyIP, proxyPort))
		if err != nil {
			return nil, fmt.Errorf("Error parsing proxy URL: %v", err)
		}
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client, nil
}

func ProcessSingleTarget(client *http.Client, url, cookie, header string) {
	fmt.Printf("Sending request to URL: %s\n", url)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bypassFound += module.HttpRequestWithVerbs(client, url, cookie, header)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bypassFound += module.HttpRequestWithHeaders(client, url, cookie, header)
	}()

	if utils.HasPath(url) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bypassFound += module.HttpRequestPathFuzzing(client, url, cookie, header)
		}()
	}

	wg.Wait()
}

func ProcessMultipleTargets(client *http.Client, file, cookie, header string) {
	fmt.Printf("Processing targets from file: %s\n", file)

	entries, err := utils.ReadTargetsFromFile(file)
	if err != nil {
		fmt.Println("Error reading targets:", err)
		return
	}

	var wg sync.WaitGroup

	for _, url := range entries {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			ProcessSingleTarget(client, url, cookie, header)
		}(url)
	}

	wg.Wait()
}
