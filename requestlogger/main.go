package requestlogger

import (
	"io"
	"flag"
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
}

// ###########################################################################

// httpLogRequest ...
func (ms *RequestLogger) httpLogRequest(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response microservice.Response
  var maxLen int

	microservice.InitResponseFromMicroService(&response, ms, msg)
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	w.WriteHeader(http.StatusOK)
	contentLen = ms.Reply(w, response)
  ms.GetLogger().Printf("| %-6.6s | %s\n", r.Method, r.URL.String())

  if ms.GetLogHeaders() {

    for name, _ := range r.Header {
      if len(name) > maxLen {
        maxLen = len(name)
      }
    }

    for name, values := range r.Header {
      for _, value := range values {
        ms.GetLogger().Printf("  %6.6s | %-*s = '%s'\n", "H", maxLen, name, value)
      }
    }

  }

  defer r.Body.Close()

  if ms.GetLogPayload() {

    var hasLines bool
    var pairs []string

    body := new(strings.Builder)
    n, _ := io.Copy(body, r.Body)
    str := body.String()

    if n > 0 {

      str = strings.Map(func(r rune) rune {
        if unicode.IsPrint(r) {
          return r
        }
         return rune('?')
      }, str)


      hasLines = true

      for hasLines = true; hasLines; hasLines = len(pairs) > 1 {

        pairs = strings.SplitN(str, "\n", 2)
        ms.GetLogger().Printf("  %6.6s | %s\n", "P", pairs[0])

        if len(pairs) > 1 {
          str = pairs[1]
        } else {
          str = ""
        }
      }
      if len(str) > 0 {
        ms.GetLogger().Printf("  %-6.6s | %s\n", "", str)
      }
    }

  }

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
	microservice.Init(&ms.MicroService, &cfg.Configuration, handler)

	return ms
}

