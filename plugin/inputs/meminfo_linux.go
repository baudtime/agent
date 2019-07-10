package inputs

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/baudtime/agent/plugin"
	. "github.com/baudtime/agent/vars"
)

const (
	memTotalMetric         = "mem_total"
	memFreeMetric          = "mem_free"
	memBuffersMetric       = "mem_buffers"
	memCachedMetric        = "mem_cached"
	memActiveAnonMetric    = "mem_active_anon"
	memInactiveAnonMetric  = "mem_inactive_anon"
	memActiveFileMetric    = "mem_active_file"
	memInactiveFileMetric  = "mem_inactive_file"
	memShmemMetric         = "mem_shmem"
	memSlabMetric          = "mem_slab"
	memSReclaimableMetric  = "mem_s_reclaimable"
	memSUnreclaimMetric    = "mem_s_unreclaim"
	swapTotalMetric        = "swap_total"
	swapFreeMetric         = "swap_free"
	memUsageBytesMetric    = "mem_usage_bytes"
	memUsagePercentMetric  = "mem_usage_percent"
	swapUsageBytesMetric   = "swap_usage_bytes"
	swapUsagePercentMetric = "swap_usage_percent"
)

func init() {
	register("meminfo", plugin.DefaultDisabled, newMeminfoCollector)
}

type MemInfo struct {
	MemTotal          uint64 `json:"mem_total"`
	MemFree           uint64 `json:"mem_free"`
	MemAvailable      uint64 `json:"mem_available"`
	Buffers           uint64 `json:"buffers"`
	Cached            uint64 `json:"cached"`
	SwapCached        uint64 `json:"swap_cached"`
	Active            uint64 `json:"active"`
	Inactive          uint64 `json:"inactive"`
	ActiveAnon        uint64 `json:"active_anon" field:"Active(anon)"`
	InactiveAnon      uint64 `json:"inactive_anon" field:"Inactive(anon)"`
	ActiveFile        uint64 `json:"active_file" field:"Active(file)"`
	InactiveFile      uint64 `json:"inactive_file" field:"Inactive(file)"`
	Unevictable       uint64 `json:"unevictable"`
	Mlocked           uint64 `json:"mlocked"`
	SwapTotal         uint64 `json:"swap_total"`
	SwapFree          uint64 `json:"swap_free"`
	Dirty             uint64 `json:"dirty"`
	Writeback         uint64 `json:"write_back"`
	AnonPages         uint64 `json:"anon_pages"`
	Mapped            uint64 `json:"mapped"`
	Shmem             uint64 `json:"shmem"`
	Slab              uint64 `json:"slab"`
	SReclaimable      uint64 `json:"s_reclaimable"`
	SUnreclaim        uint64 `json:"s_unclaim"`
	KernelStack       uint64 `json:"kernel_stack"`
	PageTables        uint64 `json:"page_tables"`
	NFS_Unstable      uint64 `json:"nfs_unstable"`
	Bounce            uint64 `json:"bounce"`
	WritebackTmp      uint64 `json:"writeback_tmp"`
	CommitLimit       uint64 `json:"commit_limit"`
	Committed_AS      uint64 `json:"committed_as"`
	VmallocTotal      uint64 `json:"vmalloc_total"`
	VmallocUsed       uint64 `json:"vmalloc_used"`
	VmallocChunk      uint64 `json:"vmalloc_chunk"`
	HardwareCorrupted uint64 `json:"hardware_corrupted"`
	AnonHugePages     uint64 `json:"anon_huge_pages"`
	HugePages_Total   uint64 `json:"huge_pages_total"`
	HugePages_Free    uint64 `json:"huge_pages_free"`
	HugePages_Rsvd    uint64 `json:"huge_pages_rsvd"`
	HugePages_Surp    uint64 `json:"huge_pages_surp"`
	Hugepagesize      uint64 `json:"hugepagesize"`
	DirectMap4k       uint64 `json:"direct_map_4k"`
	DirectMap2M       uint64 `json:"direct_map_2M"`
	DirectMap1G       uint64 `json:"direct_map_1G"`
}

type meminfoCollector struct{}

func newMeminfoCollector() (plugin.Input, error) {
	return &meminfoCollector{}, nil
}

func (c *meminfoCollector) Collect(ch chan<- plugin.Metric) error {
	mi, err := readMemInfo(procFilePath("meminfo"))
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %s", err)
	}

	ch <- plugin.Metric{Name: memTotalMetric, Value: float64(mi.MemTotal)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memFreeMetric, Value: float64(mi.MemFree)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memBuffersMetric, Value: float64(mi.Buffers)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memCachedMetric, Value: float64(mi.Cached)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memActiveAnonMetric, Value: float64(mi.ActiveAnon)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memInactiveAnonMetric, Value: float64(mi.InactiveAnon)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memActiveFileMetric, Value: float64(mi.ActiveFile)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memInactiveFileMetric, Value: float64(mi.InactiveFile)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memShmemMetric, Value: float64(mi.Shmem)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memSlabMetric, Value: float64(mi.Slab)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memSReclaimableMetric, Value: float64(mi.SReclaimable)}.With(HostLabels...)
	ch <- plugin.Metric{Name: memSUnreclaimMetric, Value: float64(mi.SUnreclaim)}.With(HostLabels...)
	ch <- plugin.Metric{Name: swapTotalMetric, Value: float64(mi.SwapTotal)}.With(HostLabels...)
	ch <- plugin.Metric{Name: swapFreeMetric, Value: float64(mi.SwapFree)}.With(HostLabels...)

	memUsed := mi.MemTotal - mi.MemFree - mi.Cached + mi.Shmem
	ch <- plugin.Metric{Name: memUsageBytesMetric, Value: float64(memUsed)}.With(HostLabels...)

	if mi.MemTotal != 0 {
		memUsagePercent := memUsed * 100 / mi.MemTotal
		if memUsagePercent > 100 {
			memUsagePercent = 100
		}
		ch <- plugin.Metric{Name: memUsagePercentMetric, Value: float64(memUsagePercent)}.With(HostLabels...)
	}

	if mi.SwapTotal != 0 {
		swapUsed := mi.SwapTotal - mi.SwapFree
		swapUsagePercent := swapUsed * 100 / mi.SwapTotal
		ch <- plugin.Metric{Name: swapUsageBytesMetric, Value: float64(swapUsed)}.With(HostLabels...)
		ch <- plugin.Metric{Name: swapUsagePercentMetric, Value: float64(swapUsagePercent)}.With(HostLabels...)
	}

	return nil
}

func readMemInfo(path string) (*MemInfo, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	// Maps a meminfo metric to its value (i.e. MemTotal --> 100000)
	statMap := make(map[string]uint64)

	var info = MemInfo{}

	for _, line := range lines {
		fields := strings.SplitN(line, ":", 2)
		if len(fields) < 2 {
			continue
		}
		valFields := strings.Fields(fields[1])
		val, _ := strconv.ParseUint(valFields[0], 10, 64)
		statMap[fields[0]] = val
	}

	elem := reflect.ValueOf(&info).Elem()
	typeOfElem := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		val, ok := statMap[typeOfElem.Field(i).Name]
		if ok {
			elem.Field(i).SetUint(val)
			continue
		}
		val, ok = statMap[typeOfElem.Field(i).Tag.Get("field")]
		if ok {
			elem.Field(i).SetUint(val)
		}
	}

	return &info, nil
}
