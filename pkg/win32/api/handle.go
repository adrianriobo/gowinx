// +build windows
package win32

import "syscall"

var (
	closeHandle = kernel32.MustFindProc("CloseHandle")
)

// https://docs.microsoft.com/en-us/windows/win32/api/handleapi/nf-handleapi-closehandle
// BOOL CloseHandle(
// 	HANDLE hObject
// );
func CloseHandle(hObject syscall.Handle) (success bool, err error) {
	r0, _, e1 := syscall.Syscall(closeHandle.Addr(), 1,
		uintptr(hObject),
		0,
		0)
	success, err = evalSyscallBool(r0, e1)
	return
}
