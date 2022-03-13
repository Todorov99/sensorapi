package sensor

import (
	"context"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"golang.org/x/sys/unix"
)

func darwinArm64InfoWithContext(ctx context.Context) ([]cpu.InfoStat, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return darwinArm64Info()
	}
}

func darwinArm64Info() ([]cpu.InfoStat, error) {
	var ret []cpu.InfoStat

	c := cpu.InfoStat{}
	c.ModelName, _ = unix.Sysctl("machdep.cpu.brand_string")
	family, _ := unix.SysctlUint32("machdep.cpu.family")
	c.Family = strconv.FormatUint(uint64(family), 10)
	model, _ := unix.SysctlUint32("machdep.cpu.model")
	c.Model = strconv.FormatUint(uint64(model), 10)
	stepping, _ := unix.SysctlUint32("machdep.cpu.stepping")
	c.Stepping = int32(stepping)
	features, err := unix.Sysctl("machdep.cpu.features")
	if err == nil {
		for _, v := range strings.Fields(features) {
			c.Flags = append(c.Flags, strings.ToLower(v))
		}
	}
	leaf7Features, err := unix.Sysctl("machdep.cpu.leaf7_features")
	if err == nil {
		for _, v := range strings.Fields(leaf7Features) {
			c.Flags = append(c.Flags, strings.ToLower(v))
		}
	}
	extfeatures, err := unix.Sysctl("machdep.cpu.extfeatures")
	if err == nil {
		for _, v := range strings.Fields(extfeatures) {
			c.Flags = append(c.Flags, strings.ToLower(v))
		}
	}
	cores, _ := unix.SysctlUint32("machdep.cpu.core_count")
	c.Cores = int32(cores)
	cacheSize, _ := unix.SysctlUint32("machdep.cpu.cache.size")
	c.CacheSize = int32(cacheSize)
	c.VendorID, _ = unix.Sysctl("machdep.cpu.vendor")

	// Use the rated frequency of the CPU. This is a static value and does not
	// account for low power or Turbo Boost modes.
	cpuFrequency, err := unix.SysctlUint64("hw.tbfrequency")
	if err != nil {
		return ret, err
	}

	c.Mhz = float64(cpuFrequency) / 1000000.0

	return append(ret, c), nil
}
