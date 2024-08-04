package websocket

import (
	"log"
	"sync"
)

type HubManager struct {
	hubs  map[string]*Hub
	mutex sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (manager *HubManager) CreateHub(matchId string) *Hub {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	hub, exists := manager.hubs[matchId]
	if exists {
		log.Printf("Hub already exists for match ID %s", matchId)
		return hub
	}

	hub = NewHub()
	manager.hubs[matchId] = hub
	go hub.Run()

	log.Printf("Created and started new hub for matchId %s", matchId)

	return hub
}

func (manager *HubManager) GetHub(matchId string) *Hub {
	hub, exists := manager.hubs[matchId]
	if !exists {
		log.Printf("Hub for matchId %s does not exist. Creating new hub.", matchId)

		return manager.CreateHub(matchId)
	}

	log.Printf("Hub for matchId %s exists. Returning existing hub.", matchId)
	return hub
}

func (manager *HubManager) DeleteHub(matchId string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.hubs, matchId)
}

func (manager *HubManager) JoinHub(matchId string, client *Client) {
	log.Printf("Client joining hub for matchId %s", matchId)

	hub := manager.GetHub(matchId)
	hub.register <- client

	log.Printf("Client joined hub for matchId %s", matchId)
}

func (manager *HubManager) LeaveHub(matchId string, client *Client) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if hub, exists := manager.hubs[matchId]; exists {
		hub.unregister <- client

		if len(hub.clients) == 0 {
			manager.DeleteHub(matchId)
		}
		log.Printf("Client left hub for matchId %s", matchId)
	}
}
