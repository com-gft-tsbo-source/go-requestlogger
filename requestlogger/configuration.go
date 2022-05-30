package requestlogger

import (
	"flag"
	"os"
  "strconv"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// LoggingConfiguration ...
type LoggingConfiguration struct {
  LineLength int
  LogRequestHeaders bool
  LogRequestPayload bool
  LogAnswerHeaders bool
  LogAnswerPayload bool
}

// ILoggingConfiguration ...
type ILoggingConfiguration interface {
	GetLineLength() int
	SetLineLength(int)
	GetLogRequestHeaders() bool
	SetLogRequestHeaders(bool)
	GetLogRequestPayload() bool
	SetLogRequestPayload(bool)
	GetLogAnswerHeaders() bool
	SetLogAnswerHeaders(bool)
	GetLogAnswerPayload() bool
	SetLogAnswerPayload(bool)
}

// ProxyConfiguration ...
type ProxyConfiguration struct {
  IsProxy bool
}

// ILoggingConfiguration ...
type IProxyConfiguration interface {
	GetIsProxy() bool
	SetIsProxy(bool)
}

// Configuration ...
type Configuration struct {
	microservice.Configuration
  LoggingConfiguration
  ProxyConfiguration
}

// IConfiguration ...
type IConfiguration interface {
	microservice.IConfiguration
  ILoggingConfiguration
  IProxyConfiguration
}

// ---------------------------------------------------------------------------

// GetLineLength ...
func (cfg *LoggingConfiguration) GetLineLength() int { return cfg.LineLength }

// SetLineLength ...
func (cfg *LoggingConfiguration) SetLineLength(v int) { cfg.LineLength = v }

// GetLogRequestHeaders ...
func (cfg *LoggingConfiguration) GetLogRequestHeaders() bool { return cfg.LogRequestHeaders }

// SetLogRequestHeaders ...
func (cfg *LoggingConfiguration) SetLogRequestHeaders(v bool) { cfg.LogRequestHeaders = v }

// GetLogRequestPayload ...
func (cfg *LoggingConfiguration) GetLogRequestPayload() bool { return cfg.LogRequestPayload }

// SetLogRequestPayload ...
func (cfg *LoggingConfiguration) SetLogRequestPayload(v bool) { cfg.LogRequestPayload = v }

// GetLogAnswerHeaders ...
func (cfg *LoggingConfiguration) GetLogAnswerHeaders() bool { return cfg.LogAnswerHeaders }

// SetLogAnswerHeaders ...
func (cfg *LoggingConfiguration) SetLogAnswerHeaders(v bool) { cfg.LogAnswerHeaders = v }

// GetLogAnswerPayload ...
func (cfg *LoggingConfiguration) GetLogAnswerPayload() bool { return cfg.LogAnswerPayload }

// SetLogAnswerPayload ...
func (cfg *LoggingConfiguration) SetLogAnswerPayload(v bool) { cfg.LogAnswerPayload = v }

// GetIsProxy ...
func (cfg *ProxyConfiguration) GetIsProxy() bool { return cfg.IsProxy }

// SetIsProxy ...
func (cfg *ProxyConfiguration) SetIsProxy(v bool) { cfg.IsProxy = v }

// InitConfigurationFromArgs ...
func InitConfigurationFromArgs(cfg *Configuration, args []string, flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.NewFlagSet("requestlogger", flag.PanicOnError)
	}

	plineLength := flagset.Int("lineLength", -1, "Maximal line length of log output.")
	plogHeaders := flagset.Bool("logHeaders", false, "Log request&answer headers.")
	plogPayload := flagset.Bool("logPayload", false, "Log request&answer payload.")
	plogRequestHeaders := flagset.Bool("logRequestHeaders", false, "Log request headers.")
	plogRequestPayload := flagset.Bool("logRequestPayload", false, "Log request payload.")
	plogAnswerHeaders := flagset.Bool("logAnswerHeaders", false, "Log answer headers.")
	plogAnswerPayload := flagset.Bool("logAnswerPayload", false, "Log answer payload.")
	pisProxy := flagset.Bool("isProxy", false, "Act as a prox.")

	microservice.InitConfigurationFromArgs(&cfg.Configuration, args, flagset)
	flagset.Parse(os.Args[1:])

	if *plineLength >= 0 {
		cfg.SetLineLength(*plineLength)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LINELENGTH")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLineLength(v)
		} else {
			cfg.SetLineLength(-1)
		}
	}

	if *plogHeaders {
		cfg.SetLogRequestHeaders(*plogHeaders)
		cfg.SetLogAnswerHeaders(*plogHeaders)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGHEADERS")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogRequestHeaders(v != 0)
			cfg.SetLogAnswerHeaders(v != 0)
		} else {
			cfg.SetLogRequestHeaders(false)
			cfg.SetLogAnswerHeaders(false)
		}
	}

	if *plogPayload {
		cfg.SetLogRequestPayload(*plogPayload)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGPAYLOAD")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogRequestHeaders(v != 0)
			cfg.SetLogAnswerHeaders(v != 0)
		}
  }

	if *plogRequestHeaders {
		cfg.SetLogRequestHeaders(*plogRequestHeaders)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGREQUESTHEADERS")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogRequestHeaders(v != 0)
		}
	}

	if *plogRequestPayload {
		cfg.SetLogRequestPayload(*plogRequestPayload)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGREQUESTPAYLOAD")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogRequestPayload(v != 0)
		}
	}

	if *plogAnswerHeaders {
		cfg.SetLogAnswerHeaders(*plogAnswerHeaders)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGANSWERHEADERS")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogAnswerHeaders(v != 0)
		}
	}

	if *plogAnswerPayload {
		cfg.SetLogAnswerPayload(*plogAnswerPayload)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGANSWERPAYLOAD")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogAnswerPayload(v != 0)
		}
	}

	if *pisProxy {
		cfg.SetIsProxy(*pisProxy)
	} else {
		ev := os.Getenv("REQUESTLOGGER_ISPROXY")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetIsProxy(v != 0)
		}
	}
}

