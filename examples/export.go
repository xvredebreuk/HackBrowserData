package main

import (
	"fmt"

	hkb "github.com/moond4rk/hackbrowserdata"
)

func main() {
	chrome, err := hkb.NewBrowser(hkb.Chrome)
	if err != nil {
		panic(err)
	}
	passwords, err := chrome.Passwords()
	if err != nil {
		panic(err)
	}
	for _, pass := range passwords {
		fmt.Printf("%+v\n", pass)
	}
	// cookies, err := chrome.Cookies()
	// if err != nil {
	// 	panic(err)
	// }

	// creditCards, err := browser.CreditCards()
	// bookmarks, err := browser.Bookmarks()
	// downloads, err := browser.Downloads()
	// extensions, err := browser.Extensions()
	// history, err := browser.History()
	// localStorage, err := browser.LocalStorage()
	// sessionStorage, err := browser.SessionStorage()

	// items := []hkb.browsingDataType{hkb.Cookie, hkb.Password, hkb.History, hkb.Bookmark}
	// browsingData, err := browser.BrowsingDatas(items)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// if err != nil {
	// 	panic(err)
	// }
	// _ = cookies
	// items := []hkb.browsingDataType{hkb.Cookie}
	// datas, err := browser.BrowsingDatas(items)
	// if err != nil {
	// 	panic(err)
	// }
	// _ = datas
}
