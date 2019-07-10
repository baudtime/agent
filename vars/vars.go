package vars

import (
	"io"
	glog "log"
	"net"
	"os"

	"github.com/baudtime/baudtime/msg"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger     log.Logger
	LocalIP    string
	HostName   string
	HostLabels []msg.Label
)

func Init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				LocalIP = ipnet.IP.String()
				break
			}
		}
	}

	if ip := net.ParseIP(LocalIP); ip == nil {
		panic("invalid ip address")
	}

	HostName, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	HostLabels = []msg.Label{{"ip", LocalIP}, {"hostname", HostName}}

	logFile := kingpin.Flag("log.file", "logs will be written to this file").String()
	logLevel := kingpin.Flag("log.level", "log level").Default("info").String()
	ConfigFilePath := kingpin.Flag("config", "configure file path").String()

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	var logWriter io.Writer

	if *logFile == "" {
		logWriter = os.Stdout
	} else {
		logWriter = &lumberjack.Logger{
			Filename:   *logFile,
			MaxSize:    256, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		}
	}

	var levelOpt level.Option
	switch *logLevel {
	case "error":
		levelOpt = level.AllowError()
	case "warn":
		levelOpt = level.AllowWarn()
	case "info":
		levelOpt = level.AllowInfo()
	case "debug":
		levelOpt = level.AllowDebug()
		if logWriter != os.Stdout {
			logWriter = io.MultiWriter(logWriter, os.Stdout)
		}
	default:
		levelOpt = level.AllowInfo()
	}

	glog.SetOutput(logWriter)
	Logger = level.NewFilter(log.NewLogfmtLogger(logWriter), levelOpt)
	Logger = log.With(Logger, "time", log.DefaultTimestamp, "caller", log.DefaultCaller)

	if *ConfigFilePath != "" {
		LoadConfig(*ConfigFilePath)
	}
}
