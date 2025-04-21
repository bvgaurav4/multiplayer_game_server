const WebSocket=require('ws')
var uuid=require('uuid-random')
const wss =new WebSocket.WebSocketServer({port:8080},()=>{
    console.log('server started')
}) 
var playersData={
    "type":"playerData"
}
//when the client connects to the server
wss.on('connection',function connection(client){
//giving a unique identity to our client
client.id=uuid();
//to stord the clients data so that v can access that later
var Clients=[];
var clientids=[];
client.send(`{"id":"${client.id}","hp":100}`)
clientids.push(client.id);
//sending message to client 
client.on('message',(data)=>{
    var dataJason=JSON.parse(data);
    if(clientids.includes(dataJason["id"]))
    {
        Clients.pop(dataJason);
        Clients.push(dataJason);
    }
    client.send(JSON.stringify(Clients))
    console.log("player message",dataJason);
})
//when the connection is closed or when the client get dissconnected
client.on('close',()=>{
    console.log("connection is closed for " + client.id)
})
})
//this listens to the client 
wss.on('listening',()=>{
    console.log('listning on 8080')
}) 
