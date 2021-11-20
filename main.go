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

var logger = logrus.New()

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
	resp, err := http.Post(c.server+"/user/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", nick, ticket)))
	if err != nil {
		logger.Warn(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Warn(err)
	}
	if gjson.Get(string(body), "success").Bool() {
		logger.Info(fmt.Sprintf("User %s login success!", nick))
		c.user = append(c.user, User{nick, ticket})
		if c.currUser.nick == "" {
			c.currUser = c.user[0]
		}
	} else {
		logger.Error(fmt.Sprintf("User %s login fail: %s!",
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
		logger.Warn("Please switch currUser!")
	}
}

func (c *Client) renewUser(internal time.Duration) {
	for {
		if c.server == "" {
			continue
		}
		time.Sleep(internal)
		for _, usr := range c.user {
			resp, err := http.Post(c.server+"/user/renew",
				"application/x-www-form-urlencoded",
				strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", usr.nick, usr.ticket)))
			resp.Body.Close()
			if err != nil {
				logger.Warn(err)
			}
		}
	}
}

func (c *Client) switchUser(nick string) {
	for _, user := range c.user {
		if user.nick == nick {
			c.currUser = user
		}
	}
	logger.Info(fmt.Sprintf("Current user has switch to %s", nick))
}

func (c *Client) createChannel(name string, ticket string) {
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
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Info(err)
	}
	logger.Info(string(body))
}

func (c *Client) joinChannel(name string, ticket string) {
	logger.Debug(c.currUser)
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
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(string(body))
}

func (c *Client) getMsg() {
	for {
		if c.server == "" {
			continue
		}

		time.Sleep(time.Microsecond * 100)
		resp, err := http.Post(c.server+"/msg/get",
			"application/x-www-form-urlencoded",
			strings.NewReader(fmt.Sprintf("nick=%s&ticket=%s", c.currUser.nick, c.currUser.ticket)))
		if err != nil {
			logger.Error(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			logger.Error(err)
		}
		message := gjson.Get(string(body), "returnObj").Array()
		for _, msg := range message {
			logger.Info(fmt.Sprintf("%s | %s | %s",
				msg.Map()["channelName"],
				msg.Map()["senderNick"],
				msg.Map()["msg"].Array()[0]))
		}
	}
}

func (c *Client) sendMsg(name string, ticket string, msg string) {
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
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
	}
}

func main() {
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&LogFormat{})
	client := Client{}
	go client.renewUser(time.Microsecond * 100)
	go client.getMsg()

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		params := append(strings.Split(text, " "), "")
		if text == "" {
			continue
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
}
