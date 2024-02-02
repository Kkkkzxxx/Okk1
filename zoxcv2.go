package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"crypto/tls"
)

func sendRequest(client *http.Client, wg *sync.WaitGroup, reqURL string, userAgent string) {
	defer wg.Done()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "de,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		fmt.Println("Usage: go run main.go <reqURL> <rps> <threads>")
		return
	}

	reqURL := args[0]
	rps, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Error parsing requests per second:", err)
		return
	}

	threads, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Error parsing threads:", err)
		return
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
      "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; MAGWJS; rv:11.0) like Gecko",
      "Mozilla/5.0 (X11; Linux x86_64; rv:31.0) Gecko/20100101 Firefox/31.0",
      "Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
   	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 RuxitSynthetic/1.0 v3408563703863544352 t8056460500199558789 ath5ee645e0 altpriv cvcv=2 smf=0",		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 RuxitSynthetic/1.0 v3366505121 t6006063806750198674 athfa3c3975 altpub cvcv=2 smf=0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 RuxitSynthetic/1.0 v2828632065406795305 t8360729428027585528 ath259cea6f altpriv cvcv=2 smf=0",
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 30 * time.Second,
			TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			ForceAttemptHTTP2: true, // Mengaktifkan HTTP/2
		},
	}

	var wg sync.WaitGroup
	wg.Add(threads * rps)

	delay := time.Second / time.Duration(rps)

	for i := 0; i < threads; i++ {
		go func() {
			for j := 0; j < rps; j++ {
				userAgent := userAgents[rand.Intn(len(userAgents))]
				sendRequest(client, &wg, reqURL, userAgent)
				time.Sleep(delay)
			}
		}()
	}

	wg.Wait()
}
