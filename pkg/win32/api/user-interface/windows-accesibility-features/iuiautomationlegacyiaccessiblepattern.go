// +build windows

package windows_accesibility_features

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

type IUIAutomationLegacyIAccessiblePattern struct {
	ole.IUnknown
}

type IUIAutomationLegacyIAccessiblePatternVtbl struct {
	ole.IUnknownVtbl
	SetValue uintptr
}

// https://github.com/mmarquee/ui-automation/blob/ec43c1449b11b5d0f3fd313367e242c6ce456bd9/src/main/java/mmarquee/uiautomation/IUIAutomationLegacyIAccessiblePattern.java
var IID_IUIAutomationLegacyIAccessiblePattern = &ole.GUID{0x828055ad, 0x355b, 0x4435, [8]byte{0x86, 0xd5, 0x3b, 0x51, 0xc1, 0x4a, 0x9b, 0x1b}}

func (pat *IUIAutomationLegacyIAccessiblePattern) VTable() *IUIAutomationLegacyIAccessiblePatternVtbl {
	return (*IUIAutomationLegacyIAccessiblePatternVtbl)(unsafe.Pointer(pat.RawVTable))
}

// https://docs.microsoft.com/en-us/windows/win32/api/uiautomationclient/nf-uiautomationclient-iuiautomationlegacyiaccessiblepattern-setvalue
// HRESULT SetValue(
// 	LPCWSTR szValue
// );
func (pat *IUIAutomationLegacyIAccessiblePattern) SetValue(value string) error {

	szValue, err := syscall.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	hr, _, _ := syscall.Syscall(
		pat.VTable().SetValue,
		2,
		uintptr(unsafe.Pointer(pat)),
		uintptr(unsafe.Pointer(szValue)),
		0)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}
