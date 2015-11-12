package s3signer

import (
	"net/http"
	"testing"
)

func NewRequestWithHeaders(method string, path string, headers ...string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	for i := 0; i < len(headers)-1; i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}
	return req
}

func TestSign(t *testing.T) {
	credentials = &Credentials{
		AccessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	tests := []struct {
		ExpectedSTS  string
		ExpectedAuth string
		Request      *http.Request
	}{
		{
			"GET\n\n\nTue, 27 Mar 2007 19:36:42 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:bWq2s1WEIj+Ydj0vQ697zp+IXMU=",
			NewRequestWithHeaders(
				"GET",
				"http://johnsmith.s3.amazonaws.com/photos/puppy.jpg",
				"Host", "johnsmith.s3.amazonaws.com",
				"Date", "Tue, 27 Mar 2007 19:36:42 +0000",
			),
		},
		{
			"PUT\n\nimage/jpeg\nTue, 27 Mar 2007 21:15:45 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:MyyxeRY7whkBe+bq8fHCL/2kKUg=",
			NewRequestWithHeaders(
				"PUT",
				"http://johnsmith.s3.amazonaws.com/photos/puppy.jpg",
				"Host", "johnsmith.s3.amazonaws.com",
				"Content-Type", "image/jpeg",
				"Content-Length", "94328",
				"Date", "Tue, 27 Mar 2007 21:15:45 +0000",
			),
		},
		{
			"GET\n\n\nTue, 27 Mar 2007 19:42:41 +0000\n/johnsmith/",
			"AWS AKIAIOSFODNN7EXAMPLE:htDYFYduRNen8P9ZfE/s9SuKy0U=",
			NewRequestWithHeaders(
				"GET",
				"http://johnsmith.s3.amazonaws.com/?prefix=photos&max-keys=50&marker=puppy",
				"Host", "johnsmith.s3.amazonaws.com",
				"Date", "Tue, 27 Mar 2007 19:42:41 +0000",
			),
		},
		{
			"GET\n\n\nTue, 27 Mar 2007 19:44:46 +0000\n/johnsmith/?acl",
			"AWS AKIAIOSFODNN7EXAMPLE:c2WLPFtWHVgbEmeEG93a4cG37dM=",
			NewRequestWithHeaders(
				"GET",
				"http://johnsmith.s3.amazonaws.com/?acl",
				"Host", "johnsmith.s3.amazonaws.com",
				"Date", "Tue, 27 Mar 2007 19:44:46 +0000",
			),
		},
		{
			"DELETE\n\n\nTue, 27 Mar 2007 21:20:26 +0000\n/johnsmith/photos/puppy.jpg",
			"AWS AKIAIOSFODNN7EXAMPLE:lx3byBScXR6KzyMaifNkardMwNk=",
			NewRequestWithHeaders(
				"DELETE",
				"http://s3.amazonaws.com/johnsmith/photos/puppy.jpg",
				"Host", "s3.amazonaws.com",
				"User-Agent", "dotnet",
				"Date", "Tue, 27 Mar 2007 21:20:27 +0000",
				"x-amz-date", "Tue, 27 Mar 2007 21:20:26 +0000",
			),
		},
		{
			"PUT\n4gJE4saaMU4BqNR0kLY+lw==\napplication/x-download\nTue, 27 Mar 2007 21:06:08 +0000\nx-amz-acl:public-read\nx-amz-meta-checksumalgorithm:crc32\nx-amz-meta-filechecksum:0x02661779\nx-amz-meta-reviewedby:joe@johnsmith.net,jane@johnsmith.net\n/static.johnsmith.net/db-backup.dat.gz",
			"AWS AKIAIOSFODNN7EXAMPLE:ilyl83RwaSoYIEdixDQcA4OnAnc=",
			NewRequestWithHeaders(
				"PUT",
				"http://static.johnsmith.net:8080/db-backup.dat.gz",
				"Host", "static.johnsmith.net",
				"User-Agent", "curl/7.15.5",
				"Date", "Tue, 27 Mar 2007 21:06:08 +0000",
				"x-amz-acl", "public-read",
				"content-type", "application/x-download",
				"Content-MD5", "4gJE4saaMU4BqNR0kLY+lw==",
				"X-Amz-Meta-ReviewedBy", "joe@johnsmith.net",
				"X-Amz-Meta-ReviewedBy", "jane@johnsmith.net",
				"X-Amz-Meta-FileChecksum", "0x02661779",
				"X-Amz-Meta-ChecksumAlgorithm", "crc32",
				"Content-Disposition", "attachment; filename=database.dat",
				"Content-Encoding", "gzip",
				"Content-Length", "5913339",
			),
		},
		{
			"GET\n\n\nWed, 28 Mar 2007 01:29:59 +0000\n/",
			"AWS AKIAIOSFODNN7EXAMPLE:qGdzdERIC03wnaRNKh6OqZehG9s=",
			NewRequestWithHeaders(
				"GET",
				"http://s3.amazonaws.com/",
				"Host", "s3.amazonaws.com",
				"Date", "Wed, 28 Mar 2007 01:29:59 +0000",
			),
		},
		{
			"GET\n\n\nWed, 28 Mar 2007 01:49:49 +0000\n/dictionary/fran%C3%A7ais/pr%c3%a9f%c3%a8re",
			"AWS AKIAIOSFODNN7EXAMPLE:DNEZGsoieTZ92F3bUfSPQcbGmlM=",
			NewRequestWithHeaders(
				"GET",
				"http://s3.amazonaws.com/dictionary/fran%C3%A7ais/pr%c3%a9f%c3%a8re",
				"Host", "s3.amazonaws.com",
				"Date", "Wed, 28 Mar 2007 01:49:49 +0000",
			),
		},
		{
			"PUT\n\nmultipart/form-data\n\nx-amz-date:Thu, 12 Nov 2015 15:15:27 +0000\n/uploads/test%20ingest.mp4?partNumber=1&uploadId=1DYY7bPzrkrXv3RdX9EPpUfmjx8ihEpGtCUE_F0Jl.dlLvjFytmlckdulxq1M0q7dDxRtLfn54JIOgPTCOvdJHF0C8jIA.Qt34CngfTB.02i_XE5_1B0NToWyDgJ_nRr",
			"AWS AKIAIOSFODNN7EXAMPLE:mDjCUH9TwVNGB8NEg9hSK8IbkN4=",
			NewRequestWithHeaders(
				"PUT",
				"https://s3.amazonaws.com/uploads/test ingest.mp4?partNumber=1&uploadId=1DYY7bPzrkrXv3RdX9EPpUfmjx8ihEpGtCUE_F0Jl.dlLvjFytmlckdulxq1M0q7dDxRtLfn54JIOgPTCOvdJHF0C8jIA.Qt34CngfTB.02i_XE5_1B0NToWyDgJ_nRr",
				"Host", "s3.amazonaws.com",
				"Content-Type", "multipart/form-data",
				"X-Amz-Date", "Thu, 12 Nov 2015 15:15:27 +0000",
			),
		},
	}

	for _, test := range tests {
		actual := stringToSign(test.Request)
		if actual != test.ExpectedSTS {
			t.Errorf("\n%s\n!=\n%s", actual, test.ExpectedSTS)
		}

		SignV2WithCredentials(test.Request, credentials)
		if test.Request.Header.Get("Authorization") != test.ExpectedAuth {
			t.Errorf("Authorization doesn't match")
		}
	}
}
