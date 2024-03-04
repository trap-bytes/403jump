# 403JUMP

[![Go Report Card](https://goreportcard.com/badge/github.com/trap-bytes/403jump)](https://goreportcard.com/report/github.com/trap-bytes/403jump)


403JUMP is a tool designed for penetration testers and bug bounty hunters to audit the security of web applications. It aims to bypass HTTP 403 (Forbidden) pages using various techniques.


![403JUMP Image](https://raw.githubusercontent.com/trap-bytes/403jump/bd50f22b15d13670947ea732e1a14f7a33253106/static/tool.png)


## Features

- **Multiple Bypass Techniques Including:**  
	- Different HTTP Verbs
	- Different Headers
	- Path Fuzzing.
- **Customization:** Allows customization of headers and cookies for more targeted testing.
- **Concurrency:** Performs actions concurrently using goroutines for efficient and fast scanning.

## Install

```
go install github.com/trap-bytes/403jump@latest
```
## Usage:

```
403jump -h
```

This will display help for the tool. Here are all the arguments it supports.

```
Usage:
  403jump [arguments]

The arguments are:
  -t string    Specify the target URL (e.g., domain.com or https://domain.com)
  -f string    Specify the file (e.g., domain.txt)
  -p string    Specify the proxy URL (e.g., 127.0.0.1:8080)
  -c string    Specify cookies (e.g., user_token=g3p21ip21h; 
  -r string    Specify headers (e.g., Myheader: test
  -timeout     Specify connection timeout in seconds
  -h           Display help

Examples:
  403jump -t domain.com
  403jump -t https://domain.com -p 127.0.0.1:8080
  403jump -f domains.txt
  403jump -c "user_token=hjljkklpo"
  403jump -r "Myheader: test"
```
