// main.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/trap-bytes/403jump/utils"
)

var bypassFound int

func main() {
	var urlFlag string
	var fileFlag string
	var proxy string
	var client *http.Client

	flag.StringVar(&urlFlag, "t", "", "Specify the target URL (e.g., domain.com or https://domain.com)")
	flag.StringVar(&fileFlag, "f", "", "Specify the file containing target URLs (e.g., domains.txt)")
	flag.StringVar(&proxy, "p", "", "Specify the proxy URL (e.g., 127.0.0.1:8080)")

	helpFlag := flag.Bool("h", false, "Display help")

	flag.Parse()

	if *helpFlag {
		fmt.Println("403jump is a tool for trying to bypass page protected by HTTP 403 status code.")
		fmt.Println("\nUsage:")
		fmt.Printf("  %s <command> [arguments]\n", os.Args[0])
		fmt.Println("\nThe arguments are:")
		fmt.Println("  -f string    Specify the file (e.g., domain.txt)")
		fmt.Println("  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)")
		fmt.Println("  -t string    Specify the target URL (e.g., domain.com or https://domain.com)")
		fmt.Println("  -h           Display help")

		fmt.Println("\nExamples:")
		fmt.Printf("  %s -t domain.com\n", os.Args[0])
		fmt.Printf("  %s -t https://domain.com -p 127.0.0.1:8080\n", os.Args[0])
		fmt.Printf("  %s -f domains.txt\n", os.Args[0])
		return
	}

	err := utils.ValidateFlags(urlFlag, fileFlag)
	if err != nil {
		fmt.Println(err)
		return
	}
	bypassFound = 0

	if len(proxy) > 0 {
		if utils.IsValidProxy(proxy) {
			fmt.Println("Using proxy:", proxy)
			client, err = CreateHTTPClientWProxy(proxy)
			if err != nil {
				fmt.Printf("Error in creating HTTP client with proxy: %v\n", err)
				return
			}
		} else {
			fmt.Println("Invalid proxy:", proxy)
			fmt.Println("Please insert a valid proxy in the ip:port format")
			return
		}
	} else {
		client = &http.Client{}
	}

	if urlFlag != "" {
		urlFlag, err := utils.ValidateUrl(urlFlag)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		} else {
			ProcessSingleTarget(client, urlFlag)
		}
	} else {
		ProcessMultipleTargets(client, fileFlag)
	}

	if bypassFound > 0 {
		bypassMessage := fmt.Sprintf("Found %d possible bypass(es)\n", bypassFound)
		coloredMessage := utils.Colorize(bypassMessage, "\033[32m") // Green color
		fmt.Print(coloredMessage)
	} else {
		noBypassMessage := "No bypass found\n"
		coloredMessage := utils.Colorize(noBypassMessage, "\033[31m") // Red color
		fmt.Print(coloredMessage)
	}

}
