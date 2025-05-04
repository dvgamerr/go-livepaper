package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	enumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	getMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	enumDisplayDevicesW = user32.NewProc("EnumDisplayDevicesW")
)

type MonitorInfo struct {
	Index      int
	DeviceName string
	Resolution string
	Position   string
	Primary    bool
}

func getMonitorInfo(hMonitor windows.Handle, index int) (MonitorInfo, error) {
	info := MonitorInfo{Index: index}

	var mi MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	ret, _, _ := getMonitorInfoW.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(&mi)),
	)

	if ret == 0 {
		return info, fmt.Errorf("GetMonitorInfoW failed")
	}

	width := mi.RcMonitor.Right - mi.RcMonitor.Left
	height := mi.RcMonitor.Bottom - mi.RcMonitor.Top
	info.Resolution = fmt.Sprintf("%dx%d", width, height)
	info.Position = fmt.Sprintf("(%d,%d)", mi.RcMonitor.Left, mi.RcMonitor.Top)
	info.Primary = (mi.DwFlags & 1) != 0 // MONITORINFOF_PRIMARY = 1

	return info, nil
}

func monitorEnumProc(hMonitor, hdcMonitor, lprcMonitor uintptr, dwData *uintptr) uintptr {
	// Convert uintptr to unsafe.Pointer according to the rules
	monitors := (*[]MonitorInfo)(unsafe.Pointer(dwData))
	info, err := getMonitorInfo(windows.Handle(hMonitor), len(*monitors))
	if err == nil {
		*monitors = append(*monitors, info)
	}
	return 1 // Continue enumeration
}

func getMonitors() []MonitorInfo {
	var monitors []MonitorInfo
	enumDisplayMonitors.Call(
		0,
		0,
		windows.NewCallback(monitorEnumProc),
		uintptr(unsafe.Pointer(&monitors)),
	)
	return monitors
}
