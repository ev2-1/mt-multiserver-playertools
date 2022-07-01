package playerTools

import (
	"github.com/HimbeerserverDE/mt-multiserver-proxy"
	"github.com/anon55555/mt"

	"sync"
)

type PosHandler struct {
	Update func(*proxy.ClientConn, mt.PlayerPos)
}

var posHandlers   []*PosHandler
var posHandlersMu sync.RWMutex

func RegisterPosHandler(h *PosHandler) {
	InitPos()

	posHandlersMu.Lock()
	defer posHandlersMu.Unlock()

	posHandlers = append(posHandlers, h)
}

var posMap = make(map[string]mt.PlayerPos)
var posMapMu sync.RWMutex

func PosMap() map[string]mt.PlayerPos {
	InitPos()

	posMapMu.RLock()
	defer posMapMu.RUnlock()

	return posMap
}

func GetPos(name string) mt.PlayerPos {
	InitPos()

	posMapMu.RLock()
	defer posMapMu.RUnlock()

	return posMap[name]
}

func setPos(name string, pos mt.PlayerPos) {
	posMapMu.Lock()
	defer posMapMu.Unlock()

	posMap[name] = pos
}

func updatePos(cc *proxy.ClientConn, pos mt.PlayerPos) {
	posHandlersMu.RLock()
	defer posHandlersMu.RUnlock()

	for _, h := range posHandlers {
		if h.Update != nil {
			h.Update(cc, pos)
		}
	}
}

var InitPosMu sync.Once

// InitPos has to be called if not using RegisterPosHandler to register handlers with proxy
func InitPos() {
	InitPosMu.Do(func() {
		proxy.RegisterPacketHandler(&proxy.PacketHandler{
			CltHandler: func(cc *proxy.ClientConn, pkt *mt.Pkt) bool {
				switch cmd := pkt.Cmd.(type) {
				case *mt.ToSrvPlayerPos:
					setPos(cc.Name(), cmd.Pos)
					updatePos(cc, cmd.Pos)
				}
			
				return false
			},
		})
	})
}
