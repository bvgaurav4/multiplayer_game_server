use reqwest;
use std::collections::HashMap;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Hello, world!");    
    let mut map = HashMap::new();
    map.insert("lang", "rust");
    map.insert("body", "json");

    let response = reqwest::get("https://www.rust-lang.org").await?;
    let body = response.text().await?;

    println!("body = {:?}", body);
    Ok(())
}
