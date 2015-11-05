package s3signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
)

func sign(stringToSign string, keys Credentials) string {
	hash := hmac.New(sha1.New, []byte(keys.SecretAccessKey))
	hash.Write([]byte(stringToSign))
	signature := make([]byte, base64.StdEncoding.EncodedLen(hash.Size()))
	base64.StdEncoding.Encode(signature, hash.Sum(nil))
	return string(signature)
}

func stringToSign(request *http.Request) string {
	var buffer bytes.Buffer
	buffer.WriteString(request.Method + "\n")
	buffer.WriteString(request.Header.Get("Content-MD5") + "\n")
	buffer.WriteString(request.Header.Get("Content-Type") + "\n")
}

func canonicalAmzHeaders(request *http.Request) string {
	var headers []string

	for header := range request.Header {
		standardized := strings.ToLower(strings.TrimSpace(header))
		if strings.HasPrefix(standardized, "x-amz") {
			headers = append(headers, standardized)
		}
	}

	sort.Strings(headers)

	//TODO(jstackhouse): Combine headers into header-name:csv
	//  Combine header fields with the same name into one
	//  "header-name:comma-separated-value-list" pair as
	//  prescribed by RFC 2616, section 4.2, without any
	//  whitespace between values. For example, the two
	//  metadata headers 'x-amz-meta-username: fred' and
	//  'x-amz-meta-username: barney' would be combined
	//  into the single header 'x-amz-meta-username: fred,barney'.

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	}
	return ""
}

func canonicalResource(request *http.Request) string {

}
