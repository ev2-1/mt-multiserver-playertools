package playerTools

import (
	proxy "github.com/HimbeerserverDE/mt-multiserver-proxy"

	"sync"
)

type PlayerList map[string]*proxy.ClientConn

func (pl PlayerList) Array() (a []string) {
	for player := range pl {
		a = append(a, player)
	}
	return
}

type PlayerListUpdateHandler struct {
	Join   func(string)
	Leave  func(string)
	Update func(PlayerList)
}

var playerListUpdateHandlers []*PlayerListUpdateHandler
var playerListUpdateHandlersMu sync.RWMutex

var playerListMu sync.RWMutex
var playerList = make(PlayerList)

func Players() PlayerList {
	playerListMu.Lock()
	defer playerListMu.Unlock()

	return playerList
}

func RegisterPlayerListUpdateHandler(h *PlayerListUpdateHandler) {
	initProxyPlayerlistHandlers()

	playerListUpdateHandlersMu.Lock()
	defer playerListUpdateHandlersMu.Unlock()

	playerListUpdateHandlers = append(playerListUpdateHandlers, h)
}

var initPlayerListUpdateHandlerOnce sync.Once

func handleLeavePlayer(cc *proxy.ClientConn) {
	name := cc.Name()

	playerListMu.Lock()
	delete(playerList, name)
	playerListMu.Unlock()

	playerListUpdateHandlersMu.RLock()
	defer playerListUpdateHandlersMu.RUnlock()

	for _, h := range playerListUpdateHandlers {
		if h.Leave != nil {
			h.Leave(name)
		}
		if h.Update != nil {
			h.Update(Players())
		}
	}
}

func handleJoinPlayer(cc *proxy.ClientConn) {
	name := cc.Name()

	playerListMu.Lock()
	playerList[name] = cc
	playerListMu.Unlock()

	playerListUpdateHandlersMu.RLock()
	defer playerListUpdateHandlersMu.RUnlock()

	for _, h := range playerListUpdateHandlers {
		if h.Join != nil {
			h.Join(name)
		}
		if h.Update != nil {
			h.Update(Players())
		}
	}
}

func initProxyPlayerlistHandlers() {
	initPlayerListUpdateHandlerOnce.Do(func() {
		proxy.RegisterClientHandler(&proxy.ClientHandler{
			Join: func(cc *proxy.ClientConn) string {
				handleJoinPlayer(cc)
				return ""
			},
			Leave: func(cc *proxy.ClientConn, _ *proxy.Leave) {
				handleLeavePlayer(cc)
			},
		})
	})
}
