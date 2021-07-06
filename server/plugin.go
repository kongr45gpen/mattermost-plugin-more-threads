package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	root *mux.Router
}

func (p *Plugin) OnActivate() error {
	p.SetupRouting()
	return nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.root.ServeHTTP(w, r)
}

func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	if post.RootId != "" {
		// Message posted in existing thread

		channel, _ := p.API.GetChannel(post.ChannelId)
		teamId := channel.TeamId

		page := 0
		for {
			members, err := p.API.GetChannelMembers(post.ChannelId, page, 30)
			if len(*members) == 0 || err != nil {
				break
			}

			for _, member := range *members {
				if member.UserId == post.UserId {
					// The thread participant should already be following this thread
					continue
				}

				followBytes, _ := p.API.KVGet(GetKVkey(member.UserId, member.ChannelId))

				autoFollowing := false
				if len(followBytes) >= 1 && followBytes[0] == 1 {
					autoFollowing = true
				}

				if autoFollowing {
					p.API.FollowThread(member.UserId, teamId, post.RootId)
				}
			}

			if len(*members) < 30 {
				break
			}

			page += 1
		}
	}
}

func GetKVkey(userId string, channelId string) string {
	return "a" + userId[:24] + channelId[:24]
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
