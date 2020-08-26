package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

var ip = "192.168.1.1"
func main() {
	if ipAddress := readUserLine("Enter IP Address : [192.168.1.1]"); ipAddress != nil && isIpv4Valid(*ipAddress) {
		ip = *ipAddress
	}
	username := "admin"
	if u := readUserLine("Enter Username : [admin]"); u != nil && *u != "" {
		username = *u
	}
	password := "admin"
	if p := readUserLine("Enter Password : [admin]"); p != nil && *p != "" {
		password = *p
	}
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
func readUserLine(prompt string) *string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(prompt)
	line, err := reader.ReadString('\n')
	line = strings.ReplaceAll(line , "\n","")
	line = strings.ReplaceAll(line , "\r","")
	if err != nil{
		return nil
	}
	return &line
}
func isIpv4Valid(host string) bool {
	return net.ParseIP(host) != nil
}
func makeRequest(path string,params string,cookies string)(response* string ,cookie* string){
	url := fmt.Sprintf("http://%s/goform/%s",ip,path)
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
	req.Header.Add("Origin", fmt.Sprintf("http://%s",ip))
	req.Header.Add("Referer", fmt.Sprintf("http://%s/lte/cmdshell.asp",ip))
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")
	if cookies != "" {
		cookies = strings.ReplaceAll(cookies,"path=/","")
		req.Header.Add("Cookie", cookies + " kz_userid=Administrator:1;")
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
