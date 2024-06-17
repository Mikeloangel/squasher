package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080"
	data := url.Values{}
	fmt.Println("Введите url")

	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		return
	}

	long = strings.TrimSpace(long)
	data.Set("url", long)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
		return
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Status code: ", response.StatusCode)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(body))
}
