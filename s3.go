package s3signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

var subresourcesS3 = []string{
	"acl",
	"lifecycle",
	"location",
	"logging",
	"notification",
	"partNumber",
	"policy",
	"requestPayment",
	"torrent",
	"uploadId",
	"uploads",
	"versionId",
	"versioning",
	"versions",
	"website",
}

// Credentials stores AWS credentials used for signing
type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

var credentials *Credentials

// SignV2 takes the HTTP request to sign. The Credentials to sign it are loaded once and stored
func SignV2(request *http.Request) {
	if credentials == nil {
		credentials = getCreds()
	}
	SignV2WithCredentials(request, credentials)
}

// SignV2WithCredentials takes the HTTP request to sign. The Credentials are passed to sign it.
func SignV2WithCredentials(request *http.Request, credentials *Credentials) {
	request.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"AWS %s:%s",
			credentials.AccessKeyID,
			signString(stringToSign(request), credentials),
		),
	)
}

func signString(stringToSign string, keys *Credentials) string {
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
	buffer.WriteString(getDateHeader(request) + "\n")
	buffer.WriteString(canonicalAmzHeaders(request))
	buffer.WriteString(canonicalResource(request))
	return buffer.String()
}

func getDateHeader(request *http.Request) string {
	if header := request.Header.Get("x-amz-date"); header != "" {
		return ""
	} else if header := request.Header.Get("Date"); header != "" {
		return header
	} else {
		return time.Now().UTC().Format(time.RFC1123Z)
	}
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

	for i, header := range headers {
		val := strings.Join(request.Header[http.CanonicalHeaderKey(header)], ",")
		headers[i] = header + ":" + strings.Replace(val, "\n", " ", -1)
	}

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	}
	return ""
}

func canonicalResource(request *http.Request) string {
	resource := ""

	// If Bucket in Host header, add it to canonical resource
	host := request.Header.Get("Host")
	if host == "" || host == "s3.amazonaws.com" {
		// Bucket Name First Part of Request URI, will be added by path
	} else if strings.HasSuffix(host, ".s3.amazonaws.com") {
		val := strings.Split(host, ".")
		resource += "/" + strings.Join(val[0:len(val)-3], ".")
	} else {
		resource += "/" + strings.ToLower(host)
	}

	resource += request.URL.EscapedPath()

	subresources := []string{}
	for _, subResource := range subresourcesS3 {
		if strings.Contains(request.URL.RawQuery, subResource) {
			if val := request.URL.Query().Get(subResource); val != "" {
				subresources = append(subresources, subResource+"="+val)
			} else {
				subresources = append(subresources, subResource)
			}
		}
	}
	if len(subresources) > 0 {
		resource += "?" + strings.Join(subresources, "&")
	}

	return resource
}
