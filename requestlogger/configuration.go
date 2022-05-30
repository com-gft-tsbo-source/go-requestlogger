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
  LogHeaders bool
  LogPayload bool
  EncodePayload bool
}

// ILoggingConfiguration ...
type ILoggingConfiguration interface {
	GetLineLength() int
	SetLineLength(int)
	GetLogHeaders() bool
	SetLogHeaders(bool)
	GetLogPayload() bool
	SetLogPayload(bool)
	GetEncodePayload() bool
	SetEncodePayload(bool)
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

// GetLogHeaders ...
func (cfg *LoggingConfiguration) GetLogHeaders() bool { return cfg.LogHeaders }

// SetLogHeaders ...
func (cfg *LoggingConfiguration) SetLogHeaders(v bool) { cfg.LogHeaders = v }

// GetLogPayload ...
func (cfg *LoggingConfiguration) GetLogPayload() bool { return cfg.LogPayload }

// SetLogPayload ...
func (cfg *LoggingConfiguration) SetLogPayload(v bool) { cfg.LogPayload = v }

// GetEncodePayload ...
func (cfg *LoggingConfiguration) GetEncodePayload() bool { return cfg.EncodePayload }

// SetEncodePayload ...
func (cfg *LoggingConfiguration) SetEncodePayload(v bool) { cfg.EncodePayload = v }

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
	plogHeaders := flagset.Bool("logHeaders", false, "Also log headers.")
	plogPayload := flagset.Bool("logPayload", false, "Also log payload.")
	pisProxy := flagset.Bool("isProxy", false, "Also log payload.")

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
		cfg.SetLogHeaders(*plogHeaders)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGHEADERS")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogHeaders(v != 0)
		} else {
			cfg.SetLogHeaders(false)
		}
	}

	if *plogPayload {
		cfg.SetLogPayload(*plogPayload)
	} else {
		ev := os.Getenv("REQUESTLOGGER_LOGPAYLOAD")
		if len(ev) > 0 {
			v, err := strconv.Atoi(ev)
			if err != nil {
				panic(err)
			}
			cfg.SetLogPayload(v != 0)
		} else {
			cfg.SetLogPayload(false)
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
		} else {
			cfg.SetIsProxy(false)
		}
	}
}

