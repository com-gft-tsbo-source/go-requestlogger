package requestlogger

import (
  "fmt"
	"io"
	"io/ioutil"
	"flag"
  "net"
	"net/http"
  "strings"
  "unicode"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// ###########################################################################
// ###########################################################################
// RequestLogger
// ###########################################################################
// ###########################################################################



// RequestLogger Encapsulates the requestlogger data
type RequestLogger struct {
	microservice.MicroService
  *LoggingConfiguration
  *ProxyConfiguration
}

// ###########################################################################

var hopHeaders = []string{
    "Connection",
    "Keep-Alive",
    "Proxy-Authenticate",
    "Proxy-Authorization",
    "Te", // canonicalized version of "TE"
    "Trailers",
    "Transfer-Encoding",
    "Upgrade",
}

func delHopHeaders(header http.Header) {
  for _, hopHeader := range hopHeaders {
    header.Del(hopHeader)
  }
}

func copyHeader(src http.Header, dst http.Header) {
  for key, values := range src {
    for _, value := range values {
      dst.Add(key, value)
    }
  }
}

func appendHostToXForwardHeader(header http.Header, host string) {
  if prior, ok := header["X-Forwarded-For"]; ok {
      host = strings.Join(prior, ", ") + ", " + host
   }
  header.Set("X-Forwarded-For", host)
}

// ###########################################################################

func (ms *RequestLogger) logHeaders(prefix string, headers *http.Header) {

  var maxLen int = 0

  for name, _ := range *headers {
    if len(name) > maxLen {
      maxLen = len(name)
    }
  }

  for name, values := range *headers {
    for _, value := range values {
      ms.GetLogger().Printf("%s %6.6s | %-*s = '%s'\n", prefix, "H", maxLen, name, value)
    }
  }
}

// ---------------------------------------------------------------------------

func (ms *RequestLogger) logPayload(prefix string, bodyStr string) {
  var hasLines bool
  var pairs []string

  bodyLen := len(bodyStr)

  if bodyLen == 0 {
    return
  }

  hasLines = true
  str := strings.Clone( bodyStr )

  for hasLines = true; hasLines; hasLines = len(pairs) > 1 {

    pairs = strings.SplitN(str, "\n", 2)

    pairs[0] = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
        return r
      }
       return rune('?')
    }, pairs[0])

    ms.GetLogger().Printf("%s %6.6s | %s\n", prefix, "P", pairs[0])

    if len(pairs) > 1 {
      str = pairs[1]
    } else {
      str = ""
    }
  }

  if len(str) > 0 {
    str = strings.Map(func(r rune) rune {
      if unicode.IsPrint(r) {
        return r
      }
       return rune('?')
    }, str)

    ms.GetLogger().Printf("%s %6.6s | %s\n", prefix, "P", str)
  }
}

// ---------------------------------------------------------------------------

// httpLogRequest ...
func (ms *RequestLogger) httpLogRequest(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response microservice.Response

  defer r.Body.Close()

  var body *strings.Builder
  var bodyStr string
  var bodyLen int64

  if ms.GetLogRequestPayload() || ms.GetIsProxy(){
    body = new(strings.Builder)
    bodyLen, _ = io.Copy(body, r.Body)
    bodyStr = body.String()
  }

  ms.GetLogger().Printf("# %-6.6s | %s\n", r.Method, r.URL.String())

  if ms.GetLogRequestHeaders() {
    ms.logHeaders("<", &r.Header)
  }

  if ms.GetLogRequestPayload() {
    ms.logPayload("<", bodyStr)
  }

  if ! ms.GetIsProxy() {
    microservice.InitResponseFromMicroService(&response, ms, msg)
	  ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
  	w.WriteHeader(http.StatusOK)
	  contentLen = ms.Reply(w, response)
  	return status, contentLen, msg
  }

  HTTPClient := &http.Client{}

  r.RequestURI = ""
  delHopHeaders(r.Header)

  if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
    appendHostToXForwardHeader(r.Header, clientIP)
  }

  r.Body =  ioutil.NopCloser(strings.NewReader(bodyStr))
  clientResponse, err := HTTPClient.Do(r)
  if err != nil {
    msg = fmt.Sprintf("Failed to proxy: %s", err.Error())
    return http.StatusInternalServerError, 0, msg
  }

  delHopHeaders(clientResponse.Header)

  if ms.GetLogAnswerHeaders() {
    ms.logHeaders(">", &clientResponse.Header)
  }

  if ms.GetLogAnswerPayload() {
    body = new(strings.Builder)
    bodyLen, _ = io.Copy(body, clientResponse.Body)
    bodyStr = body.String()
    clientResponse.Body =  ioutil.NopCloser(strings.NewReader(bodyStr))
    ms.logPayload(">", bodyStr)
  }

  copyHeader(clientResponse.Header, w.Header())
  w.WriteHeader(clientResponse.StatusCode)
  bodyLen, _ = io.Copy(w, clientResponse.Body)
  contentLen = int(bodyLen)
  msg = "Proxied"

	return status, contentLen, msg
}
// ###########################################################################

// InitFromArgs ...
func InitFromArgs(ms *RequestLogger, args []string, flagset *flag.FlagSet) *RequestLogger {
	var cfg Configuration

	if flagset == nil {
		flagset = flag.NewFlagSet("requestlogger", flag.PanicOnError)
	}

	handler := ms.DefaultHandler()
  handler.Any = ms.httpLogRequest
  handler.Get = ms.httpLogRequest
  handler.Put = ms.httpLogRequest
  handler.Post = ms.httpLogRequest
  handler.Delete = ms.httpLogRequest
  handler.Head = ms.httpLogRequest
  handler.Connect = ms.httpLogRequest
  handler.Options = ms.httpLogRequest
	InitConfigurationFromArgs(&cfg, args, flagset)
	ms.LoggingConfiguration = &cfg.LoggingConfiguration
	ms.ProxyConfiguration = &cfg.ProxyConfiguration
	microservice.Init(&ms.MicroService, &cfg.Configuration, handler)

	return ms
}

