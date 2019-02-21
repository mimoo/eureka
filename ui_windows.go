package main

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

func install(dir string) {
	// obtain HKEY_CLASSES_ROOT
	hkey := registry.Key(syscall.HKEY_CLASSES_ROOT)

	// add right-click > encrypt
	paths := map[string]string{
		// right-click on a file
		`*\shell`: "%1",
		// right-click on a folder
		`Directory\shell`: "%1",
		// right-click in the background, when inside a folder
		`Directory\background\shell`: "%V",
	}
	for path, arg := range paths {
		newk, _, err := registry.CreateKey(hkey, path+`\Encrypt\command`, registry.ALL_ACCESS)
		if err != nil {
			fmt.Println("cannot create registry key for */shell/encrypt")
			newk.Close()
			return
		}
		err = newk.SetStringValue("", `"`+dir+`\eureka.exe" -encrypt -file "`+arg+`"`)
		if err != nil {
			fmt.Println("cannot create value (Default)")
			return
		}
		newk.Close()
	}

	// add right-click > decrypt
	newk, _, err := registry.CreateKey(hkey, `.encrypted\Shell\Decrypt\command`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("cannot create registry key .encrypted")
		newk.Close()
		return
	}
	err = newk.SetStringValue("", `"`+dir+`\eureka.exe" -decrypt -file "%1"`)
	if err != nil {
		fmt.Println("cannot create value (Default)")
		return
	}
	newk.Close()

	// add eureka icon for .encrypted
	newk, _, err = registry.CreateKey(hkey, `.encrypted\DefaultIcon`, registry.ALL_ACCESS)
	if err != nil {
		fmt.Println("cannot create registry key .encrypted")
		newk.Close()
		return
	}
	err = newk.SetStringValue("", `"`+dir+`\eureka.ico"`)
	if err != nil {
		fmt.Println("cannot create value (Default)")
		return
	}
	newk.Close()
}
