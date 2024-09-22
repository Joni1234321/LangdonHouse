// go run main.go # to run this server
package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/kjk/common/siser"
)

// LogReqInfo describes info about HTTP request
type HTTPReqInfo struct {
	method    string
	uri       string
	referer   string
	ipaddr    string
	code      int           // response code, like 200, 404
	size      int64         // number of bytes of the response sent
	duration  time.Duration // how long did it take to
	userAgent string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		io.WriteString(w,
			`<html>
				<body>Background remover</body>
			</html>`)
	case http.MethodPost:
		// TODO: handle upload
		// TODO: redirect to url with new id
		http.Redirect(w, r, "[::1]/newId", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: should be guid I guess
		pathValue := r.PathValue("id")
		io.WriteString(w, fmt.Sprintf(
			`<html>
				<body>%s</body>
			</html>`, pathValue))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Request.RemoteAddress contains port, which we want to remove i.e.:
// "[::1]:58292" => "[::1]"
func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress returns ip address of the client making the request,
// taking into account http proxies
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIP
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			uri:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)

		// this runs handler h and captures information about
		// HTTP request
		m := httpsnoop.CaptureMetrics(h, w, r)

		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration
		logHTTPReq(ri)
	}
	return http.HandlerFunc(fn)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/{id}", handleItem)
	// TODO: add 404 and error handle

	// add log request middleware
	var handler http.Handler = mux
	handler = logRequestHandler(handler)

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}

	return srv
}

func logHTTPReq(ri *HTTPReqInfo) {
	var rec siser.Record
	rec.Name = "httplog"
	rec.Write(
		"method", ri.method,
		"uri", ri.uri,
		"ipaddr", ri.ipaddr,
		"code", strconv.Itoa(ri.code),
		"size", strconv.FormatInt(ri.size, 10),
		"duration", strconv.FormatInt(int64(ri.duration/time.Millisecond), 10),
		"ua", ri.userAgent,
	)
	if ri.referer != "" {
		rec.Write("referer", ri.referer)
	}

	// TODO: log this properly
	fmt.Println(string(rec.Marshal()))
}

func main() {
	server := makeHTTPServer()

	// TODO: handle server error
	_ = server.ListenAndServe()
}
