package playerTools

import (
	"github.com/HimbeerserverDE/mt-multiserver-proxy"

	"sync"
)

var _AOReadyClts PlayerList
var _AOReadyCltsMu sync.RWMutex

func AOReadyClts() PlayerList {
	_AOReadyCltsMu.RLock()
	defer _AOReadyCltsMu.RUnlock()

	return _AOReadyClts
}

func _AOReadyClt(cc *proxy.ClientConn) {
	_AOReadyCltsMu.Lock()
	defer _AOReadyCltsMu.Unlock()

	if _AOReadyClts == nil {
		_AOReadyClts = make(PlayerList)
	}

	_AOReadyClts[cc.Name()] = cc
}

// aka leave
func _AOUnreadyClt(cc *proxy.ClientConn) {
	_AOReadyCltsMu.Lock()
	defer _AOReadyCltsMu.Unlock()

	if _AOReadyClts == nil {
		_AOReadyClts = make(PlayerList)
	}

	delete(_AOReadyClts, cc.Name())
}

var _AOReadyCltInitMu sync.Once

func initAOReadyClts() {
	_AOReadyCltInitMu.Do(func() {
		proxy.RegisterClientHandler(&proxy.ClientHandler{
			AOReady: func(cc *proxy.ClientConn) {
				_AOReadyClt(cc)
			},
			Leave: func(cc *proxy.ClientConn, _ *proxy.Leave) {
				_AOUnreadyClt(cc)
			},
		})
	})
}
