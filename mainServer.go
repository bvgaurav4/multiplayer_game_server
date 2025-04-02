package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Position struct {
	x int64
	y int64
	z int64
}

type PlayerMessage struct {
	PlayerId string   `json:"player_id"`
	PoolId   string   `json:"pool_id"`
	Status   string   `json:"status"`
	Position Position `json:"position"`
	Tile     int      `json:"tilt"`
}

type ServerMessage struct {
	Src          string          `json:"src"`
	PositionList []PlayerMessage `json:"position_list"`
	MainMessage  any             `json:"main_message"`
}

type Player struct {
	Id     string
	status string
	Conn   *websocket.Conn
}

type Pool struct {
	Id     string
	Blue   []*Player
	Yellow []*Player
	count  int16
	status string
}

var (
	client []*websocket.Conn
	pools  = make(map[string]*Pool)
	count  int16
	lol    string
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func adding(player *Player) string {
	var poolsarray []*Pool
	for key := range pools {
		if pools[key].status == "waiting" {
			poolsarray = append(poolsarray, pools[key])
		}
	}
	if len(poolsarray) == 0 {
		poolId := uuid.New().String()
		lol = poolId
		pools[poolId] = &Pool{
			poolId,
			make([]*Player, 0),
			make([]*Player, 0),
			0,
			"waiting",
		}
		pool := pools[poolId]
		pool.Blue = append(pool.Blue, player)
		pool.count += 1
		pools[poolId] = pool
		return "blue"
	}
	key := 0
	if len(poolsarray[key].Blue) > len(poolsarray[key].Yellow) {
		poolsarray[key].count += 1
		poolsarray[key].Yellow = append(poolsarray[key].Yellow, &player)
		return "yellow"
	} else {
		poolsarray[key].count += 1
		poolsarray[key].Blue = append(poolsarray[key].Blue, &player)
		return "blue"
	}
}

// func adding(player Player) {
// 	var poolsarray []*Pool
// 	for key := range pools {
// 		if pools[key].status == "waiting" {
// 			poolsarray = append(poolsarray, pools[key])
// 		}
// 	}
// 	if len(poolsarray) == 0 {
// 		poolId := uuid.New().String()
// 		pools[poolId] = &Pool{
// 			poolId,
// 			make([]*Player, 0),
// 			make([]*Player, 0),
// 			0,
// 			"waiting",
// 		}
// 		pool := pools[poolId]
// 		pool.Blue = append(pool.Blue, &player)
// 		pool.count += 1
// 		pools[poolId] = pool
// 	}
// 	sort.Slice(poolsarray, func(i, j int) bool {
// 		return poolsarray[i].count > poolsarray[j].count
// 	})
// 	for key := range poolsarray {
// 		if poolsarray[key].count%2 != 0 {
// 			if len(poolsarray[key].Blue) > len(poolsarray[key].Yellow) {
// 				fmt.Println("here")
// 				poolsarray[key].count += 1
// 				poolsarray[key].Yellow = append(poolsarray[key].Yellow, &player)
// 				return
// 			} else {
// 				poolsarray[key].count += 1
// 				poolsarray[key].Blue = append(poolsarray[key].Blue, &player)
// 				return
// 			}
// 		}
// 	}
// }

func poolBroadCast(id string, msg []byte) {
	pool := pools[id]
	for i := 0; i < len(pool.Blue); i++ {
		pool.Blue[i].Conn.WriteJSON(msg)
	}
	for i := 0; i < len(pool.Yellow); i++ {
		pool.Yellow[i].Conn.WriteJSON(msg)
	}
}

func gameLogicAndMechanics(w http.ResponseWriter, r *http.Request) {
	if count >= 10 {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	playerID := uuid.New().String()
	player := Player{playerID, "waiting", conn}
	color := adding(player)
	info := make(map[string]string)
	info["player_id"] = playerID
	info["color"] = color
	info["status"] = "waiting"
	welcome := ServerMessage{Src: "server", PositionList: nil, MainMessage: info}
	count += 1
	err = conn.WriteJSON(welcome)

	for player.status == "waiting" {
	}
	info = make(map[string]string)
	info["color"] = color
	info["status"] = "starting"
	welcome = ServerMessage{Src: "server", PositionList: nil, MainMessage: info}
	err = conn.WriteJSON(welcome)
	for {
		var playerMeessage PlayerMessage
		err := conn.ReadJSON(playerMeessage)
		fmt.Printf("Received: %s\n", playerMeessage)
		// poolBroadCast(playerMeessage.PoolId, playerMeessage)
		if err != nil {
			fmt.Println("Read error:", err)
			for i := 0; i < len(client); i++ {
				if client[i] == conn {
					fmt.Println("yup found it")
				}
			}
			break
		}
	}

}
func main() {
	http.HandleFunc("/play", gameLogicAndMechanics)
	port := "8080"
	fmt.Println("WebSocket server started on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
