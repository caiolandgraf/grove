package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleCompose = `services:
  db:
    image: postgres:latest
    ports:
      - ${DB_PORT:-5432}:5432
    volumes:
      - db:/var/lib/postgresql

  redis:
    image: redis:7-alpine
    ports:
      - '6379:6379'
    volumes:
      - redis_data:/data

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - '16686:16686'

  loki:
    image: grafana/loki:3.5.0
    ports:
      - '3100:3100'
    volumes:
      - loki_data:/loki

  promtail:
    image: grafana/promtail:3.5.0
    depends_on:
      - loki

  prometheus:
    image: prom/prometheus:v3.4.1
    ports:
      - '9090:9090'
    volumes:
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:11.6.0
    ports:
      - '3000:3000'
    depends_on:
      - prometheus
      - loki
      - jaeger
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  redis_data:
  db:
  loki_data:
  prometheus_data:
  grafana_data:
`

func readFixture(t *testing.T) []string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "grove-base-docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(string(raw), "\n")
}

func assertServices(t *testing.T, content string, want []string) {
	t.Helper()
	for _, svc := range want {
		if !strings.Contains(content, "  "+svc+":") {
			t.Errorf("expected service %q in compose:\n%s", svc, content)
		}
	}
}

func assertServicesAbsent(t *testing.T, content string, absent []string) {
	t.Helper()
	for _, svc := range absent {
		if strings.Contains(content, "  "+svc+":") {
			t.Errorf("expected service %q to be removed:\n%s", svc, content)
		}
	}
}

func TestUpdateDockerComposeAllEnabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(sampleCompose), 0o644); err != nil {
		t.Fatal(err)
	}

	obs := observabilitySelection{
		Jaeger: true, Prometheus: true, Grafana: true, Loki: true, Promtail: true,
	}
	if err := updateDockerCompose(path, obs); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	assertServices(t, string(content), []string{
		"db", "redis", "jaeger", "loki", "promtail", "prometheus", "grafana",
	})
}

func TestUpdateDockerComposeAllDisabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(sampleCompose), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := updateDockerCompose(path, observabilitySelection{}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	assertServices(t, text, []string{"db", "redis"})
	assertServicesAbsent(t, text, []string{
		"jaeger", "loki", "promtail", "prometheus", "grafana",
	})
}

func TestUpdateDockerComposeRealTemplateAllEnabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	raw, err := os.ReadFile(filepath.Join("testdata", "grove-base-docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	obs := observabilitySelection{
		Jaeger: true, Prometheus: true, Grafana: true, Loki: true, Promtail: true,
	}
	if err := updateDockerCompose(path, obs); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	assertServices(t, string(content), []string{
		"db", "redis", "jaeger", "loki", "promtail", "prometheus", "grafana",
	})
}

func TestUpdateDockerComposeRealTemplateAllDisabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "docker-compose.yml")
	raw, err := os.ReadFile(filepath.Join("testdata", "grove-base-docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := updateDockerCompose(path, observabilitySelection{}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	assertServices(t, text, []string{"db", "redis"})
	assertServicesAbsent(t, text, []string{
		"jaeger", "loki", "promtail", "prometheus", "grafana",
	})
}

func TestPruneUnusedVolumesDoesNotDropServices(t *testing.T) {
	lines := readFixture(t)
	used := collectNamedVolumes(lines)
	out := pruneUnusedVolumes(lines, used)
	content := strings.Join(out, "\n")
	assertServices(t, content, []string{
		"db", "redis", "jaeger", "loki", "promtail", "prometheus", "grafana",
	})
}
