package playerTools

import (
	proxy "github.com/HimbeerserverDE/mt-multiserver-proxy"

	"sync"
)

type ActivePlayerHandler struct {
	Activate func(cc *proxy.ClientConn)
}

var activePlayerHandlers []*ActivePlayerHandler
var activePlayerHandlersMu sync.RWMutex

func RegisterActivePlayerHandler(h *ActivePlayerHandler) {
	registerProxyActivePlayerHandlers()

	activePlayerHandlersMu.Lock()
	defer activePlayerHandlersMu.Unlock()

	activePlayerHandlers = append(activePlayerHandlers, h)
}

var activePlayersMu sync.RWMutex
var activePlayers = make(PlayerList)

func ActivePlayers() PlayerList {
	activePlayersMu.Lock()
	defer activePlayersMu.Unlock()

	return activePlayers
}

func activatePlayer(cc *proxy.ClientConn) {
	activePlayersMu.Lock()
	activePlayers[cc.Name()] = true
	activePlayersMu.Unlock()

	for _, h := range activePlayerHandlers {
		if h.Activate != nil {
			h.Activate(cc)
		}
	}
}

var registerProxyActivePlayerHandlersMu sync.Once

func registerProxyActivePlayerHandlers() {
	registerProxyActivePlayerHandlersMu.Do(func() {
		proxy.RegisterClientHandler(&proxy.ClientHandler{
			StateChange: func(cc *proxy.ClientConn, _, state proxy.ClientState) {
				if state == proxy.CsActive {
					activatePlayer(cc)
				}
			},
		})

		RegisterPlayerListUpdateHandler(&PlayerListUpdateHandler{
			Leave: func(names ...string) {
				activePlayersMu.Lock()
				for _, name := range names {
					delete(activePlayers, name)
				}
				activePlayersMu.Unlock()
			},
		})
	})
}
