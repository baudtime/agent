package inputs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/baudtime/agent/plugin"
	. "github.com/baudtime/agent/vars"
	"github.com/baudtime/baudtime/msg"
	"github.com/go-kit/kit/log/level"
	"go.uber.org/multierr"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	fsUsagePercentMetric       = "fs_usage_percent"
	fsUsageBytesMetric         = "fs_usage_bytes"
	fsInodesUsagePercentMetric = "fs_inodes_usage_percent"
	fsInodesUsageNumMetric     = "fs_inodes_usage_num"
)

const (
	fsMountPointsIgnored = "^/(dev|proc|sys|boot|var/lib/docker/.+)($|/)"
	fsTypesIgnored       = "^(autofs|binfmt_misc|bpf|cgroup|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs|tmpfs|fuse.*)$"
	mountTimeout         = 30 * time.Second
)

var (
	ignoredFsMountPoints = kingpin.Flag("input.filesystem.ignored-mount-points",
		"Regexp of mount points to ignore for filesystem collector.").Default(fsMountPointsIgnored).String()
	ignoredFsTypes = kingpin.Flag("input.filesystem.ignored-fs-types",
		"Regexp of filesystem types to ignore for filesystem collector.").Default(fsTypesIgnored).String()

	filesystemLabelNames = []string{"device", "mountpoint", "fstype"}
	fsPlgLogger          = plugin.Logger("filesystem")
)

func init() {
	register("filesystem", plugin.DefaultDisabled, newFilesystemCollector)
}

type filesystemLabels struct {
	device, mountPoint, fsType, options string
}

type FilesystemStats struct {
	labels filesystemLabels
	statfs *syscall.Statfs_t
}

type filesystemCollector struct {
	ignoredFsMountPointsPattern *regexp.Regexp
	ignoredFsTypesPattern       *regexp.Regexp
}

func newFilesystemCollector() (plugin.Input, error) {
	return &filesystemCollector{
		ignoredFsMountPointsPattern: regexp.MustCompile(*ignoredFsMountPoints),
		ignoredFsTypesPattern:       regexp.MustCompile(*ignoredFsTypes),
	}, nil
}

func (c *filesystemCollector) Collect(ch chan<- plugin.Metric) error {
	allFileSystemStats, err := getFileSystemStats(c.ignoredFsMountPointsPattern, c.ignoredFsTypesPattern)

	for _, stats := range allFileSystemStats {
		if stats.statfs != nil {
			labels := []msg.Label{
				{"device", stats.labels.device},
				{"mount", stats.labels.mountPoint},
				{"type", stats.labels.fsType},
			}

			totalRoot := stats.statfs.Blocks * uint64(stats.statfs.Bsize)
			free := stats.statfs.Bfree * uint64(stats.statfs.Bsize)
			avail := stats.statfs.Bavail * uint64(stats.statfs.Bsize)
			used := totalRoot - free //so, totalRoot = used + free, totalUser = used + avail

			totalUser := used + avail
			if totalUser != 0 {
				usagePercent := (used * 100) / totalUser
				if used%totalUser != 0 {
					usagePercent += 1
				}
				if usagePercent > 100 {
					usagePercent = 100
				}
				ch <- plugin.Metric{Name: fsUsagePercentMetric, Value: float64(usagePercent)}.With(HostLabels...).With(labels...)
			}
			ch <- plugin.Metric{Name: fsUsageBytesMetric, Value: float64(used)}.With(HostLabels...).With(labels...)

			inodesTotal := stats.statfs.Files
			inodesUsed := stats.statfs.Files - stats.statfs.Ffree
			if inodesTotal != 0 {
				inodesUsagePercent := (inodesUsed * 100) / inodesTotal
				if inodesUsed%inodesTotal != 0 {
					inodesUsagePercent += 1
				}
				if inodesUsagePercent > 100 {
					inodesUsagePercent = 100
				}
				ch <- plugin.Metric{Name: fsInodesUsagePercentMetric, Value: float64(inodesUsagePercent)}.With(HostLabels...).With(labels...)
			}
			ch <- plugin.Metric{Name: fsInodesUsageNumMetric, Value: float64(inodesUsed)}.With(HostLabels...).With(labels...)
		}
	}

	return err
}

var stuckMounts = make(map[string]struct{})
var stuckMountsMtx = &sync.Mutex{}

// returns filesystem statfs.
func getFileSystemStats(ignoredMountPointsPattern, ignoredFSTypesPattern *regexp.Regexp) (map[string]*FilesystemStats, error) { //mountPoint ->FilesystemStats
	mps, err := mountPointDetails()
	if err != nil {
		return nil, err
	}

	stats := make(map[string]*FilesystemStats)
	for _, labels := range mps {
		if ignoredMountPointsPattern.MatchString(labels.mountPoint) {
			fsPlgLogger(level.Debug).Log("msg", "ignoring mount point", "mountPoint", labels.mountPoint)
			continue
		}
		if ignoredFSTypesPattern.MatchString(labels.fsType) {
			fsPlgLogger(level.Debug).Log("msg", "ignoring filesystem type", "fsType", labels.fsType)
			continue
		}
		stuckMountsMtx.Lock()
		if _, ok := stuckMounts[labels.mountPoint]; ok {
			fsPlgLogger(level.Warn).Log("msg", "Mount point is in an unresponsive state", "mountPoint", labels.mountPoint)
			stuckMountsMtx.Unlock()
			continue
		}
		stuckMountsMtx.Unlock()

		// The success channel is used do tell the "watcher" that the statfs
		// finished successfully. The channel is closed on success.
		success := make(chan struct{})
		go stuckMountWatcher(labels.mountPoint, success)

		statfs := new(syscall.Statfs_t)
		e := syscall.Statfs(rootfsFilePath(labels.mountPoint), statfs)

		stuckMountsMtx.Lock()
		close(success)
		// If the mount has been marked as stuck, unmark it and log it's recovery.
		if _, ok := stuckMounts[labels.mountPoint]; ok {
			fsPlgLogger(level.Warn).Log("msg", "Mount point has recovered, monitoring will resume", "mountPoint", labels.mountPoint)
			delete(stuckMounts, labels.mountPoint)
		}
		stuckMountsMtx.Unlock()

		if e != nil {
			err = multierr.Append(err, fmt.Errorf("Error on statfs() system call for %q: %s", rootfsFilePath(labels.mountPoint), e))
			continue
		}

		stats[labels.mountPoint] = &FilesystemStats{
			labels: labels,
			statfs: statfs,
		}
	}

	return stats, err
}

// stuckMountWatcher listens on the given success channel and if the channel closes
// then the watcher does nothing. If instead the timeout is reached, the
// mount point that is being watched is marked as stuck.
func stuckMountWatcher(mountPoint string, success chan struct{}) {
	select {
	case <-success:
		// Success
	case <-time.After(mountTimeout):
		// Timed out, mark mount as stuck
		stuckMountsMtx.Lock()
		select {
		case <-success:
			// Success came in just after the timeout was reached, don't label the mount as stuck
		default:
			fsPlgLogger(level.Error).Log("msg", "Mount point timed out, it is being labeled as stuck and will not be monitored", "mountPoint", mountPoint)
			stuckMounts[mountPoint] = struct{}{}
		}
		stuckMountsMtx.Unlock()
	}
}

func mountPointDetails() ([]filesystemLabels, error) {
	file, err := os.Open(procFilePath("1/mounts"))
	if os.IsNotExist(err) {
		// Fallback to `/proc/mounts` if `/proc/1/mounts` is missing due hidepid.
		fsPlgLogger(level.Error).Log("msg", "got err when reading root mounts, falling back to system mounts", "err", err)
		file, err = os.Open(procFilePath("mounts"))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseFilesystemLabels(file)
}

func parseFilesystemLabels(r io.Reader) ([]filesystemLabels, error) {
	var filesystems []filesystemLabels

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		if len(parts) < 4 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}

		// Ensure we handle the translation of \040 and \011
		// as per fstab(5).
		parts[1] = strings.Replace(parts[1], "\\040", " ", -1)
		parts[1] = strings.Replace(parts[1], "\\011", "\t", -1)

		filesystems = append(filesystems, filesystemLabels{
			device:     parts[0],
			mountPoint: parts[1],
			fsType:     parts[2],
			options:    parts[3],
		})
	}

	return filesystems, scanner.Err()
}
