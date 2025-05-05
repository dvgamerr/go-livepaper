package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Resolution struct {
	width  int32
	height int32
	x      int32
	y      int32
}

type Rectangle struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor Rectangle
	RcWork    Rectangle
	DwFlags   uint32
}

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	enumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	getMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	enumDisplayDevicesW = user32.NewProc("EnumDisplayDevicesW")
)

type MonitorInfo struct {
	index      int
	resolution Resolution
	rectangle  Rectangle
	primary    bool
}

func getMonitorInfo(hMonitor windows.Handle, index int) (MonitorInfo, error) {
	info := MonitorInfo{index: index}

	var mi MONITORINFO
	mi.CbSize = uint32(unsafe.Sizeof(mi))

	ret, _, _ := getMonitorInfoW.Call(
		uintptr(hMonitor),
		uintptr(unsafe.Pointer(&mi)),
	)

	if ret == 0 {
		return info, fmt.Errorf("GetMonitorInfoW failed")
	}

	width := mi.RcMonitor.bottom - mi.RcMonitor.top
	height := mi.RcMonitor.right - mi.RcMonitor.left
	info.primary = (mi.DwFlags & 1) != 0 // MONITORINFOF_PRIMARY = 1
	info.rectangle = mi.RcMonitor
	info.resolution = Resolution{
		width:  width,
		height: height,
		x:      mi.RcMonitor.left,
		y:      mi.RcMonitor.top,
	}

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
