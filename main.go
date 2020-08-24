package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	username := "admin"
	password := "admin"
	passwordMd5 := md5.Sum([]byte(password))
	passwordMd5Str := hex.EncodeToString(passwordMd5[:])
	loginResponse , loginCookie := makeRequest("checkLogin",fmt.Sprintf("login_username=%s&passwd=%s",username, passwordMd5Str),"")
	if loginResponse == nil {
		fmt.Println("Invalid login response")
		os.Exit(1)
	}else{
		if *loginResponse == "success" && *loginCookie != "" {
			fmt.Printf("Login Cookie : %s\n",*loginCookie)
			unlockResponse , _ := makeRequest("LteSetSimRestriction","SimCardCheck=0",*loginCookie)
			if *unlockResponse == "success" {
				fmt.Println("SUCCESSFULLY UNLOCKED")
				fmt.Println("Removing IMSI Prefix")
				imsi_prefix , _ := makeRequest("LteSetSimRestriction","imsi_prefix=",*loginCookie)
				if *imsi_prefix == "success" {
					fmt.Println("IMSI Prefix removed")
				}else{
					fmt.Println("unable to remove imsi prefix  :/")
				}
			}else{
				fmt.Println("Fail :/")
			}
		}else{
			fmt.Print("Invalid Credentials")
			os.Exit(1)
		}
	}
}

func makeRequest(path string,params string,cookies string)(response* string ,cookie* string){
	url := "http://192.168.1.1/goform/" + path
	method := "POST"

	payload := strings.NewReader(params)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil,nil
	}
	req.Header.Add("Proxy-Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("DNT", "1")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36 Edg/84.0.522.63")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "http://192.168.1.1")
	req.Header.Add("Referer", "http://192.168.1.1/lte/cmdshell.asp")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
	if cookies != "" {
		req.Header.Add("Cookie", cookies)
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	mCookie := res.Header.Get("Set-Cookie")
	mBody := string(body)

	cookie = &mCookie
	response = &mBody
	return
}
