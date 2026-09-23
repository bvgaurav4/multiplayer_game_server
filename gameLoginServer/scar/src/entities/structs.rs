use std::collections::*;

    pub struct Players{
    pub id : String,
    pub connection: String,
    // position: object, will be in the event object
    pub status: String,
    pub events: Vec<Event>,
    pub game_status: String // dead or alive or viewer model 
}

pub struct Match {
    pub players: HashMap<String,String>,
    pub metadata: String
}

pub struct Event{
    pub position: Position,
    pub event_name: String
}

pub struct Position{
    pub rx: f32,
    pub ry: f32,
    pub rz: f32,
    pub x: f32,
    pub y: f32,
    pub z: f32
}

pub struct MatchMetadata{
    pub score: HashMap<String,String>, // players scores 
    pub game_state: String, 
}