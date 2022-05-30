package requestlogger

import (
	"flag"
	"os"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/microservice"
)

// Configuration ...
type Configuration struct {
	microservice.Configuration
}

// IConfiguration ...
type IConfiguration interface {
	microservice.IConfiguration
}

// ---------------------------------------------------------------------------

// InitConfigurationFromArgs ...
func InitConfigurationFromArgs(cfg *Configuration, args []string, flagset *flag.FlagSet) {
	if flagset == nil {
		flagset = flag.NewFlagSet("requestlogger", flag.PanicOnError)
	}

	microservice.InitConfigurationFromArgs(&cfg.Configuration, args, flagset)
	flagset.Parse(os.Args[1:])

}
