package websocket

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	JoinLobby   = "joinLobby"
	LeaveLobby  = "leaveLobby"
	StartGame   = "startGame"
	UpdateState = "updateState"
	PlayerMove  = "playerMove"
	Initial     = "initial"
	PlayerLeave = "playerLeave"
)
