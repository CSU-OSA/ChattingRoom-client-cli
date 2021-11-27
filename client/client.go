package client

import (
	"time"
)

type Client struct {
	server   []*Server
	user     []*User
	currUser *User
	Logger   chan *map[string]string
}

func (c *Client) LoginUser(name string, ticket string, server []string) {
	serverResult := []*Server{}
	for _, srv := range server {
		isServerName := false
		for _, s := range c.server {
			if s.name == srv {
				serverResult = append(serverResult, s)
			}
		}
		if !isServerName {
			newServer := Server{url: srv}
			c.server = append(c.server, &newServer)
			serverResult = append(serverResult, &newServer)
		}
	}
	user := &User{name, ticket, serverResult}

	rslt := user.login()
	for _, r := range rslt {
		c.Logger <- r
	}

	c.user = append(c.user, user)
	if c.currUser == nil {
		c.currUser = user
	}
}

func (c *Client) RenewUser(duration time.Duration) {
	for {
		time.Sleep(duration)
		for _, user := range c.user {
			rslt := user.renew()
			for _, r := range rslt {
				c.Logger <- r
			}
		}
	}
}

func (c *Client) LogoutUser(nick string) {
	for i, user := range c.user {
		if user.nick == nick {
			rslt := user.logout()
			for _, r := range rslt {
				c.Logger <- r
			}

			if c.currUser == user {
				c.user = append(c.user[:i], c.user[i+1:]...)
				c.currUser = nil
			}

			break
		}
	}
}

func (c *Client) LogoutAllUser() {
	for _, user := range c.user {
		rslt := user.logout()
		for _, r := range rslt {
			c.Logger <- r
		}
	}
	c.user = []*User{}
	c.currUser = nil
}

func (c *Client) SwitchUser(nick string) {
	for _, user := range c.user {
		if user.nick == nick {
			c.currUser = user
			break
		}
	}
}

func (c *Client) CreateChannel(name string, ticket string) {
	rslt := c.currUser.createChannel(name, ticket)
	for _, r := range rslt {
		c.Logger <- r
	}
}

func (c *Client) JoinChannel(name string, ticket string) {
	rslt := c.currUser.joinChannel(name, ticket)
	for _, r := range rslt {
		c.Logger <- r
	}
}

func (c *Client) QuitChannel(name string) {
	rslt := c.currUser.quitChannel(name)
	for _, r := range rslt {
		c.Logger <- r
	}
}

func (c *Client) GetMsg() {
	for {
		time.Sleep(time.Microsecond * 100)
		for _, user := range c.user {
			rslt := user.getMsg()
			for _, r := range rslt {
				c.Logger <- r
			}
		}
	}
}

func (c *Client) SendMsg(channel string, msg string) {
	rslt := c.currUser.sendMsg(channel, msg)
	for _, r := range rslt {
		c.Logger <- r
	}
}
