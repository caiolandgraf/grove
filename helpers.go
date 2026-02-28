package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// ──────────────────────────────────────────────
// Terminal colors (RGB true-color)
// ──────────────────────────────────────────────

const (
	colorReset = "\033[0m"
	colorBold  = "\033[1m"
	colorDim   = "\033[2m"

	// Foreground — RGB true-color
	colorRed    = "\033[38;2;220;60;60m"
	colorGreen  = "\033[38;2;40;210;90m"
	colorYellow = "\033[38;2;230;200;40m"
	colorBlue   = "\033[38;2;80;140;255m"
	colorCyan   = "\033[38;2;80;220;220m"
	colorGray   = "\033[38;2;130;130;145m"

	// Background badges (colored bg + white foreground)
	colorBgGreen  = "\033[48;2;40;180;80m\033[38;2;255;255;255m"
	colorBgBlue   = "\033[48;2;60;120;220m\033[38;2;255;255;255m"
	colorBgRed    = "\033[48;2;195;55;55m\033[38;2;255;255;255m"
	colorBgYellow = "\033[48;2;185;140;20m\033[38;2;255;255;255m"

	// Logo gradient — bright lime-green (top) → forest green (bottom)
	logoG1 = "\033[38;2;200;40;60m"
	logoG2 = "\033[38;2;175;30;50m"
	logoG3 = "\033[38;2;150;20;40m"
	logoG4 = "\033[38;2;125;15;30m"
	logoG5 = "\033[38;2;100;10;22m"

	logoR1 = "\033[38;2;255;100;100m"
	logoR2 = "\033[38;2;245;70;70m"
	logoR3 = "\033[38;2;230;50;50m"
	logoR4 = "\033[38;2;210;35;35m"
	logoR5 = "\033[38;2;185;20;20m"
)

func success(msg string) string {
	return colorBgGreen + " ✓ " + colorReset + "  " + colorGreen + msg + colorReset
}
func info(msg string) string {
	return colorBgBlue + " i " + colorReset + "  " + msg
}
func warn(msg string) string {
	return colorBgYellow + " ! " + colorReset + "  " + colorYellow + msg + colorReset
}
func fail(msg string) string {
	return colorBgRed + " ✕ " + colorReset + "  " + colorRed + msg + colorReset
}
func done(msg string) string {
	return "  " + colorBgGreen + " DONE " + colorReset + "  " + msg
}
func nextSteps() string {
	return "  " + colorBgBlue + " Next steps " + colorReset
}
func bold(msg string) string { return colorBold + msg + colorReset }
func gray(msg string) string { return colorGray + msg + colorReset }
func badge(bg, label string) string {
	return bg + " " + label + " " + colorReset
}

// ──────────────────────────────────────────────
// buildOutputWriter
// ──────────────────────────────────────────────

// buildOutputWriter processes go compiler output line by line:
//   - Package header lines ("# module/pkg") → gray + dim, no marker.
//   - All other non-empty lines             → red, prefixed with ×.
type buildOutputWriter struct {
	w   *os.File
	buf []byte
}

func newBuildOutputWriter(w *os.File) *buildOutputWriter {
	return &buildOutputWriter{w: w}
}

func (bw *buildOutputWriter) Write(p []byte) (n int, err error) {
	bw.buf = append(bw.buf, p...)
	for {
		nl := bytes.IndexByte(bw.buf, '\n')
		if nl < 0 {
			break
		}
		bw.writeLine(string(bw.buf[:nl]))
		bw.buf = bw.buf[nl+1:]
	}
	return len(p), nil
}

func (bw *buildOutputWriter) writeLine(line string) {
	if strings.TrimSpace(line) == "" {
		fmt.Fprintln(bw.w)
		return
	}
	if strings.HasPrefix(line, "# ") {
		fmt.Fprintf(bw.w, "  %s%s%s\n", colorGray+colorDim, line, colorReset)
		return
	}
	fmt.Fprintf(bw.w, "  %s× %s%s\n", colorRed, line, colorReset)
}

// ──────────────────────────────────────────────
// indentWriter
// ──────────────────────────────────────────────

// indentWriter wraps an io.Writer and prefixes every new line with a fixed
// string so that raw subprocess output (e.g. atlas) is visually aligned with
// grove's own indented output.
type indentWriter struct {
	w     io.Writer
	pfx   []byte
	atSOL bool // true when the next byte starts a new line
}

func newIndentWriter(w io.Writer, prefix string) *indentWriter {
	return &indentWriter{w: w, pfx: []byte(prefix), atSOL: true}
}

func (iw *indentWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		if iw.atSOL {
			if _, err = iw.w.Write(iw.pfx); err != nil {
				return
			}
			iw.atSOL = false
		}
		idx := bytes.IndexByte(p, '\n')
		if idx < 0 {
			var nn int
			nn, err = iw.w.Write(p)
			n += nn
			return
		}
		var nn int
		nn, err = iw.w.Write(p[:idx+1])
		n += nn
		if err != nil {
			return
		}
		p = p[idx+1:]
		iw.atSOL = true
	}
	return
}

// ──────────────────────────────────────────────
// Case conversion
// ──────────────────────────────────────────────

// toPascalCase converts snake_case, kebab-case or camelCase to PascalCase.
//
//	"create_user_request" → "CreateUserRequest"
//	"user-profile"        → "UserProfile"
//	"userProfile"         → "UserProfile"
func toPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		b.WriteString(string(r))
	}
	return b.String()
}

// toSnakeCase converts PascalCase or camelCase to snake_case.
//
//	"UserProfile"  → "user_profile"
//	"HTTPRequest"  → "h_t_t_p_request"  (keep acronyms lowercase as-is)
func toSnakeCase(s string) string {
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// toKebabCase converts PascalCase or camelCase to kebab-case.
//
//	"UserProfile" → "user-profile"
func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

// toLowerFirst returns the string with the first rune lowercased.
func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// ──────────────────────────────────────────────
// Pluralization
// ──────────────────────────────────────────────

// toPlural returns a simple English plural of the given word.
// Handles the most common cases; for irregular forms the user can
// override inside the generated file.
func toPlural(s string) string {
	if s == "" {
		return s
	}

	lower := strings.ToLower(s)

	// Irregular mappings
	irregulars := map[string]string{
		"person": "people",
		"man":    "men",
		"woman":  "women",
		"child":  "children",
		"mouse":  "mice",
		"goose":  "geese",
		"tooth":  "teeth",
		"foot":   "feet",
		"ox":     "oxen",
		"leaf":   "leaves",
		"life":   "lives",
		"knife":  "knives",
		"wife":   "wives",
		"wolf":   "wolves",
		"half":   "halves",
		"self":   "selves",
		"elf":    "elves",
		"loaf":   "loaves",
		"potato": "potatoes",
		"tomato": "tomatoes",
		"cactus": "cacti",
		"focus":  "foci",
		"radius": "radii",
		"datum":  "data",
		"medium": "media",
		"index":  "indices",
		"matrix": "matrices",
		"vertex": "vertices",
		"axis":   "axes",
		"crisis": "crises",
	}
	if plural, ok := irregulars[lower]; ok {
		return plural
	}

	// Already plural (ending in 's') — leave it alone
	if strings.HasSuffix(lower, "ss") {
		return s + "es"
	}
	if strings.HasSuffix(lower, "s") {
		return s
	}

	// -ch, -sh, -x, -z → +es
	if strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") ||
		strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") {
		return s + "es"
	}

	// consonant + y → -ies
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := rune(lower[len(lower)-2])
		if !strings.ContainsRune("aeiou", prev) {
			return s[:len(s)-1] + "ies"
		}
	}

	// -f / -fe → -ves
	if strings.HasSuffix(lower, "fe") {
		return s[:len(s)-2] + "ves"
	}
	if strings.HasSuffix(lower, "f") {
		return s[:len(s)-1] + "ves"
	}

	// -o → +es  (after consonant)
	if strings.HasSuffix(lower, "o") && len(lower) > 1 {
		prev := rune(lower[len(lower)-2])
		if !strings.ContainsRune("aeiou", prev) {
			return s + "es"
		}
	}

	return s + "s"
}

// ──────────────────────────────────────────────
// Go module detection
// ──────────────────────────────────────────────

// getModuleName reads the module name from go.mod in the current directory.
// Falls back to a generic name if the file cannot be read.
func getModuleName() string {
	f, err := os.Open("go.mod")
	if err != nil {
		return "your/module"
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return "your/module"
}

// ──────────────────────────────────────────────
// File system helpers
// ──────────────────────────────────────────────

// fileExists reports whether a regular file exists at path.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ensureDir creates the directory (and any parents) if it does not exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}
