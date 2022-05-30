package main

import "net/http"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckStatus(resp *http.Response) {
	if resp.StatusCode != 200 {
		panic(resp.StatusCode)
	}
}
