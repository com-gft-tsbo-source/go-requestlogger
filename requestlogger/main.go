package requestlogger

import (
	"flag"
	"net/http"

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
}

// ###########################################################################

// httpLogRequest ...
func (ms *RequestLogger) httpLogRequest(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
  //	msg = fmt.Sprintf("'%s' @ '%s' called.", ms.GetName(), ms.GetVersion())
	var response microservice.Response

	microservice.InitResponseFromMicroService(&response, ms, msg)
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	w.WriteHeader(http.StatusOK)
	contentLen = ms.Reply(w, response)
	return status, contentLen, msg
}
// ###########################################################################

// InitFromArgs ...
func InitFromArgs(ms *RequestLogger, args []string, flagset *flag.FlagSet) *RequestLogger {
	var cfg Configuration

	if flagset == nil {
		flagset = flag.NewFlagSet("requestlogger", flag.PanicOnError)
	}

	InitConfigurationFromArgs(&cfg, args, flagset)
	microservice.Init(&ms.MicroService, &cfg.Configuration, nil)

//	ms.AddHandler("/template", templateHandler)

	return ms
}

