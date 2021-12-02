package main

import (
	"bufio"
	"bytes"
	pojo "chattingroom-cli/proto"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
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

type User struct {
	name string
	conn net.Conn
}

func (u *User) logout() {
	data := &pojo.RequestPOJO{
		Operation: pojo.RequestPOJO_LOGOUT,
	}
	newData, _ := proto.Marshal(data)
	_, err := u.conn.Write(newData)
	if err != nil {
		logger.Error(err)
	}
}

func (u *User) joinChannel(channel, nick, ticket string) {
	data := &pojo.RequestPOJO{
		Operation: pojo.RequestPOJO_JOIN_CHA,
		Channel: &pojo.Channel{
			Channel: channel,
			Ticket:  &ticket,
			Nick:    &nick,
		},
	}
	newData, _ := proto.Marshal(data)
	_, err := u.conn.Write(newData)
	if err != nil {
		logger.Error(err)
	}
}

func (u *User) quitChannel(channel string) {
	data := &pojo.RequestPOJO{
		Operation: pojo.RequestPOJO_QUIT_CHA,
		Channel: &pojo.Channel{
			Channel: channel,
		},
	}
	newData, _ := proto.Marshal(data)
	_, err := u.conn.Write(newData)
	if err != nil {
		logger.Error(err)
	}
}

func (u *User) sendMsg(channel, content string) {
	data := &pojo.RequestPOJO{
		Operation: pojo.RequestPOJO_SENDMSG,
		Message: &pojo.RequestMessage{
			Channel: channel,
			Content: content,
		},
	}
	newData, _ := proto.Marshal(data)
	_, err := u.conn.Write(newData)
	if err != nil {
		logger.Error(err)
	}
}

func (u *User) getMsg() {
	for {
		time.Sleep(time.Microsecond * 100)

		data := &pojo.RequestPOJO{
			Operation: pojo.RequestPOJO_GETMSG,
		}
		newData, _ := proto.Marshal(data)
		_, err := u.conn.Write(newData)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (u *User) heartbeat() {
	for {
		time.Sleep(time.Second * 1)
		data := &pojo.RequestPOJO{
			Operation: pojo.RequestPOJO_HEARTBEAT,
		}
		newData, _ := proto.Marshal(data)
		_, err := u.conn.Write(newData)
		if err != nil {
			logger.Error(err)
		}
	}
}

type Client struct {
	server   string
	user     []*User
	currUser *User
}

func (c *Client) loginUser(name string) {
	if c.server == "" {
		logger.Error("Please set server first!")
		return
	}

	conn, err := net.Dial("tcp", c.server)
	if err != nil {
		logger.Error(err)
		return
	}
	user := User{name, conn}
	go user.heartbeat()
	go user.getMsg()

	c.user = append(c.user, &user)

	if c.currUser == nil {
		c.currUser = &user
	}
}

func (c *Client) logoutUser(name string) {
	for i, user := range c.user {
		if user.name == name {
			user.logout()
			c.user = append(c.user[:i], c.user[i+1:]...)
			break
		}
	}
	if c.currUser.name == name {
		c.currUser = nil
		logger.Warn("Please switch user manually!")
	}
}

func (c *Client) logoutAllUser() {
	for _, user := range c.user {
		user.logout()
	}
	c.user = []*User{}
	c.currUser = nil
}

func (c *Client) switchUser(name string) {
	for _, user := range c.user {
		if user.name == name {
			c.currUser = user
			logger.Info(fmt.Sprintf("Current user has switch to %s", name))
			break
		}
	}
}

func (c *Client) getMessage() {
	for {
		if c.currUser == nil {
			continue
		}
		var buf [4096]byte
		n, err := c.currUser.conn.Read(buf[:])
		if err != nil {
			logger.Error(err)
			continue
		}

		resp := &pojo.ResponsePOJO{}
		err = proto.Unmarshal(buf[:n], resp)
		if err != nil {
			logger.Error(err)
			continue
		}
		if resp.GetType() == pojo.ResponsePOJO_MESSAGE {
			for _, message := range resp.GetMessage() {
				logger.Info(message.GetChannel() +
					"｜" +
					message.GetFromNick() +
					"｜" +
					message.GetContent())
			}
		}
	}
}

var logger = logrus.New()

func main() {
	var client = Client{}
	go client.getMessage()

	parseCommand := func(text string) {
		params := append(strings.Split(text, " "), "")
		if text == "" {
			return
		} else if strings.HasPrefix(text, "/server") {
			client.server = params[1]
			logger.Info(fmt.Sprintf("Server has changed to %s", client.server))
		} else if strings.HasPrefix(text, "/user") {
			switch params[1] {
			case "login":
				client.loginUser(params[2])
			case "logout":
				client.logoutUser(params[2])
			case "switch":
				client.switchUser(params[2])
			}
		} else if strings.HasPrefix(text, "/channel") {
			switch params[1] {
			case "join":
				client.currUser.joinChannel(params[2], params[3], params[4])
			case "quit":
				client.currUser.quitChannel(params[2])
			}
		} else {
			client.currUser.sendMsg("PublicChannel", text)
		}
	}

	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		client.logoutAllUser()
		os.Exit(1)
	}()

	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&LogFormat{})

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
		if strings.HasPrefix(text, "/exit") {
			client.logoutAllUser()
			break
		}
		parseCommand(text)
	}
}
