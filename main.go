package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Аргументы
	job := flag.String("job", "sysmetrics", "job name for Pushgateway")
	instance := flag.String("instance", "", "instance label (default = hostname)")
	url := flag.String("url", "https://push.it25.su", "Pushgateway URL")
	install := flag.String("install", "", "install metric in format app=value, e.g. kompas3d=1")
	user := flag.String("user", "", "basic auth username")
	pass := flag.String("pass", "", "basic auth password")
	flag.Parse()

	if *install == "" {
		log.Fatalf("must specify --install app=value")
	}

	// instance = hostname если не задан
	if *instance == "" {
		h, _ := os.Hostname()
		*instance = h
	}

	// Парсим install
	parts := strings.SplitN(*install, "=", 2)
	if len(parts) != 2 {
		log.Fatalf("invalid --install format, expected app=value, got %q", *install)
	}
	app := parts[0]
	val := parts[1]

	// Формируем метрику в формате Prometheus
	metricsText := fmt.Sprintf("install{app=%q} %s\n", app, val)

	// Формируем URL
	pushURL := fmt.Sprintf("%s/metrics/job/%s/instance/%s", strings.TrimRight(*url, "/"), *job, *instance)

	// Отправляем POST
	req, err := http.NewRequest("POST", pushURL, strings.NewReader(metricsText))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	// Basic Auth
	if *user != "" && *pass != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(*user + ":" + *pass))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to push metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatalf("push failed: status %d", resp.StatusCode)
	}

	fmt.Printf("✅ pushed metric to %s\n%s", pushURL, metricsText)
}
