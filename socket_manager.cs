using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using WebSocketSharp;
using Newtonsoft.Json.Linq;

public class scoker_manager : MonoBehaviour
{
    WebSocket socket;//manages our websockets

    public GameObject player; //creating a game object when a client joins the game
    public PlayerData playerdata;
    void Start()
    {
        socket = new WebSocket("ws://localhost:8080");
        socket.Connect();//this establish a connection between client and server
        // here the client is connecting to the server
        socket.OnMessage += (sender, e) =>
         {
              if (e.IsText)
              {
                  JObject jsonObj = JObject.Parse(e.Data);
                  if (jsonObj["id"] != null)
                  {
                      PlayerData tempPlayerData = JsonUtility.FromJson<PlayerData>(e.Data);
                      playerdata = tempPlayerData;
                      Debug.Log("player id is " + playerdata.id);
                      return;
                  }
              }
         };
        socket.OnClose += (sender,e)=>
        {
                Debug.Log(e.Code);
                Debug.Log(e.Reason);
                Debug.Log("connection closed");
        };

    }

    void Update()
    {
        if (socket == null)
        {
            return;
        }
        if (player != null && playerdata.id != "")
        {
            playerdata.xPos = player.transform.position.x;
            playerdata.yPos = player.transform.position.y;
            playerdata.zPos = player.transform.position.z;

            System.DateTime epochstart = new System.DateTime(1970, 1, 1, 8, 0, 0, System.DateTimeKind.Utc);
            double timestamp = (System.DateTime.UtcNow - epochstart).TotalSeconds;
            playerdata.timestamp = timestamp;
            string playerDataJSON = JsonUtility.ToJson(playerdata);
            socket.Send(playerDataJSON);
        }
    }
    private void OnDestroy()
    {
        socket.Close();
    }
}
