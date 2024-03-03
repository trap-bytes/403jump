package utils

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func ValidateFlags(urlFlag, fileFlag string) error {
	if urlFlag != "" && fileFlag != "" {
		return fmt.Errorf("You can supply either a single target or a file, but not both\n Example usage: ./403jump -t domain.com")
	}

	if urlFlag == "" && fileFlag == "" {
		return fmt.Errorf("Please provide a target.\n Example usage: ./403jump -t domain.com\n Use -h for help.")
	}

	return nil
}

func ValidateUrl(inputURL string) (string, error) {
	u, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %v", err)
	}

	if u.Scheme == "" {
		inputURL = "https://" + inputURL
		u, _ = url.Parse(inputURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("Invalid URL scheme")
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	_, err = net.LookupHost(host)
	if err != nil {
		return "", err
	}

	if port != "" {
		inputURL = fmt.Sprintf("%s://%s:%s%s", u.Scheme, host, port, u.RequestURI())
	} else {
		inputURL = fmt.Sprintf("%s://%s%s", u.Scheme, host, u.RequestURI())
	}

	return inputURL, nil
}

func ReadTargetsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// Skip empty lines
			continue
		}
		validUrl, err := ValidateUrl(line)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
		} else {
			entries = append(entries, validUrl)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func Colorize(text string, colorCode string) string {
	resetColor := "\033[0m"
	return colorCode + text + resetColor
}

func HasPath(inputURL string) bool {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		fmt.Println("Skipping domain fuzzing due to an error parsing the URL path")
		return false
	}
	return parsedURL.Path != ""
}

func IsValidProxy(input string) bool {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return false
	}

	ip := parts[0]
	portStr := parts[1]

	if net.ParseIP(ip) == nil {
		return false
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return false
	}
	return true
}

func PrintBanner() {
	azureBold := "\033[1;36m"

	fmt.Println(Colorize("  _  _    ___ ____      _                       ", azureBold))
	fmt.Println(Colorize(" | || |  / _ \\___ \\    (_)                      ", azureBold))
	fmt.Println(Colorize(" | || |_| | | |__) |    _ _   _ _ __ ___  _ __  ", azureBold))
	fmt.Println(Colorize(" |__   _| | | |__ <    | | | | | '_ ` _ \\| '_ \\ ", azureBold))
	fmt.Println(Colorize("    | | | |_| |__) |   | | |_| | | | | | | |_) |", azureBold))
	fmt.Println(Colorize("    |_|  \\___/____/    | |\\__,_|_| |_| |_| .__/ ", azureBold))
	fmt.Println(Colorize("                      _/ |               | |    ", azureBold))
	fmt.Println(Colorize("                     |__/                |_|    ", azureBold))
	fmt.Println("")
	fmt.Print(Colorize("403JUMP - HTTP 403 Bypass Tool\n\n", azureBold))
}
