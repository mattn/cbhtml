package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	GMEM_MOVEABLE = 0x0002
	GMEM_DDESHARE = 0x2000
)

var (
	modkernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procGlobalAlloc             = modkernel32.NewProc("GlobalAlloc")
	procGlobalLock              = modkernel32.NewProc("GlobalLock")
	procGlobalUnlock            = modkernel32.NewProc("GlobalUnlock")
	procGlobalFree              = modkernel32.NewProc("GlobalFree")
	moduser32                   = syscall.NewLazyDLL("user32.dll")
	procOpenClipboard           = moduser32.NewProc("OpenClipboard")
	procCloseClipboard          = moduser32.NewProc("CloseClipboard")
	procSetClipboardData        = moduser32.NewProc("SetClipboardData")
	procRegisterClipboardFormat = moduser32.NewProc("RegisterClipboardFormatW")
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	t, _, err := procRegisterClipboardFormat.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("HTML Format"))))
	if t == 0 {
		log.Fatal(err)
	}

	r, _, err := procOpenClipboard.Call(0)
	if r == 0 {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	header := "Version:0.9\r\n" +
		"StartHTML:%010d\r\n" +
		"EndHTML:%010d\r\n" +
		"StartFragment:%010d\r\n" +
		"EndFragment:%010d\r\n" +
		"<html><body><!--StartFragment -->\r\n"

	footer := "\r\n<!--EndFragment--></body</html>"

	buf.WriteString(fmt.Sprintf(header, 140, 140+len(b), 140, 140+len(b)))
	buf.Write(b)
	buf.WriteString(footer)

	r, _, err = procGlobalAlloc.Call(GMEM_MOVEABLE|GMEM_DDESHARE, uintptr(buf.Len()+4))
	if r == 0 {
		log.Fatal(err)
	}
	h, _, err := procGlobalLock.Call(r)
	if h == 0 {
		log.Fatal(err)
	}
	defer procGlobalUnlock.Call(h)

	p := (*[1<<50 - 1]byte)(unsafe.Pointer(h))
	copy(p[:buf.Len()], buf.Bytes())
	r, _, err = procSetClipboardData.Call(t, h)
	if r == 0 {
		log.Fatal(err)
	}
}
