package inputs

import (
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// The path of the proc filesystem.
	procPath   = kingpin.Flag("path.procfs", "procfs mountpoint.").Default("/proc").String()
	sysPath    = kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	rootfsPath = kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
)

func procFilePath(name string) string {
	return filepath.Join(*procPath, name)
}

func sysFilePath(name string) string {
	return filepath.Join(*sysPath, name)
}

func rootfsFilePath(name string) string {
	return filepath.Join(*rootfsPath, name)
}
