package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type LogFormat struct{}

func (f LogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer != nil {
		buf = entry.Buffer
	} else {
		buf = &bytes.Buffer{}
	}

	buf.WriteByte('[')
	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString("] [")
	buf.WriteString(strings.ToUpper(entry.Level.String()))
	buf.WriteString("]: ")
	buf.WriteString(entry.Message)
	buf.WriteString(" \n")

	ret := append([]byte(nil), buf.Bytes()...) // copy buffer
	return ret, nil
}

var url = "http://localhost:8003/chat"
var logger = logrus.New()

func login(nick string, tickt string) {
	resp, err := http.Post(url+"/user/login", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", nick, tickt)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	logger.Info(string(body))
}

func renew(nick string, tickt string) {
	for {
		time.Sleep(time.Microsecond * 100)
		resp, err := http.Post(url+"/user/renew", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", nick, tickt)))
		resp.Body.Close()
		if err != nil {
			logger.Warn(err)
		}
	}
}

func create(usrNick string, usrTickt string, name string, tickt string) {
	data := fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s", usrNick, usrTickt, name)
	if tickt != "" {
		data += fmt.Sprintf("&ticket=%s", tickt)
	}
	resp, err := http.Post(url+"/channel/create", "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	logger.Info(string(body))
}

func join(usrNick string, usrTickt string, name string, tickt string) {
	data := fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s", usrNick, usrTickt, name)
	if tickt != "" {
		data += fmt.Sprintf("&ticket=%s", tickt)
	}
	resp, err := http.Post(url+"/channel/join", "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	logger.Info(string(body))
}

func get(nick string, tickt string) {
	for {
		time.Sleep(time.Microsecond * 100)
		resp, err := http.Post(url+"/msg/get", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", nick, tickt)))
		if err != nil {
			logger.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Warn(err)
		}
		message := gjson.Get(string(body), "returnObj").Array()
		for _, msg := range message {
			logger.Info(msg.String())
		}
	}
}

func send(usrNick string, usrTickt string, name string, tickt string, msg string) {
	resp, err := http.Post(url+"/msg/send", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s&ticket=%s&msg=%s", usrNick, usrTickt, name, tickt, msg)))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
}

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&LogFormat{})

	login("Jigsaw", "123456")
	go renew("Jigsaw", "123456")
	go get("Jigsaw", "123456")

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		} else if text == "/exit" {
			break
		} else if text == "/create" {
			create("Jigsaw", "123456", "game", "chance")
		} else if text == "/join" {
			join("Jigsaw", "123456", "PublicChannel", "")
		} else {
			send("Jigsaw", "123456", "PublicChannel", "", text)
		}
	}
}
