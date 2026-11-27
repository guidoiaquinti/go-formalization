package islocaldevice

import (
	"runtime"
	"syscall"
	"unsafe"
)

var proc = syscall.NewLazyDLL("ntdll.dll").NewProc("RtlIsDosDeviceName_U")

func IsDOSName(path string) bool {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}
	result, _, _ := proc.Call(uintptr(unsafe.Pointer(name)))
	runtime.KeepAlive(name)
	return uint32(result) != 0
}
