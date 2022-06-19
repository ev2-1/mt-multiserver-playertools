package playerTools

import (
	"github.com/HimbeerserverDE/mt-multiserver-proxy"

	"sync"
	"fmt"
)

var srvLists = make(map[string]PlayerList)
var srvListsMu sync.RWMutex

func ServerPlayerLists() map[string]PlayerList {
	srvListsMu.RLock()
	defer srvListsMu.RUnlock()

	return srvLists
}

func ServerPlayerList(server string) PlayerList {
	srvListsMu.RLock()
	defer srvListsMu.RUnlock()

	return srvLists[server]
}

func ServerPlayers(server string) int {
	srvListsMu.RLock()
	defer srvListsMu.RUnlock()

	return len(srvLists[server])
}

func srvJoin(cc *proxy.ClientConn, server string) {
	name := cc.Name()

	fmt.Println("Join", name, server)

	srvListsMu.Lock()
	
	if srvLists[server] == nil {
		srvLists[server] = make(PlayerList)
	}

	srvLists[server][name] = cc
	srvListsMu.Unlock()

	updateSrvPlayerList(server)
	updateSrvPlayerListJoin(name, server)
	updateSrvPlayerListGlobal()
}

func srvLeave(cc *proxy.ClientConn, server string) {
	name := cc.Name()

	fmt.Println("Leave", name, server)

	srvListsMu.Lock()
	if srvLists[server] != nil {
		delete(srvLists[server], name)
	}
	srvListsMu.Unlock()

	updateSrvPlayerList(server)
	updateSrvPlayerListLeave(name, server)
	updateSrvPlayerListGlobal()
}

func updateSrvPlayerList(srv string) {
	list := ServerPlayerList(srv)

	srvPlayerListHandlersMu.RLock()
	defer srvPlayerListHandlersMu.RUnlock()

	for _, h := range srvPlayerListHandlers {
		if h.Update != nil {
			h.Update(list, srv)
		}
	}
}

func updateSrvPlayerListGlobal() {
	list := ServerPlayerLists()

	srvPlayerListHandlersMu.RLock()
	defer srvPlayerListHandlersMu.RUnlock()

	for _, h := range srvPlayerListHandlers {
		if h.UpdateGlobal != nil {
			h.UpdateGlobal(list)
		}
	}
}

func updateSrvPlayerListJoin(name, server string) {
	srvPlayerListHandlersMu.RLock()
	defer srvPlayerListHandlersMu.RUnlock()

	for _, h := range srvPlayerListHandlers {
		if h.Join != nil {
			h.Join(name, server)
		}
	}
}

func updateSrvPlayerListLeave(name, server string) {
	srvPlayerListHandlersMu.RLock()
	defer srvPlayerListHandlersMu.RUnlock()

	for _, h := range srvPlayerListHandlers {
		if h.Leave != nil {
			h.Leave(name, server)
		}
	}
}

var initSrvListsMu sync.Once

func initSrvLists() {
	initSrvListsMu.Do(func() {
		fmt.Println("HEREHREHHREHREHRHEHREHRHEHREHRHE")
	
		proxy.RegisterClientHandler(&proxy.ClientHandler{
			// Join
			AOReady: func(cc *proxy.ClientConn) {
				srvJoin(cc, cc.ServerName())
			},
			Leave: func(cc *proxy.ClientConn, _ *proxy.Leave) {
				srvLeave(cc, cc.ServerName())
			},
			Hop: func(cc *proxy.ClientConn, s, d string) {
				srvLeave(cc, s)
				srvJoin(cc, d)
			},
		})
	})
}

type SrvPlayerListHandler struct {
	Join         func(name, server string)
	Leave        func(name, server string)
	Update       func(names PlayerList, server string)
	UpdateGlobal func(map[string]PlayerList)
}

var srvPlayerListHandlers []*SrvPlayerListHandler
var srvPlayerListHandlersMu sync.RWMutex

func RegisterSrvPlayerListHandler(h *SrvPlayerListHandler) {
	initSrvLists()

	srvPlayerListHandlersMu.Lock()
	defer srvPlayerListHandlersMu.Unlock()

	srvPlayerListHandlers = append(srvPlayerListHandlers, h)
}
