package client

import (
	"fmt"

	"github.com/tidwall/gjson"

	"chattingroom/requests"
)

type User struct {
	nick   string
	ticket string
	server string
}

func (u *User) createChannel(name string, ticket string) (string, error) {
	if u.server == "" {
		return "Please bind server first!", nil
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}
	if ticket != "" {
		data["ticket"] = ticket
	}

	resp, err := requests.Post(u.server+"/channel/create", data)
	if err != nil {
		return "", err
	}

	if gjson.Get(resp, "success").Bool() {
		return fmt.Sprintf("Create channel %s success!", name), nil
	} else {
		return fmt.Sprintf("Create channel %s fail: %s",
			name,
			gjson.Get(resp, "msg")), nil
	}
}

func (u *User) joinChannel(name string, ticket string) (string, error) {
	if u.server == "" {
		return "Please bind server first!", nil
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}
	if ticket != "" {
		data["ticket"] = ticket
	}

	resp, err := requests.Post(u.server+"/channel/join", data)
	if err != nil {
		return "", err
	}

	if gjson.Get(resp, "success").Bool() {
		return fmt.Sprintf("Join channel %s success!", name), nil
	} else {
		return fmt.Sprintf("Join channel %s fail: %s",
			name,
			gjson.Get(resp, "msg")), nil
	}
}

func (u *User) quitChannel(name string) (string, error) {
	if u.server == "" {
		return "Please bind server first!", nil
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}

	resp, err := requests.Post(u.server+"/channel/quit", data)
	if err != nil {
		return "", err
	}

	if gjson.Get(resp, "success").Bool() {
		return fmt.Sprintf("Join channel %s success!", name), nil
	} else {
		return fmt.Sprintf("Join channel %s fail: %s",
			name,
			gjson.Get(resp, "msg")), nil
	}
}

func (u *User) getMsg() ([]string, error) {
	if u.server == "" {
		return []string{}, nil
	}

	data := map[string]string{
		"nick":   u.nick,
		"ticket": u.ticket,
	}

	resp, err := requests.Post(u.server+"/msg/get", data)
	if err != nil {
		return []string{}, err
	}

	message := gjson.Get(resp, "returnObj").Array()
	result := []string{}
	for _, msg := range message {
		result = append(result, fmt.Sprintf("%s | %s | %s",
			msg.Map()["channelName"],
			msg.Map()["senderNick"],
			msg.Map()["msg"].Array()[0]))
	}
	return result, nil
}

func (u *User) sendMsg(name string, msg string) (string, error) {
	if u.server == "" {
		return "Please bind server first!", nil
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
		"msg":       msg,
	}

	_, err := requests.Post(u.server+"/msg/send", data)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return "", nil
}
