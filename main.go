package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// 格式化 log
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

type Client struct {
	server   string
	user     []User
	currUser User
}

type User struct {
	nick   string
	ticket string
}

func (c *Client) loginUser(nick string, ticket string) {
	if c.server == "" {
		logger.Error("Please set server first!")
		return
	}

	resp, err := http.Post(c.server+"/user/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", nick, ticket)))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	if gjson.Get(string(body), "success").Bool() {
		logger.Info(fmt.Sprintf("User %s login success!", nick))
		c.user = append(c.user, User{nick, ticket})
		if c.currUser.nick == "" {
			c.currUser = c.user[0]
		}
	} else {
		logger.Error(fmt.Sprintf("User %s login fail: %s",
			nick,
			gjson.Get(string(body), "msg").String()))
	}
}

func (c *Client) logoutUser(nick string) {
	for i, user := range c.user {
		if user.nick == nick {
			c.user = append(c.user[:i], c.user[i+1:]...)
			logger.Info(fmt.Sprintf("User %s Has logout!", nick))
			break
		}
	}
	if c.currUser.nick == nick {
		c.currUser = User{}
		logger.Warn("Please switch user manually!")
	}
}

func (c *Client) renewUser(internal time.Duration) {
	for {
		time.Sleep(internal)
		for _, usr := range c.user {
			resp, err := http.Post(c.server+"/user/renew",
				"application/x-www-form-urlencoded",
				strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", usr.nick, usr.ticket)))
			resp.Body.Close()
			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func (c *Client) switchUser(nick string) {
	for _, user := range c.user {
		if user.nick == nick {
			c.currUser = user
			logger.Info(fmt.Sprintf("Current user has switch to %s", nick))
			break
		}
	}
}

func (c *Client) createChannel(name string, ticket string) {
	if c.server == "" || c.currUser.nick == "" {
		logger.Error("Please set server and login first!")
		return
	}

	data := fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s",
		c.currUser.nick,
		c.currUser.ticket,
		name)
	if ticket != "" {
		data += fmt.Sprintf("&ticket=%s", ticket)
	}
	resp, err := http.Post(c.server+"/channel/create",
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	if gjson.Get(string(body), "success").Bool() {
		logger.Info(fmt.Sprintf("Create channel %s success!", name))
	} else {
		logger.Error(fmt.Sprintf("Create channel %s fail: %s",
			name,
			gjson.Get(string(body), "msg").String()))
	}
}

func (c *Client) joinChannel(name string, ticket string) {
	if c.server == "" || c.currUser.nick == "" {
		logger.Error("Please set server and login first!")
		return
	}

	data := fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s",
		c.currUser.nick,
		c.currUser.ticket,
		name)
	if ticket != "" {
		data += fmt.Sprintf("&ticket=%s", ticket)
	}
	resp, err := http.Post(c.server+"/channel/join",
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}

	if gjson.Get(string(body), "success").Bool() {
		logger.Info(fmt.Sprintf("Join channel %s success!", name))
	} else {
		logger.Error(fmt.Sprintf("Join channel %s fail: %s",
			name,
			gjson.Get(string(body), "msg").String()))
	}
}

func (c *Client) getMsg() {
	for {
		time.Sleep(time.Microsecond * 100)

		if c.server == "" || c.currUser.nick == "" {
			continue
		}

		resp, err := http.Post(c.server+"/msg/get",
			"application/x-www-form-urlencoded",
			strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", c.currUser.nick, c.currUser.ticket)))
		if err != nil {
			logger.Error(err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Error(err)
			continue
		}
		message := gjson.Get(string(body), "returnObj").Array()
		for _, msg := range message {
			t, _ := time.ParseInLocation("2006-01-02 15:04:05", msg.Map()["recTime"].String(), time.Local)
			logger.WithTime(t).Info(fmt.Sprintf("%s | %s | %s",
				msg.Map()["channelName"],
				msg.Map()["senderNick"],
				msg.Map()["msg"].Array()[0]))
		}
	}
}

func (c *Client) sendMsg(name string, ticket string, msg string) {
	if c.server == "" || c.currUser.nick == "" {
		logger.Error("Please set server and login first!")
		return
	}

	resp, err := http.Post(c.server+"/msg/send",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("usrNick=%s&usrTicket=%s&name=%s&ticket=%s&msg=%s",
			c.currUser.nick,
			c.currUser.ticket,
			name,
			ticket,
			msg)))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
}

var logger = logrus.New()
var client = Client{}

func parseCommand(text string) {
	params := append(strings.Split(text, " "), "")
	if text == "" {
		return
	} else if strings.HasPrefix(text, "/server") {
		client.server = params[1]
		logger.Info(fmt.Sprintf("Server has changed to %s", client.server))
	} else if strings.HasPrefix(text, "/user") {
		switch params[1] {
		case "login":
			client.loginUser(params[2], params[3])
		case "logout":
			client.logoutUser(params[2])
		case "switch":
			client.switchUser(params[2])
		}
	} else if strings.HasPrefix(text, "/channel") {
		switch params[1] {
		case "create":
			client.createChannel(params[2], params[3])
		case "join":
			client.joinChannel(params[2], params[3])
		}
	} else {
		client.sendMsg("PublicChannel", "", text)
	}
}

func main() {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&LogFormat{})

	go client.renewUser(time.Microsecond * 100)
	go client.getMsg()

	f, err := os.Open("./.chattingroomrc")
	if err == nil {
		logger.Info("Load chattingroomrc:")
		reader := bufio.NewReader(f)
		for {
			text, err := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			parseCommand(text)
			if err == io.EOF {
				break
			}
		}
	}
	f.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		parseCommand(text)
	}
}
