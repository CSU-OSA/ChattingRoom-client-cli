package client

import (
	"fmt"

	"github.com/tidwall/gjson"

	"chattingroom/requests"
)

type Server struct {
	name string
	url  string
}

func (s *Server) getChannelList() ([]string, error) {
	resp, err := requests.Get(s.url+"/channel/list", nil)
	if err != nil {
		return nil, err
	}

	rslt := []string{}
	for _, channel := range gjson.Get(resp, "@this").Array() {
		rslt = append(rslt, (fmt.Sprintf("Channel: %s", channel)))
	}
	return rslt, nil
}

func (s *Server) getUserList() ([]string, error) {
	resp, err := requests.Get(s.url+"/user/list", nil)
	if err != nil {
		return nil, err
	}

	rslt := []string{}
	for _, channel := range gjson.Get(resp, "@this").Array() {
		rslt = append(rslt, (fmt.Sprintf("User: %s", channel)))
	}
	return rslt, nil
}
