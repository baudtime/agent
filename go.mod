module github.com/baudtime/agent

go 1.13

replace (
	git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git => github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
)

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Sirupsen/logrus v1.4.0 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4 // indirect
	github.com/baudtime/baudtime v0.1.3
	github.com/beevik/ntp v0.2.0
	github.com/cespare/xxhash/v2 v2.1.0
	github.com/docker/docker v1.13.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/go-kit/kit v0.9.0
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/mdlayher/taskstats v0.0.0-20190313225729-7cbba52ee072
	github.com/opencontainers/runc v0.1.1 // v1.0.0-rc9
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/vishvananda/netlink v1.0.0 // indirect
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f // indirect
	go.uber.org/multierr v1.1.0
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
