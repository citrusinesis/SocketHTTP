package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	html     string = "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n"
	favicon  string = "HTTP/1.1 200 OK\r\nContent-Type: image/x-icon\r\n\r\n"
	notFound string = "HTTP/1.1 404 Not Found\r\n"
)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	request := make([]byte, 1024)
	_, err := conn.Read(request)
	if err != nil {
		return
	}

	headers := strings.Split(string(request), "\r\n")
	resource := strings.Split(headers[0], " ")[1]

	fmt.Println(
		time.Now().Format(time.Stamp),
		strings.Split(headers[0], " "),
	)

	conn.Write(getResponse(resource))
}

func getResponse(resource string) (response []byte) {
	var header string
	files := getFileList("./resource")

	if !contains(files, resource[1:]) {
		response = []byte(notFound)
		return
	}

	switch resource {
	case "/":
		resource = "/index.html"
		header = html
	case "/favicon.ico":
		header = favicon
	default:
		header = html
	}

	file, err := os.ReadFile("./resource" + resource)
	if err != nil {
		panic(err)
	}

	response = []byte(header)
	response = append(response, file...)
	return
}

func getFileList(dirPath string) (list []string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	list = append(list, "")

	for _, file := range files {
		if strings.Contains(file.Name(), ".html") ||
			strings.Contains(file.Name(), ".ico") {
			list = append(list, file.Name())
		}
	}

	return
}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
