# S3 Signer [![GoDoc](https://godoc.org/github.com/itsjamie/s3-signer?status.svg)](https://godoc.org/github.com/itsjamie/s3-signer) [![Build Status](https://travis-ci.org/itsjamie/s3-signer.svg?branch=master)](https://travis-ci.org/itsjamie/s3-signer) [![Coverage Status](https://coveralls.io/repos/itsjamie/s3-signer/badge.svg?branch=master)](https://coveralls.io/r/itsjamie/s3-signer?branch=master)

s3signer is a package that implements the Amazon S3 Signature Version 2

## Getting Started
Construct the http.Request that you intend to send to the Amazon API.
Simply pass that and a populated credentials struct to `s3signer.SignV2()`

## Resources
* [Signing and Authenticating REST Requests](http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html)

## TODO
- [ ] Implement Signature Version 4
- [ ] Build Example Browser-based Multipart Uploads using XHR2
- [ ] Build Example Browser-based Uploads using POST with policies
