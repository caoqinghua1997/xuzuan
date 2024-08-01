package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

// vars
var Client *http.Client

type myPayLoad map[string]int

var userName string
var key string
var orderProductID int
var freq int

func main() {
	for {
		time.Sleep(time.Second * 6)
		getOrder()
	}
}

func getOrder() {
	var mySlice []myPayLoad
	mySlice = append(mySlice, myPayLoad{"formatId": orderProductID, "price": 0})
	tamp := time.Now().UnixNano() / 1000000
	data := map[string]any{
		"formatIds": mySlice,
		"count":     1,
	}
	fmt.Println(data)
	jsonData, _ := json.Marshal(data)
	sign := genSign(jsonData)
	url := fmt.Sprintf("http://open.xuzuan.cn/api/MarketOrder/ReceivingOrder?timestamp=%d&userName=%s&sign=%s", tamp, userName, sign)
	req := Server(url, jsonData)
	fmt.Println(req)
	if strings.Contains(req, "成功") {
		fmt.Println("接到订单啦， 停止运行")
		// QQ通知
		// Notify()
		os.Exit(0)
	}
}

func Notify() {
	//  QQ通知
	// curl  http://ip:port/send_private_msg?user_id=QQ&message="接到订单啦!请处理！"
}

func Server(API string, payload []byte) string {
	var data = bytes.NewReader(payload)
	fmt.Println(API)
	req, err := http.NewRequest("POST", API, data)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func init() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error conf file: %s \n", err))
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	Client = &http.Client{
		Jar: jar,
	}

	userName = viper.GetString("username")
	key = viper.GetString("key")
	orderProductID = viper.GetInt("productID")
	freq = viper.GetInt("freq")
	if freq < 3 {
		freq = 5
	}
}

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

func genSign(s []byte) string {
	var tmp = string(s) + key
	sign := MD5(tmp)
	return sign
}
