package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Command             []string
	DependenciesCommand [][]string
	Cwd                 string
}

type Benchmark struct {
	Name  string
	Value float64
}

var SERVERS = map[string]Server{
	"php": {
		Command: []string{
			"php",
			"-d", "opcache.enable_cli=1",
			"-S", "127.0.0.1:8008",
			"webservers/php/webserver.php",
		},
		Cwd: "",
	},
	"go": {
		Command: []string{"./webserver"}, // Run the compiled binary
		DependenciesCommand: [][]string{
			{"go", "mod", "init", "webserver"},
			{"go", "build", "-o", "webserver", "webserver.go"}, // Build the binary
		},
		Cwd: "webservers/go",
	},
	"python": {
		Command: []string{"python3", "webserver.py"},
		Cwd:     "webservers/python",
	},
}

func installDependencies(server Server) {
	if len(server.DependenciesCommand) > 0 {
		cwd := server.Cwd
		for _, command := range server.DependenciesCommand {
			cmd := exec.Command(command[0], command[1:]...)
			if cwd != "" {
				cmd.Dir = cwd
			}
			cmd.Stdout = io.Discard
			cmd.Stderr = io.Discard
			err := cmd.Run()
			if err != nil {
				//fmt.Println("Error installing dependencies")
			}
		}
	}
}

func startServer(server Server) *exec.Cmd {
	cwd := server.Cwd
	cmd := exec.Command(server.Command[0], server.Command[1:]...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
	time.Sleep(2 * time.Second)
	return cmd
}

func stopServer(cmd *exec.Cmd) {
	err := cmd.Process.Kill()
	if err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}
	cmd.Wait()
	resp, err := http.Get("http://127.0.0.1:8008/")
	if err == nil {
		if resp.StatusCode == 200 {
			fmt.Println("Server is still running")
		}
		resp.Body.Close()
	}
}

func sendRequest(urlStr string, data url.Values) (float64, int) {
	startTime := time.Now()
	resp, err := http.PostForm(urlStr, data)
	elapsedTime := time.Since(startTime).Seconds()
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return 0, 0
	}
	defer resp.Body.Close()
	return elapsedTime, resp.StatusCode
}

func randomString(n int) string {
	letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func runBenchmark(server Server, language string) float64 {
	var batchTimes []float64
	batchSize := 16
	requestsPerBatch := 8
	routes := []string{"sha256", "base64", "urlencode"}

	resp, err := http.Get("http://127.0.0.1:8008/")
	if err != nil {
		fmt.Println("Server is not reachable")
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		fmt.Println("unexpected status code:", resp.StatusCode)
		return 0
	}

	getResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0
	}

	if string(getResponse) != language {
		fmt.Println("Server is not the correct language")
		return 0
	}

	for batchNum := 0; batchNum < batchSize; batchNum++ {
		routesBatch := make([]string, requestsPerBatch)
		dataBatch := make([]url.Values, requestsPerBatch)

		for i := 0; i < requestsPerBatch; i++ {
			value := randomString(64)
			route := routes[rand.Intn(len(routes))]
			routesBatch[i] = route
			dataBatch[i] = url.Values{"input": {value}}
		}

		var wg sync.WaitGroup
		wg.Add(requestsPerBatch)
		var mu sync.Mutex
		runtimes := []float64{}

		for i := 0; i < requestsPerBatch; i++ {
			go func(route string, data url.Values) {
				defer wg.Done()
				rt, status := sendRequest(fmt.Sprintf("http://127.0.0.1:8008/%s", route), data)
				if status != 200 {
					fmt.Printf("Request failed: %d\n", status)
				}
				mu.Lock()
				runtimes = append(runtimes, rt)
				mu.Unlock()
			}(routesBatch[i], dataBatch[i])
		}

		wg.Wait()
		avgRuntime := 0.0
		for _, rt := range runtimes {
			avgRuntime += rt
		}
		avgRuntime /= float64(len(runtimes))
		batchTimes = append(batchTimes, avgRuntime)
		time.Sleep(333 * time.Millisecond)
	}

	totalTime := 0.0
	for _, bt := range batchTimes {
		totalTime += bt
	}
	overallAvg := totalTime / float64(len(batchTimes))
	return overallAvg
}

func main() {
	results := make(map[string]float64)

	for serverName, server := range SERVERS {
		fmt.Printf("Running benchmark on %s server...\n", serverName)
		installDependencies(server)
		serverProcess := startServer(server)
		results[serverName] = runBenchmark(server, serverName)
		stopServer(serverProcess)
	}

	benchmarks := make([]Benchmark, 0, len(results))
	for name, value := range results {
		benchmarks = append(benchmarks, Benchmark{Name: name, Value: value})
	}

	sort.Slice(benchmarks, func(i, j int) bool {
		return benchmarks[i].Value < benchmarks[j].Value
	})

	fastest := benchmarks[0]
	printCoolArt(fastest.Name)

	maxValue := benchmarks[len(benchmarks)-1].Value

	fmt.Println("\nBenchmark Results:")
	sum := 0.0
	for _, benchmark := range benchmarks {
		barLength := int((benchmark.Value / maxValue) * 50)
		bar := strings.Repeat("█", barLength)
		color := getColor(benchmark.Value, maxValue)
		fmt.Printf("%s%s%s: avg time = %f seconds \n[%s%s%s]\n",
			color, benchmark.Name, resetColor(),
			benchmark.Value,
			color, bar, resetColor(),
		)
		sum += benchmark.Value
	}

	avg := sum / float64(len(results))
	fmt.Printf("\nOverall Average Batch Time: %f seconds\n", avg)
}

func printCoolArt(name string) {
	fmt.Printf("\n🎉🏆 %s #1 Benchmark Winner 🏆🎉", strings.ToUpper(name))
}

func getColor(value, maxValue float64) string {
	percentage := value / maxValue
	switch {
	case percentage > 0.8:
		return "\033[1;31m" // Red
	case percentage > 0.5:
		return "\033[1;33m" // Yellow
	default:
		return "\033[1;32m" // Green
	}
}

// resetColor resets the color to default
func resetColor() string {
	return "\033[0m"
}
