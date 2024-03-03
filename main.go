package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/trap-bytes/403jump/utils"
)

var bypassFound int64

func main() {

	utils.PrintBanner()

	var urlFlag string
	var fileFlag string
	var proxy string
	var cookie string
	var header string
	var timeout int
	var client *http.Client

	flag.StringVar(&urlFlag, "t", "", "Specify the target URL (e.g., domain.com or https://domain.com)")
	flag.StringVar(&fileFlag, "f", "", "Specify the file containing target URLs (e.g., domains.txt)")
	flag.StringVar(&proxy, "p", "", "Specify the proxy URL (e.g., 127.0.0.1:8080)")
	flag.StringVar(&cookie, "c", "", "Specify cookies (e.g., user_token=g3p21ip21h; )")
	flag.StringVar(&header, "r", "", "Specify headers (e.g., Myheader: test )")
	flag.IntVar(&timeout, "timeout", 10, "Specify connection timeout in seconds")

	helpFlag := flag.Bool("h", false, "Display help")

	flag.Parse()

	if *helpFlag {
		fmt.Println("403jump is a tool for trying to bypass page protected by HTTP 403 status code.")
		fmt.Println("\nUsage:")
		fmt.Printf("  %s [arguments]\n", os.Args[0])
		fmt.Println("\nThe arguments are:")
		fmt.Println("  -t string    Specify the target URL (e.g., domain.com or https://domain.com)")
		fmt.Println("  -f string    Specify the file (e.g., domain.txt)")
		fmt.Println("  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)")
		fmt.Println("  -c string    Specify cookies (e.g., user_token=g3p21ip21h; ")
		fmt.Println("  -r string    Specify headers (e.g., Myheader: test")
		fmt.Println("  -timeout		Specify connection timeout")
		fmt.Println("  -h           Display help")

		fmt.Println("\nExamples:")
		fmt.Printf("  %s -t domain.com\n", os.Args[0])
		fmt.Printf("  %s -t https://domain.com -p 127.0.0.1:8080\n", os.Args[0])
		fmt.Printf("  %s -f domains.txt\n", os.Args[0])
		fmt.Printf("  %s -c \"user_token=hjljkklpo\"\n", os.Args[0])
		fmt.Printf("  %s -r \"Myheader: test\"\n", os.Args[0])

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
			client, err = CreateHTTPClientWProxy(proxy, timeout)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("Invalid proxy:", proxy)
			fmt.Println("Please insert a valid proxy in the ip:port format")
			return
		}
	} else {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}

	if urlFlag != "" {
		urlFlag, err := utils.ValidateUrl(urlFlag)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			return
		} else {
			ProcessSingleTarget(client, urlFlag, cookie, header)
		}
	} else {
		ProcessMultipleTargets(client, fileFlag, cookie, header)
	}

	if bypassFound > 0 {
		bypassMessage := fmt.Sprintf("\nFound %d possible bypass(es)\n", bypassFound)
		coloredMessage := utils.Colorize(bypassMessage, "\033[1;32m") // Green color
		fmt.Print(coloredMessage)
	} else {
		noBypassMessage := "\nNo bypass found\n"
		coloredMessage := utils.Colorize(noBypassMessage, "\033[1;31m") // Red color
		fmt.Print(coloredMessage)
	}

}
