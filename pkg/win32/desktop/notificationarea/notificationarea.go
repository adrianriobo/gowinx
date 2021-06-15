// +build windows
package notificationarea

import (
	"fmt"
	"syscall"

	win32wam "github.com/adrianriobo/gowinx/pkg/win32/api/windows-and-messages"
	win32toolbar "github.com/adrianriobo/gowinx/pkg/win32/ux/commands/toolbar"
	win32windows "github.com/adrianriobo/gowinx/pkg/win32/ux/windows"
)

// systemtray aka notification area, it is composed of notifications icons (offering display the status and various functions)
// distributerd across:
// * visible area on the right side of the taskbar (class: Shell_TrayWnd)
// * hidden area as overflowwindow ( class: NotifyIconOverflowWindow)

const (
	NOTIFICATION_AREA_VISIBLE_WINDOW_CLASS string = "Shell_TrayWnd"
	NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS  string = "NotifyIconOverflowWindow"
	TOOLBARWINDOWS32_ID                    int32  = 1504
)

func GetHiddenNotificationAreaRect() (rect win32wam.RECT, err error) {
	// Show notification area (hidden)
	if err = ShowHiddenNotificationArea(); err == nil {
		if toolbarHandler, err := getNotificationAreaToolbarByWindowClass(NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS); err == nil {
			if _, err = win32wam.GetWindowRect(toolbarHandler, &rect); err == nil {
				fmt.Printf("Rect for system tray t:%d,l:%d,r:%d,b:%d\n", rect.Top, rect.Left, rect.Right, rect.Bottom)
			}
		}
	}
	if err != nil {
		fmt.Printf("error getting hidden notification area rect: %v\n", err)
	}
	return
}

func ShowHiddenNotificationArea() (err error) {
	if handler, err := getNotificationAreaWindowByClass(NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS); err == nil {
		win32wam.ShowWindow(handler, win32wam.SW_SHOWNORMAL)
	}
	return
}

func getNotificationAreaWindowByClass(className string) (handler syscall.Handle, err error) {
	if handler, err = win32windows.FindWindowByClass(className); err != nil {
		fmt.Printf("error getting handler on notification area for windows class: %s, error: %v\n", className, err)
	}
	return
}

func getNotificationAreaToolbarByWindowClass(className string) (handler syscall.Handle, err error) {
	if windowHandler, err := getNotificationAreaWindowByClass(className); err == nil {
		if handler, err = win32wam.GetDlgItem(windowHandler, TOOLBARWINDOWS32_ID); err != nil {
			fmt.Printf("error getting toolbar handler on notification area for windows class: %s, error: %v\n", className, err)
		}
	}
	return
}

func GetIconPositionByTitle(buttonText string) (int, int, error) {
	toolbarHandlers, _ := findToolbars()
	for i, toolbarHandler := range toolbarHandlers {
		fmt.Printf("trying on toolbar %d\n", i)
		if x, y, err := win32toolbar.GetButtonClickablePosition(toolbarHandler, buttonText); err == nil {
			return x, y, nil
		}
	}
	return -1, -1, fmt.Errorf("button %s not found on toolbar\n", buttonText)

}

// The notification area is composed of elements, app notification icons use to be placed
// at the toolbars
func findToolbars() ([]syscall.Handle, error) {
	handler, _ := win32windows.FindWindowByClass(NOTIFICATION_AREA_VISIBLE_WINDOW_CLASS)
	toolbars, _ := win32toolbar.FindToolbars(handler)
	toolbarHandler, err := getNotificationAreaToolbarByWindowClass(NOTIFICATION_AREA_HIDDEN_WINDOW_CLASS)
	if err == nil {
		toolbars = append(toolbars, toolbarHandler)
	}
	return toolbars, nil
}
