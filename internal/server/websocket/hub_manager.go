package websocket

import "sync"

type HubManager struct {
	hubs  map[string]*Hub
	mutex sync.RWMutex
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]*Hub),
	}
}

func (manager *HubManager) GetHub(matchId string) *Hub {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	hub, exists := manager.hubs[matchId]
	if !exists {
		return manager.CreateHub(matchId)
	}

	return hub
}

func (manager *HubManager) CreateHub(matchId string) *Hub {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	hub := NewHub()
	manager.hubs[matchId] = hub
	go hub.Run()

	return hub
}

func (manager *HubManager) DeleteHub(matchId string) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.hubs, matchId)
}

func (manager *HubManager) JoinHub(matchId string, client *Client) {
	hub := manager.GetHub(matchId)
	hub.register <- client
}

func (manager *HubManager) LeaveHub(matchId string, client *Client) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if hub, exists := manager.hubs[matchId]; exists {
		hub.unregister <- client

		if len(hub.clients) == 0 {
			manager.DeleteHub(matchId)
		}
	}
}
