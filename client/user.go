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

func (u *User) login() []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"nick":   u.nick,
		"ticket": u.ticket,
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/user/login", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Info",
				"message": fmt.Sprintf("User %s login success!", u.nick)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf(
					"User %s login fail: %s",
					u.nick,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) renew() []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"nick":   u.nick,
		"ticket": u.ticket,
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/user/renew", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Debug",
				"message": fmt.Sprintf("User %s renew success!", u.nick)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf(
					"User %s renew fail: %s",
					u.nick,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) logout() []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"nick":   u.nick,
		"ticket": u.ticket,
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/user/logout", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Info",
				"message": fmt.Sprintf("User %s logout success!", u.nick)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf(
					"User %s logout fail: %s",
					u.nick,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) createChannel(name string, ticket string) []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}
	if ticket != "" {
		data["ticket"] = ticket
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/channel/create", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Info",
				"message": fmt.Sprintf("Create channel %s success!", name)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf(
					"Create channel %s fail: %s",
					name,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) joinChannel(name string, ticket string) []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}
	if ticket != "" {
		data["ticket"] = ticket
	}

	for _, server := range u.server {

		resp, err := requests.Post(server.url+"/channel/join", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Info",
				"message": fmt.Sprintf("Join channel %s success!", name)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf("Join channel %s fail: %s",
					name,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) quitChannel(name string) []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/channel/quit", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
			continue
		}

		if gjson.Get(resp, "success").Bool() {
			rslt = append(rslt, &map[string]string{
				"level":   "Info",
				"message": fmt.Sprintf("Join channel %s success!", name)})
		} else {
			rslt = append(rslt, &map[string]string{
				"level": "Error",
				"message": fmt.Sprintf(
					"Join channel %s fail: %s",
					name,
					gjson.Get(resp, "msg"))})
		}
	}

	return rslt
}

func (u *User) getMsg() []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"nick":   u.nick,
		"ticket": u.ticket,
	}

	for _, server := range u.server {
		resp, err := requests.Post(server.url+"/msg/get", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{
				"level":   "Error",
				"message": err.Error()})
		}

		message := gjson.Get(resp, "returnObj").Array()
		for _, msg := range message {
			m := fmt.Sprintf("%s | %s | %s",
				msg.Map()["channelName"],
				msg.Map()["senderNick"],
				msg.Map()["msg"].Array()[0])
			isNew := true
			for _, r := range rslt {
				if (*r)["message"] == m {
					isNew = false
					break
				}
			}
			if isNew {
				rslt = append(rslt, &map[string]string{
					"level":   "Info",
					"message": m})
			}
		}
	}

	return rslt
}

func (u *User) sendMsg(name string, msg string) []*map[string]string {
	rslt := []*map[string]string{}

	if len(u.server) == 0 {
		rslt = append(rslt, &map[string]string{
			"level":   "Error",
			"message": "Please bind server first!"})
	}

	data := map[string]string{
		"usrNick":   u.nick,
		"usrTicket": u.ticket,
		"name":      name,
		"msg":       msg,
	}

	for _, server := range u.server {
		_, err := requests.Post(server.url+"/msg/send", data)
		if err != nil {
			rslt = append(rslt, &map[string]string{
				"level":   "Error",
				"message": err.Error()})
			continue
		}
	}

	return rslt
}
