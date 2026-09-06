use std::io::{Read, Write};
use std::net::TcpStream;
use std::str;

fn main() -> std::io::Result<()> {
    // 1. Connect to the server socket
    let mut stream = TcpStream::connect("127.0.0.1:7878")?;
    println!("Successfully connected to server on port 8080");

    // 2. Send data to the server
    let msg = "Hello from the standard client!";
    stream.write_all(msg.as_bytes())?;
    println!("Sent: {}", msg);

    // 3. Receive a response from the server
    let mut buffer = [0; 512];
    let bytes_read = stream.read(&mut buffer)?;
    
    let response = str::from_utf8(&buffer[..bytes_read])
        .unwrap_or("Invalid UTF-8 data");
    println!("Received: {}", response);

    Ok(())
}