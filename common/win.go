//go:build windows

package common

import (
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32  = syscall.NewLazyDLL("user32.dll")
	kernel  = syscall.NewLazyDLL("kernel32.dll")
	dwmapi  = syscall.NewLazyDLL("dwmapi.dll")
	pGetFG  = user32.NewProc("GetForegroundWindow")
	pGetWin = user32.NewProc("GetWindow")
	pVis    = user32.NewProc("IsWindowVisible")
	pIconic = user32.NewProc("IsIconic")
	pTxtLen = user32.NewProc("GetWindowTextLengthW")
	pTxt    = user32.NewProc("GetWindowTextW")
	pPid    = user32.NewProc("GetWindowThreadProcessId")
	pShow   = user32.NewProc("ShowWindow")
	pPost   = user32.NewProc("PostMessageW")
	pOpen   = kernel.NewProc("OpenProcess")
	pClose  = kernel.NewProc("CloseHandle")
	pQName  = kernel.NewProc("QueryFullProcessImageNameW")
	pDwmGet = dwmapi.NewProc("DwmGetWindowAttribute")
)

func procName(h uintptr) string {
	var pid uint32
	pPid.Call(h, uintptr(unsafe.Pointer(&pid)))
	ph, _, _ := pOpen.Call(0x1000, 0, uintptr(pid)) // PROCESS_QUERY_LIMITED_INFORMATION
	if ph == 0 {
		return ""
	}
	defer pClose.Call(ph)
	buf := make([]uint16, 512)
	n := uint32(len(buf))
	if r, _, _ := pQName.Call(ph, 0, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n))); r == 0 {
		return ""
	}
	return strings.ToLower(filepath.Base(syscall.UTF16ToString(buf[:n])))
}

func ok(h uintptr) bool {
	if v, _, _ := pVis.Call(h); v == 0 {
		return false
	}
	if i, _, _ := pIconic.Call(h); i != 0 {
		return false
	}
	l, _, _ := pTxtLen.Call(h)
	if l == 0 {
		return false
	}
	var cloaked int32
	pDwmGet.Call(h, 14, uintptr(unsafe.Pointer(&cloaked)), 4)
	if cloaked != 0 {
		return false
	}
	buf := make([]uint16, l+1)
	pTxt.Call(h, uintptr(unsafe.Pointer(&buf[0])), l+1)
	title := syscall.UTF16ToString(buf)
	if title == "Program Manager" {
		return false
	}
	// Skip Vicinae's launcher, but allow its Settings window.
	if strings.HasPrefix(procName(h), "vicinae") {
		return strings.HasPrefix(title, "Vicinae Settings")
	}
	return true
}

// Target returns the active window, skipping Vicinae's own windows.
func Target() uintptr {
	h, _, _ := pGetFG.Call()
	for h != 0 && !ok(h) {
		h, _, _ = pGetWin.Call(h, 2) // GW_HWNDNEXT
	}
	return h
}

func Maximize(h uintptr) { pShow.Call(h, 3) }            // SW_MAXIMIZE
func Close(h uintptr)    { pPost.Call(h, 0x0010, 0, 0) } // WM_CLOSE
