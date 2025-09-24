package headers

import (
	"bytes"
	"fmt"
	"strings"
)
//validate header names per HTTP spec

func isToken(str string) bool {
	for _, ch := range str {
		found := false
		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			found = true
		}
		switch ch {
		case '#', '$', '&', '%', '\'',  '*', '+', '-', '.', '^', '_','`','|', '~':
			found = true
		}
		if !found {
			return false
		}
	}
	return true
}

var rn = []byte("\r\n")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed filed line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed filed name")
	}
	return string(name), string(value), nil
}

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{}}
}
//manage header values

func (h *Headers) Get(name string) (string,bool) {
	str, ok := h.headers[strings.ToLower(name)]
	return str, ok
}

func (h *Headers) Replace(name, value string){
	name = strings.ToLower(name)
	h.headers[name] = value
}

func (h *Headers) Delete(name string){
	name = strings.ToLower(name)
	delete(h.headers, name)
}

func (h *Headers) Set(name, value string) {
	name = strings.ToLower(name)

	if v, ok := h.headers[name]; ok {
		h.headers[name] = fmt.Sprintf("%s, %s", v, value)
	} else {
		h.headers[name] = value
	}
}

func (h *Headers) ForEach(cb func(n, v string)) {
	for n, v := range h.headers {
		cb(n, v)
	}
}
//converts raw bytes to header map
func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}
		end := read + idx

		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		fmt.Printf("header :\"%s\n", string(data[read:end]))

		name, value, err := parseHeader(data[read:end])
		if err != nil {
			return 0, false, err
		}

		if !isToken(name) {
			return 0, false, fmt.Errorf("malformed header name")
		}
		read += idx + len(rn)
		h.Set(name,value)
	}
	return read, done, nil
}
