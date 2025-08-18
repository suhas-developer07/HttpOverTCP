package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/suhasdeveloper07/httpovertcp/internal/request"
	"github.com/suhasdeveloper07/httpovertcp/internal/response"
	"github.com/suhasdeveloper07/httpovertcp/internal/server"
)

const port = 42069

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {

		h := response.GetDefaultHeaders(0)
		body := respond200()
		status := response.StatusOk
		// Handle the request and write the response
		if req.RequestLine.RequestTarget == "/yourproblem" {
			w.WriteStatusLine(response.StatusBadRequest)
			body = respond400()
			status = response.StatusBadRequest

		} else if req.RequestLine.RequestTarget == "/myporblem" {
			body = respond500()
			status = response.StatusInternalServerError
		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/video"){
			f,_ := os.ReadFile("assets/vim.mp4")
			h.Replace("content-type","vidio/mp4")
			h.Replace("content-length",fmt.Sprintf("%d",len(f)))

			w.WriteStatusLine(response.StatusOk)
			w.WriteHeaders(*h)
			w.WriteBody(f)

		}else if strings.HasPrefix(req.RequestLine.RequestTarget,"/httpbin/stream"){
			target := req.RequestLine.RequestTarget
			res,err := http.Get("https://httpbin.org" + target[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				status = response.StatusInternalServerError
			}else {
				w.WriteStatusLine(response.StatusOk)
				h.Delete("Content-Lenght")
				h.Set("transfer-encoding","chunked")
				h.Replace("content-type","text/plain")
				w.WriteHeaders(*h)

				for {
					data := make([]byte, 32)
					n,err := res.Body.Read(data)
					if err != nil {
						break
					}

					w.WriteBody([]byte(fmt.Sprintf("%x\r\n",n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n\r\n"))
				return
			}
		}
		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-type","text/html")
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
