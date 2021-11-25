package client

import (
	"fmt"

	"github.com/tidwall/gjson"

	"chattingroom/requests"
)

type User struct {
	nick   string
	ticket string
	server []*Server
}

func (u *User) createChannel(name string, ticket string) ([]string, error) {
	if len(u.server) == 0 {
		return []string{"Please bind server first!"}, nil
	}

	rslt := make([]string, 0)
	for _, server := range u.server {
		data := map[string]string{
			"usrNick":   u.nick,
			"usrTicket": u.ticket,
			"name":      name,
		}
		if ticket != "" {
			data["ticket"] = ticket
		}

		resp, err := requests.Post(server.url+"/channel/create", data)
		if err != nil {
			return nil, err
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, fmt.Sprintf("Create channel %s success!", name))
		} else {
			rslt = append(rslt, fmt.Sprintf("Create channel %s fail: %s",
				name,
				gjson.Get(resp, "msg")))
		}
	}
	return rslt, nil
}

func (u *User) joinChannel(name string, ticket string) ([]string, error) {
	if len(u.server) == 0 {
		return []string{"Please bind server first!"}, nil
	}
	rslt := make([]string, 0)
	for _, server := range u.server {
		data := map[string]string{
			"usrNick":   u.nick,
			"usrTicket": u.ticket,
			"name":      name,
		}
		if ticket != "" {
			data["ticket"] = ticket
		}

		resp, err := requests.Post(server.url+"/channel/join", data)
		if err != nil {
			return nil, err
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, fmt.Sprintf("Join channel %s success!", name))
		} else {
			rslt = append(rslt, fmt.Sprintf("Join channel %s fail: %s",
				name,
				gjson.Get(resp, "msg")))
		}
	}
	return rslt, nil
}

func (u *User) quitChannel(name string) ([]string, error) {
	if len(u.server) == 0 {
		return []string{"Please bind server first!"}, nil
	}
	rslt := make([]string, 0)
	for _, server := range u.server {
		data := map[string]string{
			"usrNick":   u.nick,
			"usrTicket": u.ticket,
			"name":      name,
		}

		resp, err := requests.Post(server.url+"/channel/quit", data)
		if err != nil {
			return nil, err
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, fmt.Sprintf("Join channel %s success!", name))
		} else {
			rslt = append(rslt, fmt.Sprintf("Join channel %s fail: %s",
				name,
				gjson.Get(resp, "msg")))
		}
	}
	return rslt, nil
}

func (u *User) getMsg() ([]string, error) {
	if len(u.server) == 0 {
		return []string{"Please bind server first!"}, nil
	}
	rslt := make([]string, 0)
	for _, server := range u.server {
		data := map[string]string{
			"nick":   u.nick,
			"ticket": u.ticket,
		}

		resp, err := requests.Post(server.url+"/msg/get", data)
		if err != nil {
			return []string{}, err
		}

		message := gjson.Get(resp, "returnObj").Array()
		for _, msg := range message {
			rslt = append(rslt, fmt.Sprintf("%s | %s | %s",
				msg.Map()["channelName"],
				msg.Map()["senderNick"],
				msg.Map()["msg"].Array()[0]))
		}
	}
	return rslt, nil
}

func (u *User) sendMsg(name string, msg string) ([]string, error) {
	if len(u.server) == 0 {
		return []string{"Please bind server first!"}, nil
	}

	rslt := make([]string, 0)
	for _, server := range u.server {

		data := map[string]string{
			"usrNick":   u.nick,
			"usrTicket": u.ticket,
			"name":      name,
			"msg":       msg,
		}

		_, err := requests.Post(server.url+"/msg/send", data)
		if err != nil {
			return nil, err
		}
	}
	return rslt, nil
}
