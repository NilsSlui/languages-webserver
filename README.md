# Web Server Benchmark

This project benchmarks a simple webserver implemented in multiple languages, measuring their performance on handling various HTTP requests. The benchmark evaluates the average response time under simulated loads.

## How It Works

### Supported Languages
- **PHP**
- **Go**
- **Python**

### Benchmark Workflow
1. **Setup**: The script (attempts to) install dependencies for each language 
2. **Server Start**: For each language, a web server is started on localhost:8008.
3. **Validation**: The server is checked for proper setup by issuing a test request.
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

### Prerequisites
- **Go**: Version 1.18 or higher
- **PHP**: Version 7.4 or higher
- **Python**: Version 3.6or higher 
- **Node.js**: Version 22 or higher

## If the installation of dependencies fails, you can manually install the dependencies for each language
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
🎉🏆 GO #1 Benchmark Winner 🏆🎉
Benchmark Results:
go: avg time = 0.001743 seconds 
[██████████]
php: avg time = 0.002708 seconds 
[████████████████]
node: avg time = 0.004474 seconds 
[███████████████████████████]
python: avg time = 0.008137 seconds 
[██████████████████████████████████████████████████]

Overall Average Batch Time: 0.004265 seconds
```

