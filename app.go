package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	port := flag.Int("p", 9095, "server port, default: 9095")
	flag.Parse()

	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(w, "hello")
	}
	mirrorHandler := func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(10 * time.Millisecond)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			body = []byte(fmt.Sprintf("<read failed: %v>", err))
		}

		url := req.URL.Path
		if req.URL.Query().Encode() != "" {
			url += "?" + req.URL.Query().Encode()
		}

		content := fmt.Sprintf(`Your Request server from port %d
===============
Method: %s
URL   : %s
Header: %v
Body  : %s
`, *port, req.Method, url, req.Header, body)

		_, _ = io.WriteString(w, content)
	}

	http.HandleFunc("/", mirrorHandler)
	http.HandleFunc("/pipeline/activity/1", helloHandler)
	http.HandleFunc("/pipeline/activity/2", helloHandler)

	fmt.Printf("Server started at port: %d\n", *port)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
