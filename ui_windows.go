package main

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

// TODO: clean this?
// TODO: will what I wrote work with win10?
const (
	// right-click on a file
	win7file = `*\shellex\ContextMenuHandlers`
	// right-click on a folder
	win7Folder = `Directory\shell`
	// right-click in the background, when inside a folder
	win7FolderBackground = `Directory\background`
	other                = `Directory\Background\shell\encrypt`
	// right-click
	win10ContextMenu = `*\shellex\ContextMenuHandlers`
	other2           = `*\shell\encrypt`
)

// TODO: when to call that?
// TODO: and what is the real path of the app for the commands?
func install() {
	// obtain HKEY_CLASSES_ROOT
	k := registry.Key(syscall.HKEY_CLASSES_ROOT)

	// add right-click > encrypt
	newk, _, err := registry.CreateKey(k, `*\shell\Encrypt\command`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("cannot create registry key for */shell/encrypt")
		newk.Close()
		return
	}
	err = newk.SetStringValue("", `"C:\Users\David\Downloads\eureka-master\eureka.exe" -encrypt -file "%1"`)
	if err != nil {
		fmt.Println("cannot create value (Default)")
		return
	}
	newk.Close()

	// add right-click > decrypt
	newk, _, err = registry.CreateKey(k, `.encrypted\Shell\Decrypt\command`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("cannot create registry key .encrypted")
		newk.Close()
		return
	}
	err = newk.SetStringValue("", `"C:\Users\David\Downloads\eureka-master\eureka.exe" -decrypt -file "%1"`)
	if err != nil {
		fmt.Println("cannot create value (Default)")
		return
	}
	newk.Close()

	// add eureka icon for .encrypted
	newk, _, err = registry.CreateKey(k, `.encrypted\DefaultIcon`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("cannot create registry key .encrypted")
		newk.Close()
		return
	}
	err = newk.SetStringValue("", `"C:\Users\David\Downloads\eureka-master\eureka.ico"`)
	if err != nil {
		fmt.Println("cannot create value (Default)")
		return
	}
	newk.Close()
}
