package playerTools

import (
	proxy "github.com/HimbeerserverDE/mt-multiserver-proxy"

	"sync"
)

type PlayerList map[string]bool

type PlayerListUpdateHandler struct {
	Join   func(...string)
	Leave  func(...string)
	Update func(PlayerList)
}

var playerListUpdateHandlers []*PlayerListUpdateHandler
var playerListUpdateHandlersMu sync.RWMutex

var playerListMu sync.RWMutex
var playerList PlayerList

func Players() PlayerList {
	playerListMu.Lock()
	defer playerListMu.Unlock()

	return playerList
}

func RegisterPlayerListUpdateHandler(h *PlayerListUpdateHandler) {
	registerProxyPlayerlistHandlers()

	playerListUpdateHandlersMu.Lock()
	defer playerListUpdateHandlersMu.Unlock()

	playerListUpdateHandlers = append(playerListUpdateHandlers, h)
}

var initPlayerListUpdateHandlerOnce sync.Once

func handleLeavePlayer(name string) {
	playerListMu.Lock()
	delete(playerList, name)
	playerListMu.Unlock()

	playerListUpdateHandlersMu.RLock()
	defer playerListUpdateHandlersMu.RUnlock()
	
	for _, h:= range playerListUpdateHandlers {
		if h.Leave != nil {
			h.Leave(name)
		}
		if h.Update != nil {
			h.Update(Players())
		}
	}
}

func handleJoinPlayer(name string) {
	playerListMu.Lock()
	playerList[name] = true
	playerListMu.Unlock()

	playerListUpdateHandlersMu.RLock()
	defer playerListUpdateHandlersMu.RUnlock()

	for _, h:= range playerListUpdateHandlers {
		if h.Leave != nil {
			h.Leave(name)
		}
		if h.Update != nil {
			h.Update(Players())
		}
	}
}

func registerProxyPlayerlistHandlers() {
	initPlayerListUpdateHandlerOnce.Do(func() {
		proxy.RegisterClientHandler(&proxy.ClientHandler{
			Join: func(cc *proxy.ClientConn) string {
				handleJoinPlayer(cc.Name())
				return ""
			},
			Leave: func(cc *proxy.ClientConn, _ *proxy.Leave) {
				handleLeavePlayer(cc.Name())
			},
		})
	})
}
