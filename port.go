package main

import (
	"encoding/json"
	"strings"

	"github.com/ruokeqx/cpolarPorter/cpolar"
)

type PortManager struct {
	portMap map[string]string
}

func NewPortManager() *PortManager {
	return &PortManager{
		portMap: make(map[string]string),
	}
}

func (pm *PortManager) Update(tunnels []cpolar.Tunnel) bool {
	var updated bool = false
	newPortMap := make(map[string]string)

	for _, tunnel := range tunnels {
		srcPort := pm.extractPort(tunnel.Addr)
		dstPort := pm.extractPort(tunnel.PublicUrl)
		if srcPort != "" && dstPort != "" {
			newPortMap[srcPort] = dstPort
		}
	}

	if len(newPortMap) != len(pm.portMap) {
		updated = true
	} else {
		for srcPort, dstPort := range newPortMap {
			if existingDstPort, exists := pm.portMap[srcPort]; !exists || existingDstPort != dstPort {
				updated = true
				break
			}
		}
	}

	pm.portMap = newPortMap
	return updated
}

func (pm *PortManager) extractPort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func (pm *PortManager) Marshal() ([]byte, error) {
	return json.Marshal(pm.portMap)
}
