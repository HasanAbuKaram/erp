package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"text/template"
	"time"
)

// ServiceStatus holds the status of a service
type ServiceStatus struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

// checkServiceHTTP checks if an HTTP service is reachable by sending a GET request
func checkServiceHTTP(wg *sync.WaitGroup, mu *sync.Mutex, name, host, port string, services *[]ServiceStatus) {
	defer wg.Done()

	address := net.JoinHostPort(host, port)
	url := fmt.Sprintf("http://%s", address)
	status := "down"

	// Attempt to send a GET request
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		status = "up"
	}
	if resp != nil {
		resp.Body.Close()
	}

	mu.Lock()
	*services = append(*services, ServiceStatus{Name: name, Status: status, Address: address})
	mu.Unlock()

	log.Printf("The connection to %s is: %s", name, status)
}

// checkServiceTCP checks if a TCP service (like a database) is reachable by attempting to establish a connection
func checkServiceTCP(wg *sync.WaitGroup, mu *sync.Mutex, name, host, port string, services *[]ServiceStatus) {
	defer wg.Done()

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second) // 5-second timeout
	status := "down"
	if err == nil {
		status = "up"
		conn.Close()
	}

	mu.Lock()
	*services = append(*services, ServiceStatus{Name: name, Status: status, Address: address})
	mu.Unlock()

	log.Printf("The connection to %s is: %s", name, status)
}

// restartContainer restarts a Docker container and logs the output
func restartContainer(containerName string) error {
	cmd := exec.Command("docker", "restart", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error restarting container %s: %s", containerName, string(output))
	}
	return nil
}

// stopContainer stops a Docker container and logs the output
func stopContainer(containerName string) error {
	cmd := exec.Command("docker", "stop", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error stopping container %s: %s", containerName, string(output))
	}
	return nil
}

// startContainer starts a Docker container and logs the output
func startContainer(containerName string) error {
	cmd := exec.Command("docker", "start", containerName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error starting container %s: %s", containerName, string(output))
	}
	return nil
}

// dashboardHandler serves the dashboard page
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var services []ServiceStatus

	// Define the list of services to check
	serviceList := []struct {
		name   string
		host   string
		port   string
		isHTTP bool
	}{
		{"ERP", "erp", "8080", true},
		{"Static Server", "static-server", "8081", true},
		{"Supply Chain", "supply-chain", "8082", true},
		{"db_supply_chain", "db_supply_chain", "5432", false},
	}

	// Check each service status concurrently
	for _, service := range serviceList {
		wg.Add(1)
		if service.isHTTP {
			go checkServiceHTTP(&wg, &mu, service.name, service.host, service.port, &services)
		} else {
			go checkServiceTCP(&wg, &mu, service.name, service.host, service.port, &services)
		}
	}

	wg.Wait()

	// Define the HTML template
	dashboardHTML := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Service Dashboard</title>
		<style>
			body { font-family: Arial, sans-serif; }
			table { width: 60%; margin: 50px auto; border-collapse: collapse; }
			th, td { padding: 10px; text-align: left; border: 1px solid #ddd; }
			th { background-color: #f4f4f4; }
			.up { color: green; }
			.down { color: red; }
			button { margin: 5px; padding: 10px; border: none; border-radius: 5px; color: white; cursor: pointer; }
			.start { background-color: green; }
			.stop { background-color: red; }
		</style>
		<script>
			function refreshPage() {
				setTimeout(function(){
					location.reload();
				}, 5000); // Refresh every 5 seconds
			}
		</script>
	</head>
	<body onload="refreshPage()">
		<h1 style="text-align: center;">Service Dashboard</h1>
		<table>
			<tr><th>Service</th><th>Status</th><th>Actions</th></tr>
		{{range .}}
			<tr>
				<td>{{.Name}}</td>
				<td class="{{.Status}}">{{.Status}}</td>
				<td>
					{{if eq .Status "down"}}
					<form action="/start" method="post" style="display:inline;">
						<input type="hidden" name="service" value="{{.Name}}">
						<button type="submit" class="start">Start</button>
					</form>
					{{else}}
					<form action="/stop" method="post" style="display:inline;">
						<input type="hidden" name="service" value="{{.Name}}">
						<button type="submit" class="stop">Stop</button>
					</form>
					{{end}}
				</td>
			</tr>
		{{end}}
		</table>
	</body>
	</html>`

	tmpl, err := template.New("dashboard").Parse(dashboardHTML)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template parsing error: %v", err)
		return
	}
	err = tmpl.Execute(w, services)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
		return
	}
}

// handleControl handles start and stop requests
func handleControl(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		service := r.FormValue("service")
		var err error
		switch r.URL.Path {
		case "/start":
			err = startContainer(service)
		case "/stop":
			err = stopContainer(service)
		default:
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(w, "Failed to control container", http.StatusInternalServerError)
			log.Printf("Control error: %v", err)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// func main() {
// 	http.HandleFunc("/dashboard", dashboardHandler)
// 	http.HandleFunc("/start", handleControl)
// 	http.HandleFunc("/stop", handleControl)

// 	log.Println("Starting server on :8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatalf("Server failed to start: %v", err)
// 	}
// }
