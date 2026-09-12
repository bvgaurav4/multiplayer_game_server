use std::net::{TcpListener,TcpStream};
mod connectors;
mod config;

use std::io::{Read,Write};
use connectors::sockets::{socket};
use config::kafka_config::create_kafka_producer;
fn main()  {
    let a_ = create_kafka_producer();
    let _ = socket();
}