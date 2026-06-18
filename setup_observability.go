package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type observabilitySelection struct {
	Jaeger     bool
	Prometheus bool
	Grafana    bool
	Loki       bool
	Promtail   bool
}

func promptObservability() observabilitySelection {
	defaults := observabilitySelection{
		Jaeger:     true,
		Prometheus: true,
		Grafana:    true,
		Loki:       true,
		Promtail:   true,
	}

	info, err := os.Stdin.Stat()
	if err != nil || (info.Mode()&os.ModeCharDevice) == 0 {
		return defaults
	}

	fmt.Println()
	fmt.Println("  " + bold("Observability stack"))
	fmt.Println("  " + gray("Press Enter to accept the default."))

	reader := bufio.NewReader(os.Stdin)
	obs := observabilitySelection{
		Jaeger: askYesNo(reader, "Enable Jaeger tracing?", defaults.Jaeger),
		Prometheus: askYesNo(
			reader,
			"Enable Prometheus metrics?",
			defaults.Prometheus,
		),
		Loki: askYesNo(reader, "Enable Loki logs?", defaults.Loki),
		Promtail: askYesNo(
			reader,
			"Enable Promtail log shipper?",
			defaults.Promtail,
		),
		Grafana: askYesNo(
			reader,
			"Enable Grafana dashboards?",
			defaults.Grafana,
		),
	}

	if obs.Promtail && !obs.Loki {
		fmt.Println("  " + warn("Promtail requires Loki; disabling Promtail."))
		obs.Promtail = false
	}

	if obs.Grafana && !obs.Jaeger && !obs.Prometheus && !obs.Loki {
		fmt.Println("  " + warn("Grafana enabled without data sources."))
	}

	fmt.Println()
	return obs
}

func askYesNo(reader *bufio.Reader, prompt string, defaultYes bool) bool {
	suffix := "[y/N]"
	if defaultYes {
		suffix = "[Y/n]"
	}

	for {
		fmt.Printf("  %s %s ", prompt, gray(suffix))
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return defaultYes
		}

		value := strings.TrimSpace(strings.ToLower(input))
		switch value {
		case "":
			return defaultYes
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("  " + warn("Please answer y or n."))
		}
	}
}

func (o observabilitySelection) summary() string {
	parts := make([]string, 0, 5)
	if o.Jaeger {
		parts = append(parts, "jaeger")
	}
	if o.Prometheus {
		parts = append(parts, "prometheus")
	}
	if o.Grafana {
		parts = append(parts, "grafana")
	}
	if o.Loki {
		parts = append(parts, "loki")
	}
	if o.Promtail {
		parts = append(parts, "promtail")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ", ")
}

func configureObservability(
	projectDir string,
	obs observabilitySelection,
) error {
	if obs.Promtail && !obs.Loki {
		obs.Promtail = false
	}

	if err := updateEnvExample(projectDir, obs); err != nil {
		return err
	}

	composePath, err := resolveComposePath(projectDir)
	if err != nil {
		return err
	}
	if composePath != "" {
		if err := updateDockerCompose(composePath, obs); err != nil {
			return err
		}
	}

	return nil
}

func resolveComposePath(projectDir string) (string, error) {
	candidates := []string{
		"docker-compose.yml",
		"infra/compose.yml",
		"infra/docker-compose.yml",
	}

	for _, rel := range candidates {
		path := filepath.Join(projectDir, rel)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf(
		"docker compose file not found (checked: %s)",
		strings.Join(candidates, ", "),
	)
}

func updateEnvExample(projectDir string, obs observabilitySelection) error {
	path := filepath.Join(projectDir, ".env.example")
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	lines = upsertEnvValue(lines, "OTEL_ENABLED", boolString(obs.Jaeger))
	lines = upsertEnvValue(lines, "METRICS_ENABLED", boolString(obs.Prometheus))

	if !containsEnvKey(lines, "OTEL_ENABLED") {
		lines = insertBeforeEnv(
			lines,
			"OTEL_SERVICE_NAME",
			"OTEL_ENABLED",
			boolString(obs.Jaeger),
		)
	}
	if !containsEnvKey(lines, "METRICS_ENABLED") {
		lines = insertBeforeHeader(
			lines,
			"# Application",
			"METRICS_ENABLED",
			boolString(obs.Prometheus),
		)
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func upsertEnvValue(lines []string, key, value string) []string {
	prefix := key + "="
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = prefix + value
		}
	}
	return lines
}

func containsEnvKey(lines []string, key string) bool {
	prefix := key + "="
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func insertBeforeEnv(lines []string, key, newKey, value string) []string {
	needle := key + "="
	insert := newKey + "=" + value
	for i, line := range lines {
		if strings.HasPrefix(line, needle) {
			return append(lines[:i], append([]string{insert}, lines[i:]...)...)
		}
	}
	return append(lines, insert)
}

func insertBeforeHeader(lines []string, header, key, value string) []string {
	insert := key + "=" + value
	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			return append(lines[:i], append([]string{insert}, lines[i:]...)...)
		}
	}
	return append(lines, insert)
}

func updateDockerCompose(path string, obs observabilitySelection) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lines = removeServiceBlocks(lines, obs)
	lines = filterGrafanaDependsOn(lines, obs)
	usedVolumes := collectNamedVolumes(lines)
	lines = removeVolumeBlocks(lines, obs)
	lines = pruneUnusedVolumes(lines, usedVolumes)

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

func removeServiceBlocks(
	lines []string,
	obs observabilitySelection,
) []string {
	disabled := map[string]bool{
		"jaeger":     !obs.Jaeger,
		"loki":       !obs.Loki,
		"promtail":   !obs.Promtail,
		"prometheus": !obs.Prometheus,
		"grafana":    !obs.Grafana,
	}

	var out []string
	inServices := false

	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !inServices {
			out = append(out, line)
			if isTopLevelKey(line) && trimmed == "services:" {
				inServices = true
			}
			i++
			continue
		}

		if isTopLevelKey(line) {
			inServices = false
			out = append(out, line)
			i++
			continue
		}

		if isIndentedKey(line, 2) {
			name := strings.TrimSuffix(trimmed, ":")
			if disabled[name] {
				i++
				for i < len(lines) {
					next := lines[i]
					if isTopLevelKey(next) || isIndentedKey(next, 2) {
						break
					}
					i++
				}
				continue
			}
		}

		out = append(out, line)
		i++
	}

	return out
}

func filterGrafanaDependsOn(
	lines []string,
	obs observabilitySelection,
) []string {
	if !obs.Grafana {
		return lines
	}

	disabled := map[string]bool{
		"jaeger":     !obs.Jaeger,
		"loki":       !obs.Loki,
		"prometheus": !obs.Prometheus,
	}

	var out []string
	inGrafana := false

	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if isIndentedKey(line, 2) {
			name := strings.TrimSuffix(trimmed, ":")
			inGrafana = name == "grafana"
			out = append(out, line)
			i++
			continue
		}

		if inGrafana && strings.HasPrefix(line, "    depends_on:") {
			kept := []string{}
			j := i + 1
			for j < len(lines) {
				next := lines[j]
				if !strings.HasPrefix(next, "      - ") {
					break
				}
				name := strings.TrimSpace(strings.TrimPrefix(next, "      - "))
				if !disabled[name] {
					kept = append(kept, next)
				}
				j++
			}
			if len(kept) > 0 {
				out = append(out, line)
				out = append(out, kept...)
			}
			i = j
			continue
		}

		out = append(out, line)
		i++
	}

	return out
}

func removeVolumeBlocks(
	lines []string,
	obs observabilitySelection,
) []string {
	disabled := map[string]bool{
		"loki_data":       !obs.Loki,
		"prometheus_data": !obs.Prometheus,
		"grafana_data":    !obs.Grafana,
	}

	var out []string
	inVolumes := false

	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !inVolumes {
			out = append(out, line)
			if isTopLevelKey(line) && trimmed == "volumes:" {
				inVolumes = true
			}
			i++
			continue
		}

		if isTopLevelKey(line) {
			inVolumes = false
			out = append(out, line)
			i++
			continue
		}

		if isIndentedKey(line, 2) {
			name := strings.TrimSuffix(trimmed, ":")
			if disabled[name] {
				i++
				for i < len(lines) {
					next := lines[i]
					if isTopLevelKey(next) || isIndentedKey(next, 2) {
						break
					}
					i++
				}
				continue
			}
		}

		out = append(out, line)
		i++
	}

	return out
}

func collectNamedVolumes(lines []string) map[string]bool {
	used := make(map[string]bool)
	inServices := false
	inServiceVolumes := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !inServices {
			if isTopLevelKey(line) && trimmed == "services:" {
				inServices = true
			}
			continue
		}

		if isTopLevelKey(line) {
			break
		}

		if isIndentedKey(line, 2) {
			inServiceVolumes = false
			continue
		}

		if isIndentedKey(line, 4) && trimmed == "volumes:" {
			inServiceVolumes = true
			continue
		}

		if inServiceVolumes {
			if isIndentedKey(line, 4) && trimmed != "volumes:" {
				inServiceVolumes = false
				continue
			}
			if strings.HasPrefix(trimmed, "- ") {
				entry := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				if entry == "" {
					continue
				}
				if strings.HasPrefix(entry, "./") ||
					strings.HasPrefix(entry, "/") ||
					strings.HasPrefix(entry, "${") {
					continue
				}
				name := entry
				if idx := strings.Index(entry, ":"); idx >= 0 {
					name = strings.TrimSpace(entry[:idx])
				}
				if name != "" {
					used[name] = true
				}
			}
		}
	}

	return used
}

func pruneUnusedVolumes(lines []string, used map[string]bool) []string {
	if len(used) == 0 {
		return lines
	}

	var out []string
	inVolumes := false
	keptAny := false
	volumeHeaderIndex := -1

	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if !inVolumes {
			if isTopLevelKey(line) && trimmed == "volumes:" {
				inVolumes = true
				keptAny = false
				volumeHeaderIndex = len(out)
				out = append(out, line)
				i++
				continue
			}
			out = append(out, line)
			i++
			continue
		}

		if isTopLevelKey(line) {
			if !keptAny && volumeHeaderIndex >= 0 {
				out = out[:volumeHeaderIndex]
			}
			inVolumes = false
			volumeHeaderIndex = -1
			keptAny = false
			out = append(out, line)
			i++
			continue
		}

		if isIndentedKey(line, 2) {
			name := strings.TrimSuffix(trimmed, ":")
			if used[name] {
				keptAny = true
				out = append(out, line)
				i++
				for i < len(lines) {
					next := lines[i]
					if isTopLevelKey(next) || isIndentedKey(next, 2) {
						break
					}
					out = append(out, next)
					i++
				}
				continue
			}

			i++
			for i < len(lines) {
				next := lines[i]
				if isTopLevelKey(next) || isIndentedKey(next, 2) {
					break
				}
				i++
			}
			continue
		}

		out = append(out, line)
		i++
	}

	if inVolumes && !keptAny && volumeHeaderIndex >= 0 {
		out = out[:volumeHeaderIndex]
	}

	return out
}

func isTopLevelKey(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed != "" && !strings.HasPrefix(line, " ") &&
		strings.HasSuffix(trimmed, ":")
}

func isIndentedKey(line string, indent int) bool {
	if indent <= 0 {
		return false
	}
	prefix := strings.Repeat(" ", indent)
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(line, prefix) &&
		!strings.HasPrefix(line, prefix+" ") &&
		strings.HasSuffix(trimmed, ":")
}
