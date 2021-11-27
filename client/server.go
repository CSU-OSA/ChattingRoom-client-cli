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

func (s *Server) getChannelList() []*map[string]string {
	rslt := []*map[string]string{}

	resp, err := requests.Get(s.url+"/channel/list", nil)
	if err != nil {
		rslt = append(rslt, &map[string]string{"level": "Error", "message": err.Error()})
	}

	for _, channel := range gjson.Get(resp, "@this").Array() {
		rslt = append(rslt, &map[string]string{
			"level":   "Info",
			"message": fmt.Sprintf("Channel: %s", channel)})
	}

	return rslt
}
