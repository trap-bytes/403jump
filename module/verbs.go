package module

import (
	"fmt"
	"net/http"

	"github.com/trap-bytes/403jump/utils"
)

var Verbs = [9]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "TRACE", "PATCH", "INVENTED"}

func HttpRequestWithVerbs(client *http.Client, url string) int {
	bypass := 0

	for _, verb := range Verbs {
		req, err := http.NewRequest(verb, url, nil)
		if err != nil {
			fmt.Printf("Error creating a %s request for %s: %v\n", verb, url, err)
			return 0
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error performing a %s request for %s: %v\n", verb, url, err)
			return 0
		}
		defer resp.Body.Close()
		if resp.StatusCode < http.StatusBadRequest {
			bypassMessage := fmt.Sprintf("Possible bypass found for %s with %s method. Status code: %d\n", url, verb, resp.StatusCode)
			coloredMessage := utils.Colorize(bypassMessage, "\033[32m") // Green color
			fmt.Print(coloredMessage)
			bypass++
		}
	}

	return bypass
}
