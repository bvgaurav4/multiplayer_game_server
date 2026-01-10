use axum::{routing::get, Router, response::IntoResponse};
use redis::Commands;
// use redis::TypedCommands;

async fn hello() -> Result<(), ()> {
    // let client = redis::Client::open("redis://127.0.0.1/")
    //     .map_err(|e| e.to_string())?;

    // let mut conn = client
    //     .get_connection()
    //     .map_err(|e| e.to_string())?;

    // conn.set("game:players", 10)
    //     .map_err(|e| e.to_string())?;

    // let players: i32 = conn
    //     .get("game:players")
    //     .map_err(|e| e.to_string())?;

    // Ok(format!("Hello, players = {}", players))
        let mut r = match redis::Client::open("redis://127.0.0.1") {
            Ok(client) => {
                match client.get_connection() {
                    Ok(conn) => conn,
                    Err(e) => {
                        println!("Failed to connect to Redis: {e}");
                        return Ok(());
                    }
                }
            },
            Err(e) => {
                println!("Failed to create Redis client: {e}");
                return Ok(());
            }
        };
        return Ok(());

}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/hello", get(hello))
        .route("/hello2", get(hello));

    let listener = tokio::net::TcpListener::bind("0.0.0.0:5000")
        .await
        .unwrap();

    axum::serve(listener, app)
        .await
        .unwrap();
}
