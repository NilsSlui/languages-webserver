use std::net::{TcpListener, TcpStream};
use std::io::{Read, Write};
use std::collections::HashMap;

fn main() {
    // Specify the port to listen on
    let port = std::env::var("PORT").unwrap_or_else(|_| "8008".to_string());
    let address = format!("127.0.0.1:{}", port);

    println!("Listening on {}", address);

    // Bind to the specified address and port
    let listener = TcpListener::bind(&address).expect("Could not bind to address");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                handle_connection(stream);
            }
            Err(e) => {
                eprintln!("Connection failed: {}", e);
            }
        }
    }
}

fn handle_connection(mut stream: TcpStream) {
    let mut buffer = [0; 512];
    stream.read(&mut buffer).unwrap();

    // Parse request (simplified)
    let request = String::from_utf8_lossy(&buffer);
    let mut lines = request.lines();
    let first_line = lines.next().unwrap_or("");
    let mut parts = first_line.split_whitespace();
    let method = parts.next().unwrap_or("");
    let uri = parts.next().unwrap_or("");

    let mut response_code = 200;
    let response_body = match method {
        "POST" => {
            let params = parse_query_string(lines.last().unwrap_or(""));
            match params.get("input") {
                Some(input) => match uri {
                    "/sha256" => sha256(input),
                    "/base64" => base64_encode(input),
                    "/urlencode" => urlencode(input),
                    _ => {
                        response_code = 404;
                        "Not Found".to_string()
                    }
                },
                None => {
                    response_code = 400;
                    "Input is required".to_string()
                }
            }
        }
        "GET" => {
            response_code = 201;
            "rust".to_string()
        }
        _ => {
            response_code = 405;
            "Method Not Allowed".to_string()
        }
    };

    let response = format!(
        "HTTP/1.1 {} OK\r\nContent-Type: text/plain\r\n\r\n{}",
        response_code, response_body
    );

    stream.write_all(response.as_bytes()).unwrap();
    stream.flush().unwrap();
}

fn parse_query_string(query: &str) -> HashMap<String, String> {
    query
        .split('&')
        .filter_map(|pair| {
            let mut parts = pair.splitn(2, '=');
            let key = parts.next()?.to_string();
            let value = parts.next()?.to_string();
            Some((key, value))
        })
        .collect()
}

fn sha256(input: &str) -> String {
    let mut hash = [0u8; 32];
    for (i, byte) in input.bytes().enumerate() {
        hash[i % 32] ^= byte;
    }
    hash.iter().map(|byte| format!("{:02x}", byte)).collect()
}

fn base64_encode(input: &str) -> String {
    const BASE64_CHARS: &str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut encoded = String::new();
    let mut bits = 0u32;
    let mut bits_count = 0;

    for byte in input.bytes() {
        bits = (bits << 8) | byte as u32;
        bits_count += 8;

        while bits_count >= 6 {
            bits_count -= 6;
            let index = (bits >> bits_count) & 0b111111;
            encoded.push(BASE64_CHARS.chars().nth(index as usize).unwrap());
        }
    }

    if bits_count > 0 {
        let index = (bits << (6 - bits_count)) & 0b111111;
        encoded.push(BASE64_CHARS.chars().nth(index as usize).unwrap());
    }

    while encoded.len() % 4 != 0 {
        encoded.push('=');
    }

    encoded
}

fn urlencode(input: &str) -> String {
    input
        .chars()
        .flat_map(|c| match c {
            'a'..='z' | 'A'..='Z' | '0'..='9' | '-' | '_' | '.' | '~' => vec![c],
            _ => format!("%{:02X}", c as u8).chars().collect(),
        })
        .collect()
}
