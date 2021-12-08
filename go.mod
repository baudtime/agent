module github.com/baudtime/agent

go 1.13

replace (
	git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git => github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4 // indirect
	github.com/baudtime/baudtime v0.1.3
	github.com/beevik/ntp v0.2.0
	github.com/cespare/xxhash/v2 v2.1.0
	github.com/go-kit/kit v0.9.0
	github.com/mdlayher/taskstats v0.0.0-20190313225729-7cbba52ee072
	github.com/opencontainers/runc v1.0.3 // v1.0.0-rc9
	go.uber.org/multierr v1.1.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
