// use std::env;
use rdkafka::config::ClientConfig;
use rdkafka::producer::FutureProducer;


fn create_kafka_producer() -> FutureProducer {
    let bootstrap_servers =
        env::var("KAFKA_BOOTSTRAP_SERVERS")
            .unwrap_or_else(|_| "localhost:9092".to_string());

    ClientConfig::new()
        .set("bootstrap.servers", &bootstrap_servers)
        .set("acks", "all")
        .set("enable.idempotence", "true")
        .create()
        .expect("Failed to create Kafka producer")
}