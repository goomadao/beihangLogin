package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Acid     string `json:"acid"`
	Enc_ver  string `json:"enc_ver"`
}

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client = http.Client{
		Transport: tr,
	}
	baseUrl      = "https://gw.buaa.edu.cn/cgi-bin/"
	callbackName = "laji_srun"
)

func Login(username, password string) error {
	var (
		ac_id   = "1"
		enc_ver = "srun_bx1"
		n       = "200"
		Type    = "1"
		os      = "Mac OS"
		name    = "Macintosh"
		// os           = "Windows 10"
		// name         = "Windows"
		double_stack = "0"
	)
	challenge, err := getResponse("challenge", url.Values{
		"username": []string{username},
	})
	if err != nil {
		return err
	}
	token, exist := challenge["challenge"].(string)
	if !exist {
		return errors.New("There is no challenge field in challenge response.")
	}
	ip, exist := challenge["client_ip"].(string)
	if !exist {
		return errors.New("There is no client_ip field in challenge response.")
	}

	params, err := getLoginParams(username, password, ac_id, ip, enc_ver, n, Type, os, name, double_stack, token)
	if err != nil {
		return err
	}

	loginResponse, err := getResponse("login", params)
	if err != nil {
		return err
	}
	res, exist := loginResponse["error"].(string)
	if !exist {
		return errors.New("There is no error field in login response.")
	}
	if res == "ok" {
		suc_msg, exist := loginResponse["suc_msg"].(string)
		if exist && suc_msg == "login_ok" {
			return nil
		} else if exist {
			// return errors.New(suc_msg)
			fmt.Println(suc_msg)
		} else {
			fmt.Println("unknown error")
			// return errors.New("unknown error")
		}
		return nil
	}
	error_msg, exist := loginResponse["error_msg"].(string)
	if exist && len(error_msg) > 0 {
		res += ", " + error_msg
	}
	return errors.New(res)
}

func Logout(username string) error {
	conn, err := net.Dial("tcp", "10.200.21.4:443")
	if err != nil {
		return err
	}
	defer conn.Close()
	ip := conn.LocalAddr().String()
	resp, err := getResponse("logout", url.Values{
		"action":   []string{"logout"},
		"username": []string{username},
		"ac_id":    []string{"1"},
		"ip":       []string{ip[:strings.LastIndex(ip, ":")]},
	})
	if err != nil {
		return err
	}
	res, exist := resp["error"].(string)
	if !exist {
		return errors.New("There is no error field in logout response.")
	}
	if res == "ok" {
		return nil
	}
	error_msg, exist := resp["error_msg"].(string)
	if exist && len(error_msg) > 0 {
		res += ", " + error_msg
	}
	return errors.New(res)
}

func Status() (map[string]interface{}, error) {
	resp, err := getResponse("status", url.Values{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func getResponse(reqType string, params url.Values) (map[string]interface{}, error) {
	reqUrl := baseUrl
	if reqType == "challenge" {
		reqUrl += "get_challenge"
	} else if reqType == "login" || reqType == "logout" {
		reqUrl += "srun_portal"
	} else if reqType == "status" {
		reqUrl += "rad_user_info"
	}
	params.Set("callback", callbackName)
	params.Set("_", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	reqUrl += "?" + params.Encode()
	Logger.Debug(reqUrl)
	// if reqType == "logout" {
	jar, err := cookiejar.New(nil)
	cookieUrl, _ := url.Parse(baseUrl)
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{Name: "cookie", Value: "14511497"})
	jar.SetCookies(cookieUrl, cookies)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	// }
	resp, err := client.Get(reqUrl)
	if err != nil {
		// Logger.Error("Get "+reqType+" response failed", zap.Error(err))
		Logger.Error("Get " + reqType + " response failed")
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Logger.Error("Get "+reqType+" body failed", zap.Error(err))
		Logger.Error("Get " + reqType + " body failed")
		return nil, err
	}
	body = body[len(callbackName)+1 : len(body)-1]
	var out bytes.Buffer
	err = json.Indent(&out, body, "", "\t")
	if err != nil {
		Logger.Debug("Indent " + reqType + " response failed")
		Logger.Debug(reqType + " response: " + string(body))
		// return nil, err
	} else {
		Logger.Debug(reqType + " response: " + out.String())
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Error("Unmarshal " + reqType + " body failed")
		return nil, err
	}
	return response, nil
}

func getLoginParams(username, password, ac_id, ip, enc_ver, n, Type, os, name, double_stack, token string) (url.Values, error) {
	infoJson, err := json.Marshal(Info{
		Username: username,
		Password: password,
		Ip:       ip,
		Acid:     ac_id,
		Enc_ver:  enc_ver,
	})
	if err != nil {
		return nil, err
	}
	info := GetEncodedInfo(string(infoJson), token)

	password = GetEncodedPassword(password, token)

	chkstr := token + username + token + password + token + ac_id + token + ip + token + n + token + Type + token + info
	chksum := GetEncodedChkstr(chkstr)

	return url.Values{
		"action":   []string{"login"},
		"username": []string{username},
		"password": []string{"{MD5}" + password},
		"ac_id":    []string{ac_id},
		"ip":       []string{ip},
		"chksum":   []string{chksum},
		"info":     []string{info},
		"n":        []string{n},
		"type":     []string{Type},
		// "os":           []string{os},
		// "name":         []string{name},
		"double_stack": []string{double_stack},
	}, nil
}
