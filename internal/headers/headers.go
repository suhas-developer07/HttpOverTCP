package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed filed line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name,[]byte(" ")){
		return "","", fmt.Errorf("malformed filed name")
	}
	return string(name),string(value),nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:],rn)
		if idx == -1 {
			break
		}

		// EMPTY HEADER
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		fmt.Printf("header :\"%s\n",string(data[read:idx]))

		name, value, err := parseHeader(data[read:read+idx])
		if err  != nil {
			return 0,false, err
		}

		read += idx + len(rn)
		h[name] = value
	}
	return read, done,nil
}
