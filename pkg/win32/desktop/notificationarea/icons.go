// +build windows
package notificationarea

import (
	"fmt"
	"syscall"
	"unsafe"

	win32api "github.com/adrianriobo/gowinx/pkg/win32/api"
)

const MEM_COMMIT = 0x1000
const PAGE_READWRITE = 0x04

func GetHiddenIconsCount() (int32, error) {
	return getIconsCountByWindowClass(NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS)
}

// FIXME how to identify ToolBar32 always instance 3?
// func GetVisibleIconsCount() (int32, error) {
// 	handler, _ := ux.FindWindowByClass(NOTIFICATION_AREA_VISIBLE_WINDOW_CLASS)
// 	toolbarHandler, _ := findElementsbyClass(handler, "ToolbarWindow32")
// 	buttonsCount, _ := win32.SendMessage(toolbarHandler, win32.TB_BUTTONCOUNT, 0, 0)
// 	return int32(buttonsCount), nil
// }

func getIconsCountByWindowClass(className string) (int32, error) {
	var err error
	if toolbarHandler, err := getNotificationAreaToolbarByWindowClass(className); err == nil {
		buttonsCount, _ := win32api.SendMessage(toolbarHandler, win32api.TB_BUTTONCOUNT, 0, 0)
		return int32(buttonsCount), nil
	}
	return 0, err
}

// To implement
// func GetIconPosition(title string) (x, y int32) {
// }

func GetIconPosition(rect win32api.RECT) (x, y int32) {
	x = rect.Left + 10
	y = rect.Top + 10
	fmt.Printf("Crc icon will be clicked at x: %d y: %d\n", x, y)
	return
}

func GetIconByTittle(title string) syscall.Handle {
	toolbarHandlers, _ := findToolbars()
	for i, toolbarHandler := range toolbarHandlers {
		fmt.Printf("Looking for %s at toolbar index %d\n", title, i)
		iconHandler, iconIndex, err := findElementByTitle(toolbarHandler, title)
		if err == nil {
			fmt.Printf("We found the icon for %s at index %d\n", title, iconIndex)
			return iconHandler
		}
	}
	return 0
}

func GetIconRectByTittle(title string) (rect win32api.RECT, err error) {
	toolbarHandlers, _ := findToolbars()
	for i, toolbarHandler := range toolbarHandlers {
		fmt.Printf("Looking for %s at toolbar index %d\n", title, i)
		iconHandler, iconIndex, err := findElementByTitle(toolbarHandler, title)
		if err == nil {
			fmt.Printf("We found the icon for %s at index %d\n", title, iconIndex)
			rect, err = getControlRect(iconHandler)
		}
	}
	return
}

func getControlRect(controlHandler syscall.Handle) (rect win32api.RECT, err error) {
	if _, err = win32api.GetWindowRect(controlHandler, &rect); err == nil {
		fmt.Printf("Rect for control t:%d,l:%d,r:%d,b:%d\n", rect.Top, rect.Left, rect.Right, rect.Bottom)
	} else {
		fmt.Printf("error getting control area rect: %v\n", err)
	}
	return
}

func GetButtonsTexts() {
	// var err error
	if toolbarHandler, err := getNotificationAreaToolbarByWindowClass(NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS); err == nil {
		buttonsCount, _ := win32api.SendMessage(toolbarHandler, win32api.TB_BUTTONCOUNT, 0, 0)
		for i := 0; i < int(buttonsCount); i++ {
			text, _ := getButtonText(toolbarHandler, int32(i))
			fmt.Printf("The name of the button at index %d, is %s\n", i, text)
		}
	}
}

func getButtonText(toolbarHandler syscall.Handle, buttonIndex int32) (text string, err error) {
	rect, _ := getControlRect(toolbarHandler)
	fmt.Printf("Get rect top: %d, left: %d \n", rect.Top, rect.Left)
	var tbProcessID uint32
	toolbarThreadId, _ := win32api.GetWindowThreadProcessId(toolbarHandler, &tbProcessID)
	fmt.Printf("ProcessId is %d ThreadId is %d \n", tbProcessID, toolbarThreadId)
	processHandler, _ := win32api.OpenProcessAllAccess(false, tbProcessID)
	fmt.Printf("ProcessHandler is %d \n", processHandler)
	n := make([]byte, 256)
	infoBaseAddress, _ := win32api.VirtualAllocEx(processHandler, 0, 256, MEM_COMMIT, PAGE_READWRITE)
	fmt.Printf("Base adrress is %d \n", infoBaseAddress)
	p := &n[0]
	length, _ := win32api.SendMessage(
		toolbarHandler,
		win32api.TB_GETBUTTONTEXT,
		uintptr(buttonIndex),
		infoBaseAddress)

	if length > 0 {
		index, _ := win32api.SendMessage(
			toolbarHandler,
			win32api.TB_COMMANDTOINDEX,
			uintptr(buttonIndex),
			0)
		var numRead uintptr
		if dataRead, _ := win32api.ReadProcessMemory(processHandler, infoBaseAddress,
			uintptr(unsafe.Pointer(p)),
			length*2,
			&numRead); !dataRead {
			fmt.Print("Nothing read \n")
		} else {
			fmt.Printf("Button with index %d is %s\n", index, string(n[:numRead]))

		}
	} else {
		fmt.Printf("Error requesting Buttontext %v\n", err)
	}

	return
}

// func GetNotifyToolbarHandler() (win.HWND, error) {
// 	if handler, err := GetNotifyIconOverflowWindowHandler(); err != nil {
// 		return win.HWND(0), err
// 	} else {
// 		if toolbarHandler := win.GetDlgItem(handler, NIOW_TOOLBAR32_ID); toolbarHandler > 0 {
// 			return toolbarHandler, nil
// 		}
// 	}
// 	return win.HWND(0), fmt.Errorf("Error getting NotifyToolbarHandler")
// }

// Kernel32.VirtualFreeEx(
// 	hProcess,
// 	ipRemoteBuffer,
// 	UIntPtr.Zero,
// 	MemAllocationType.RELEASE );

// Kernel32.CloseHandle( hProcess );
