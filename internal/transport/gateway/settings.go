package gateway

import "net/http"

// Default allowed request headers.
var defaultAllowedHTTPHeaders = []string{
	"Accept",
	"Accept-Language",
	"Accept-Encoding",
	"Content-Type",
	"Content-Language",
	"Origin",
	"Authorization",
	"User-Agent",
	// "X-CSRF-Token",
	"Strict-Transport-Security",
	"X-Content-Type-Options",
	"Content-Security-Policy",
	"X-Requested-With",
	"Cache-Control",
	"Content-Length",
	"Cookie",
	"Host",
	"Pragma",
	"Referer",
}

// Default allowed request http methods, also using for http gateway.
var defaultAllowedHTTPMethods = []string{
	http.MethodPost,
}
