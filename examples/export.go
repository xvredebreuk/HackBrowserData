package main

import (
	"fmt"

	hkb "github.com/moond4rk/hackbrowserdata"
)

func main() {
	browser, err := hkb.NewBrowser(hkb.Firefox)
	if err != nil {
		panic(err)
	}
	passwords, err := browser.Passwords()
	if err != nil {
		panic(err)
	}
	fmt.Println(len(passwords))
	// all, err := browser.AllBrowsingData()
	// if err != nil {
	// 	panic(err)
	// }
}
