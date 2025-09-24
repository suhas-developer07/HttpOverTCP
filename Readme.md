
# HTTP Over TCP

## 📖 Overview

**HTTP Over TCP** is a minimal **HTTP/1.1 implementation written from scratch in Go**, built directly on top of raw **TCP sockets**.

Instead of using Go’s built-in `net/http` package, this project manually handles:

* **TCP connections**
* **Request parsing** (request line, headers, body)
* **Response writing** (status line, headers, body, chunked encoding)
* **Basic server abstraction** similar to a mini framework

This project is designed for **learning and experimentation** with networking and HTTP internals.

---

## 🚀 Features

* Parse **HTTP/1.1 request lines, headers, and body**
* Validate and manage headers
* Response writer with:

  * Status codes (`200`, `400`, `500`, …)
  * Static file serving
  * Streaming (chunked transfer encoding)
* Simple **server abstraction** (`server.Serve`)
* Two runnable demos:

  * `tcplistener`: raw TCP listener that prints parsed requests
  * `httpserver`: a full HTTP server with routes

---

## 📂 Project Structure

```
.
├── cmd/
│   ├── httpserver/         # Example HTTP server with routes & responses
│   │   └── main.go
│   └── tcplistener/        # Raw TCP listener that parses HTTP requests
│       └── main.go
│       
├── internal/
│   ├── headers/            # Header parsing & management
│   │   ├── headers.go
│   │   └── headers_test.go
│   ├── request/            # HTTP request parsing state machine
│   │   ├── request.go
│   │   └── request_test.go
│   ├── response/           # HTTP response writer
│   │   └── response.go
│   └── server/             # Server abstraction (Serve, Close, etc.)
│       └── server.go
├── message.txt             # Sample data for testing
└── Readme.md
```

---

## 🛠️ Getting Started

### Clone the repo

```bash
git clone https://github.com/suhasdeveloper07/httpovertcp.git
cd httpovertcp
```

### Run the raw TCP listener

This simply prints the incoming HTTP request structure:

```bash
go run ./cmd/tcplistener
```

Then in another terminal:

```bash
curl -v http://localhost:42069/
```

### Run the HTTP server

This version demonstrates actual responses:

```bash
go run ./cmd/httpserver
```

---

## 📡 Example Routes (`httpserver`)

* `/` → **200 OK** HTML page
* `/yourproblem` → **400 Bad Request**
* `/myporblem` → **500 Internal Server Error**
* `/video` → serves `assets/vim.mp4`
* `/httpbin/stream/<n>` → streams data from [httpbin.org](https://httpbin.org/stream)

---

## 🧑‍💻 Example: Writing Your Own Handler

You can build custom handlers using `server.Serve`:

```go
server, _ := server.Serve(42069, func(w *response.Writer, req *request.Request) {
    body := []byte("Hello from custom handler!")
    h := response.GetDefaultHeaders(len(body))
    h.Replace("Content-Type", "text/plain")

    w.WriteStatusLine(response.StatusOk)
    w.WriteHeaders(*h)
    w.WriteBody(body)
})
defer server.Close()
```

---

## 📚 Learning Value

This project is not meant for production use. Instead, it is a **step-by-step exploration** of how HTTP really works:

* From raw **TCP streams** → structured **HTTP requests**
* From **status lines** → full **HTTP responses**
* From simple **listeners** → a mini **HTTP server**

---

## 📜 License

MIT License © 2025 [suhasdeveloper07](https://github.com/suhasdeveloper07)

---
