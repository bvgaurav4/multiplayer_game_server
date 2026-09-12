use std::net::{TcpListener,TcpStream};
use std::io::{Read,Write};

fn handle_stream(mut stream: TcpStream) {
    let mut buffer = [0; 1024];
    loop {
        match stream.read(&mut buffer) {
            Ok(0) => {
                println!("Client disconnected");
                break;
            }
            Ok(n) => {
                println!("Received: {}", String::from_utf8_lossy(&buffer[..n]));
                stream.write_all(&buffer[..n]).unwrap(); // Echo back
            }
            Err(e) => {
                eprintln!("Error: {}", e);
                break;
            }
        }
    }
}
pub fn socket() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:7878")?;
    // Binding to Socket address.
    println!("Server listening on port 7878");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                std::thread::spawn(|| handle_stream(stream));
            }
            Err(e) => {
                eprintln!("Connection failed: {}", e);
            }
        }
    }

    Ok(())
}