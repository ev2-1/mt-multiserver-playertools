package playerTools

import (
	"github.com/HimbeerserverDE/mt-multiserver-proxy"
	"github.com/anon55555/mt"

	"sync"
)

type ClickHandler struct {
	Pos [3]int16
	Srv string

	Handler func(cc *proxy.ClientConn)
}

func RegisterClick(c *ClickHandler) {
	initRegisterClick()

	clickHandlersMu.Lock()
	defer clickHandlersMu.Unlock()
	clickHandlers = append(clickHandlers, c)
}

var clickHandlers []*ClickHandler
var clickHandlersMu sync.RWMutex

var initRegisterClickMu sync.Once

func initRegisterClick() {
	initRegisterClickMu.Do(func() {
		proxy.RegisterInteractionHandler(proxy.InteractionHandler{
			Type: proxy.AnyInteraction,
			Handler: func(cc *proxy.ClientConn, cmd *mt.ToSrvInteract) bool {
				var pointedPos [3]int16
				switch p := cmd.Pointed.(type) {
				case *mt.PointedNode:
					pointedPos = p.Under
				default:
					return false
				}

				if cmd.Action == mt.Dig { // if click
					go func() {
						clickHandlersMu.RLock()
						defer clickHandlersMu.RUnlock()
						for _, v := range clickHandlers {
							if v.Pos == pointedPos && v.Srv == cc.ServerName() {
								v.Handler(cc)
							}
						}
					}()
				}

				return false
			},
		})
	})
}
