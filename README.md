# Web Server Benchmark

This project benchmarks a simple webserver implemented in multiple languages. The benchmark evaluates the average response time under simulated loads.

## How It Works

### Supported Languages
- **Go**: Version 1.18 or higher
- **PHP**: Version 7.4 or higher
- **Python**: Version 3.6or higher 
- **Node.js**: Version 22 or higher
- **Ruby**: Version 2.6 or higher
- **Rust**: 2021 

### Benchmark Workflow
1. **Setup**: The script (attempts to) install dependencies for each language 
2. **Server Start**: For each language, a web server is started on localhost:8008.
4. **Simulated Load**:
   - 16 batches of requests are sent, with 8 simultaneous requests per batch.
   - Each request targets one of the following endpoints:
     - `/sha256`: Calculates the SHA-256 hash of the input.
     - `/base64`: Encodes the input in Base64.
     - `/urlencode`: URL-encodes the input string.
   - Random 64 character strings are used as input for each request.
5. **Performance Measurement**:
   - The overall average time across batches is recorded.
   - An ascii bar chart is displayed

## Running the Benchmark

### If the installation of dependencies fails, you can manually install the dependencies for each language
- **Python**: Flask
  ```bash
  pip install Flask
  ```

### Steps
1. Clone the repository and navigate to the root directory:
   ```bash
   git clone https://github.com/NilsSlui/languages-webserver.git
   cd languages-webserver
   ```

2. Run the benchmark script:
   ```bash 
   go run benchmark.go
   ```

## Results
```
Benchmark Results:
go: avg 0.001706 seconds 
[█████████]
php: avg 0.002920 seconds 
[████████████████]
rust: avg 0.003013 seconds 
[████████████████]
node: avg 0.004119 seconds 
[██████████████████████]
ruby: avg 0.006387 seconds 
[███████████████████████████████████]
python: avg 0.009086 seconds 
[██████████████████████████████████████████████████]

Overall Average Batch Time: 0.004539 secondsg
```

