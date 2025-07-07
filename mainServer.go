package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Position struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	RX    float64 `json:"rx"`
	RY    float64 `json:"ry"`
	RZ    float64 `json:"rz"`
	State bool    `json:"state"`
	Color string  `json:"color"`
	alive bool    `json:"alive"`
}

type PlayerMessage struct {
	PlayerId string   `json:"player_id"`
	PoolId   string   `json:"pool_id"`
	Status   string   `json:"status"`
	Position Position `json:"position"`
}

type ServerMessage struct {
	Src          string              `json:"src"`
	PositionList map[string]Position `json:"position_list"`
	MainMessage  any                 `json:"main_message"`
}

type Player struct {
	Id       string `json:"id"`
	PoolId   string `json:"pool_id"`
	Color    string `json:"color"`
	Status   string `json:"status"`
	Mutex    sync.Mutex
	position Position        `json:"position"`
	Conn     *websocket.Conn `json:"conn"`
}

type Pool struct {
	Id              string
	Blue            []*Player
	Yellow          []*Player
	count           int16
	status          string
	playerPositions map[string]Position
}

var (
	client              []*websocket.Conn
	pools               = make(map[string]*Pool)
	count               int16
	lol                 string
	poolMu              sync.Mutex
	preSetLoctionBlue   [5]Position
	preSetLoctionyellow [5]Position
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
		playerPositions := map[string]Position{}
		pos := Position{
			X: 0,
			Y: 0,
			Z: 0,
		}
		playerPositions[player.Id] = pos
		pools[poolId] = &Pool{
			Id:              poolId,
			Blue:            make([]*Player, 0),
			Yellow:          make([]*Player, 0),
			count:           0,
			status:          "waiting",
			playerPositions: playerPositions,
		}
		pool := pools[poolId]
		pool.Blue = append(pool.Blue, player)
		pool.count += 1
		pools[poolId] = pool
		player.PoolId = poolId
		return "blue"
	}
	key := 0
	if len(poolsarray[key].Blue) > len(poolsarray[key].Yellow) {
		poolsarray[key].count += 1
		poolsarray[key].Yellow = append(poolsarray[key].Yellow, player)
		player.PoolId = poolsarray[key].Id
		return "yellow"
	} else {
		poolsarray[key].count += 1
		poolsarray[key].Blue = append(poolsarray[key].Blue, player)
		player.PoolId = poolsarray[key].Id
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

func poolBroadCast(id string, positionlist map[string]Position, msg any) {
	fmt.Println("in broadcast")
	fmt.Println("pools")
	fmt.Println(pools)

	pool, ok := pools[id]
	if ok {
		fmt.Println("pool")
		fmt.Println(pool)
		if pool != nil {
			welcome := ServerMessage{Src: "server", PositionList: positionlist, MainMessage: msg}
			fmt.Println("blue")

			for _, player := range pool.Blue {
				if player != nil && player.Conn != nil {
					fmt.Println(player.Mutex)
					// player.Mutex.Lock()
					err := player.Conn.WriteJSON(welcome)
					// player.Mutex.Unlock()
					if err != nil {
						fmt.Println("Error writing JSON to Blue player:", err)
					}
				}
			}
			fmt.Println("yellow")

			for _, player := range pool.Yellow {
				if player != nil && player.Conn != nil {
					// player.Mutex.Lock()
					err := player.Conn.WriteJSON(welcome)
					// player.Mutex.Unlock()
					if err != nil {
						fmt.Println("Error writing JSON to Yellow player:", err)
					}
				}
			}
			fmt.Println("done")
			return
		}
	}
}

func gameLogicAndMechanics(w http.ResponseWriter, r *http.Request) {
	if count >= 2 {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()
	playerID := uuid.New().String()
	player := Player{
		Id:     playerID,
		Status: "waiting",
		Conn:   conn,
	}
	color := adding(&player)
	info := make(map[string]string)
	info["player_id"] = playerID
	info["color"] = color
	info["status"] = "waiting"
	info["pool_id"] = player.PoolId
	welcome := ServerMessage{Src: "server", PositionList: nil, MainMessage: info}
	count += 1
	player.Mutex.Lock()
	err = conn.WriteJSON(welcome)
	player.Mutex.Unlock()
	if err != nil {
		return
	}

	for {
		for player.Status == "waiting" {

		}

		countDuration := 5

		if player.Status == "Starting" {
			for i := 0; i <= countDuration; i++ {
				mgs := make(map[string]any)
				mgs["status"] = "Starting"
				mgs[""] = countDuration - i
				countDown := ServerMessage{Src: "server", PositionList: nil, MainMessage: mgs}
				player.Mutex.Lock()
				err = conn.WriteJSON(countDown)
				player.Mutex.Unlock()
				time.Sleep(1 * time.Second)

				if err != nil {
					//pop
					changeStatus(player.PoolId, "waiting")
				}
			}
		}

		changeStatus(player.PoolId, "Playing")

		mgs := make(map[string]any)
		mgs["status"] = "Playing"
		countDown := ServerMessage{Src: "server", PositionList: nil, MainMessage: mgs}
		player.Mutex.Lock()
		err = conn.WriteJSON(countDown)
		player.Mutex.Unlock()

		fmt.Println("here now")
		for player.Status == "Playing" {
			var playerMeessage PlayerMessage
			_, msg, err := conn.ReadMessage()

			// if err != nil {
			// 	count -= 1
			// 	pools[player.PoolId].count -= 1
			// 	if player.Color == "blue" {
			// 		i := 0
			// 		for i = 0; i < len(pools[player.PoolId].Blue); i++ {
			// 			if &pools[player.PoolId].Blue[i].Id == &player.Id {
			// 				break
			// 			}
			// 		}
			// 		pools[player.PoolId].Blue = append(pools[player.PoolId].Blue[:i], pools[player.PoolId].Blue[i:]...)
			// 		return
			// 	}
			// 	if player.Color == "yellow" {
			// 		i := 0
			// 		for i = 0; i < len(pools[player.PoolId].Yellow); i++ {
			// 			if &pools[player.PoolId].Yellow[i].Id == &player.Id {
			// 				break
			// 			}
			// 		}
			// 		pools[player.PoolId].Yellow = append(pools[player.PoolId].Yellow[:i], pools[player.PoolId].Yellow[i:]...)
			// 		return
			// 	}
			// }

			err = json.Unmarshal(msg, &playerMeessage)
			fmt.Println("bruh")
			if err == nil {
				fmt.Printf("Received: %+v\n", playerMeessage)
				poolMu.Lock()
				pools[player.PoolId].playerPositions[playerID] = playerMeessage.Position
				poolMu.Unlock()
				fmt.Println(player)
				fmt.Println(pools[player.PoolId].playerPositions)
				fmt.Println("pools", pools)
				poolBroadCast(player.PoolId, pools[player.PoolId].playerPositions, "")
			} else {
				fmt.Println("somthing is wrong", err)
				break
			}
		}
	}

}

func changeStatus(poolId string, status string) {
	pool := pools[poolId]
	for i := 0; i < len(pool.Blue); i++ {
		pool.Blue[i].Status = status
	}
	for i := 0; i < len(pool.Yellow); i++ {
		pool.Yellow[i].Status = status
	}
}

func settingPree() {

	preSetLoctionBlue[0] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionBlue[1] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionBlue[2] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionBlue[3] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionBlue[4] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}

	preSetLoctionyellow[0] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionyellow[1] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionyellow[2] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionyellow[3] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
	preSetLoctionyellow[4] = Position{X: 0, Y: 0, Z: 0, RX: 0, RY: 0, RZ: 0}
}

func startGame() {
	for {
		for keys := range pools {
			if pools[keys].count == 1 && pools[keys].status != "Starting" {
				pools[keys].status = "Starting"
				fmt.Println("Starting Game ")
				info := make(map[string]string)
				info["status"] = "Starting"
				for i := range pools[keys].Blue {
					pools[keys].Blue[i].Status = "Starting"
					pools[keys].Blue[i].position = preSetLoctionBlue[i]
				}
				for i := range pools[keys].Yellow {
					pools[keys].Yellow[i].Status = "Starting"
					pools[keys].Yellow[i].position = preSetLoctionyellow[i]
				}
				poolBroadCast(keys, nil, info)
			}
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

func main() {
	settingPree()
	go startGame()
	http.HandleFunc("/play", gameLogicAndMechanics)
	port := "8080"
	fmt.Println("WebSocket server started on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
