package playerTools

import (
	"github.com/HimbeerserverDE/mt-multiserver-proxy"
)

func GetPlayerByName(name string) *proxy.ClientConn {
	for clt := range proxy.Clts() {
		if clt.Name() == name {
			return clt
		}
	}

	return nil
}
