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
// filteredWriter
// ──────────────────────────────────────────────

// filteredWriter wraps an io.Writer and silently drops any line whose trimmed
// content starts with one of the configured prefixes.  Used to suppress
// framework-internal log lines (e.g. "gest: watching for changes…") that
// should not appear in grove's own output.
type filteredWriter struct {
	w       io.Writer
	buf     []byte
	filters []string
}

func newFilteredWriter(w io.Writer, filters ...string) *filteredWriter {
	return &filteredWriter{w: w, filters: filters}
}

func (fw *filteredWriter) Write(p []byte) (n int, err error) {
	fw.buf = append(fw.buf, p...)
	for {
		nl := bytes.IndexByte(fw.buf, '\n')
		if nl < 0 {
			break
		}
		line := string(fw.buf[:nl])
		fw.buf = fw.buf[nl+1:]

		trimmed := strings.TrimSpace(line)
		suppressed := false
		for _, f := range fw.filters {
			if strings.HasPrefix(trimmed, f) {
				suppressed = true
				break
			}
		}
		if !suppressed {
			fmt.Fprintln(fw.w, line)
		}
	}
	return len(p), nil
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
// atlasOutputWriter
// ──────────────────────────────────────────────

// atlasOutputWriter parses the raw text output of Atlas CLI commands and
// re-renders it with Grove's colour palette and badge style.
//
// It recognises the following Atlas output patterns:
//
//	Migrating to version 20240101120000 (3 migrations in total):
//	  -- migrating version 20240101120000
//	    -> CREATE TABLE ...
//	  -- ok (12.3ms)
//	  -- error (12.3ms)
//	  -------------------------
//	  -- 81ms
//	  -- 2 migrations
//	  -- 9 sql statements
//	No migration files to execute
//	The migration directory is synced with the database, no migration files to execute
type atlasOutputWriter struct {
	w   io.Writer
	buf []byte
}

func newAtlasOutputWriter(w io.Writer) *atlasOutputWriter {
	return &atlasOutputWriter{w: w}
}

func (aw *atlasOutputWriter) Write(p []byte) (n int, err error) {
	aw.buf = append(aw.buf, p...)
	for {
		nl := bytes.IndexByte(aw.buf, '\n')
		if nl < 0 {
			break
		}
		aw.writeLine(string(aw.buf[:nl]))
		aw.buf = aw.buf[nl+1:]
	}
	return len(p), nil
}

func (aw *atlasOutputWriter) Flush() {
	if len(aw.buf) > 0 {
		aw.writeLine(string(aw.buf))
		aw.buf = nil
	}
}

func (aw *atlasOutputWriter) writeLine(line string) {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return
	}

	// ── "Migrating to version …" header ──────────────────────────────────────
	if strings.HasPrefix(trimmed, "Migrating to version ") {
		fmt.Fprintf(aw.w, "\n  %s%s%s\n", colorBold, trimmed, colorReset)
		return
	}

	// ── "-- migrating version …" section opener ───────────────────────────────
	if strings.HasPrefix(trimmed, "-- migrating version ") {
		version := strings.TrimPrefix(trimmed, "-- migrating version ")
		fmt.Fprintf(aw.w,
			"\n  %s%s%s  %s%s%s\n",
			colorBgBlue, " MIGRATE ", colorReset,
			colorBold, version, colorReset,
		)
		return
	}

	// ── "-- ok (…)" success line ──────────────────────────────────────────────
	if strings.HasPrefix(trimmed, "-- ok") {
		elapsed := ""
		if start := strings.Index(trimmed, "("); start >= 0 {
			if end := strings.Index(trimmed, ")"); end > start {
				elapsed = trimmed[start+1 : end]
			}
		}
		if elapsed != "" {
			fmt.Fprintf(aw.w,
				"  %s%s%s  %s%s%s\n",
				colorBgGreen, " OK ", colorReset,
				colorGray, elapsed, colorReset,
			)
		} else {
			fmt.Fprintf(aw.w, "  %s%s%s\n", colorBgGreen, " OK ", colorReset)
		}
		return
	}

	// ── "-- error (…)" failure line ───────────────────────────────────────────
	if strings.HasPrefix(trimmed, "-- error") {
		elapsed := ""
		if start := strings.Index(trimmed, "("); start >= 0 {
			if end := strings.Index(trimmed, ")"); end > start {
				elapsed = trimmed[start+1 : end]
			}
		}
		if elapsed != "" {
			fmt.Fprintf(aw.w,
				"  %s%s%s  %s%s%s\n",
				colorBgRed, " ERROR ", colorReset,
				colorGray, elapsed, colorReset,
			)
		} else {
			fmt.Fprintf(aw.w, "  %s%s%s\n", colorBgRed, " ERROR ", colorReset)
		}
		return
	}

	// ── SQL statement lines ("-> …") ─────────────────────────────────────────
	if strings.HasPrefix(trimmed, "-> ") {
		sql := strings.TrimPrefix(trimmed, "-> ")
		// Upper-case SQL keywords get a cyan tint; the rest stays plain.
		keyword, rest := splitSQLKeyword(sql)
		if keyword != "" {
			fmt.Fprintf(aw.w,
				"    %s%s%s %s\n",
				colorCyan, keyword, colorReset, rest,
			)
		} else {
			fmt.Fprintf(aw.w, "    %s%s%s\n", colorGray, sql, colorReset)
		}
		return
	}

	// ── Continuation lines of a multi-line SQL statement ─────────────────────
	// Atlas indents them with spaces but no "-> " prefix after the first line.
	// They only appear inside a migration block (after a "-> "), so we just
	// render them indented and dimmed.
	if strings.HasPrefix(line, "      ") || strings.HasPrefix(line, "\t") {
		fmt.Fprintf(
			aw.w,
			"    %s%s%s\n",
			colorGray+colorDim,
			trimmed,
			colorReset,
		)
		return
	}

	// ── Summary separator "---…" ──────────────────────────────────────────────
	if strings.HasPrefix(trimmed, "---") {
		fmt.Fprintf(
			aw.w,
			"\n  %s%s%s\n",
			colorGray+colorDim,
			strings.Repeat("─", 40),
			colorReset,
		)
		return
	}

	// ── Summary stat lines ("-- 81ms", "-- 2 migrations", "-- 9 sql …") ──────
	if strings.HasPrefix(trimmed, "-- ") {
		stat := strings.TrimPrefix(trimmed, "-- ")
		fmt.Fprintf(aw.w, "  %s%s%s\n", colorGray, stat, colorReset)
		return
	}

	// ── "No migration files to execute" / already synced ─────────────────────
	lower := strings.ToLower(trimmed)
	if strings.Contains(lower, "no migration") ||
		strings.Contains(lower, "synced with the database") ||
		strings.Contains(lower, "nothing to do") {
		fmt.Fprintf(aw.w,
			"\n  %s%s%s  %s%s%s\n\n",
			colorBgGreen, " UP TO DATE ", colorReset,
			colorGray, trimmed, colorReset,
		)
		return
	}

	// ── Fallback — indent and dim ─────────────────────────────────────────────
	fmt.Fprintf(aw.w, "  %s%s%s\n", colorGray, trimmed, colorReset)
}

// splitSQLKeyword splits a SQL statement string into the first keyword and the
// remainder. Returns ("", "") if the string doesn't start with a known keyword.
func splitSQLKeyword(s string) (keyword, rest string) {
	keywords := []string{
		"CREATE UNIQUE INDEX",
		"CREATE INDEX",
		"CREATE EXTENSION",
		"CREATE TABLE",
		"CREATE SEQUENCE",
		"CREATE TYPE",
		"CREATE VIEW",
		"CREATE FUNCTION",
		"CREATE TRIGGER",
		"ALTER TABLE",
		"ALTER COLUMN",
		"ALTER SEQUENCE",
		"DROP TABLE",
		"DROP INDEX",
		"DROP COLUMN",
		"DROP CONSTRAINT",
		"DROP SEQUENCE",
		"DROP TYPE",
		"DROP VIEW",
		"INSERT INTO",
		"UPDATE",
		"DELETE FROM",
		"SELECT",
		"COMMENT ON",
	}
	upper := strings.ToUpper(s)
	for _, kw := range keywords {
		if strings.HasPrefix(upper, kw) {
			return kw, s[len(kw):]
		}
	}
	return "", ""
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

// toWords converts a PascalCase or snake_case name to space-separated words
// with each word title-cased. Consecutive uppercase sequences are kept as one
// word so acronyms are preserved. Used for gest.Describe() labels.
//
//	"BookCool"        → "Book Cool"
//	"AuthService"     → "Auth Service"
//	"order_items"     → "Order Items"
//	"ISBN"            → "ISBN"
//	"parseJSON"       → "Parse JSON"
//	"UserAuthService" → "User Auth Service"
func toWords(s string) string {
	// First split on underscores (handles snake_case input).
	// Then split each segment on PascalCase boundaries, keeping consecutive
	// uppercase runs (acronyms) together.
	var words []string
	for _, seg := range strings.Split(s, "_") {
		if seg == "" {
			continue
		}
		runes := []rune(seg)
		start := 0
		for i := 1; i < len(runes); i++ {
			curr := runes[i]
			prev := runes[i-1]
			// Start a new word when transitioning from lower to upper,
			// or from a run of uppers into a lower (e.g. "JSONParser" → "JSON", "Parser").
			splitHere := false
			if unicode.IsUpper(curr) && unicode.IsLower(prev) {
				splitHere = true
			} else if unicode.IsLower(curr) && unicode.IsUpper(prev) && i-start > 1 {
				// back up one so the last upper joins the new word
				words = append(words, string(runes[start:i-1]))
				start = i - 1
				continue
			}
			if splitHere {
				words = append(words, string(runes[start:i]))
				start = i
			}
		}
		words = append(words, string(runes[start:]))
	}
	// Title-case each word.
	for i, w := range words {
		if w == "" {
			continue
		}
		r := []rune(w)
		// If the entire word is uppercase (acronym), leave it as-is.
		allUpper := true
		for _, c := range r {
			if unicode.IsLower(c) {
				allUpper = false
				break
			}
		}
		if !allUpper {
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
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

// toSingular returns the singular form of an English word.
// It is the inverse of toPlural and handles the same irregular and suffix
// rules so that inputs like "books" → "book", "categories" → "category".
// If the word is already singular it is returned unchanged.
func toSingular(s string) string {
	if s == "" {
		return s
	}

	lower := strings.ToLower(s)

	// Irregular plurals — map plural → singular.
	irregulars := map[string]string{
		"people":   "person",
		"men":      "man",
		"women":    "woman",
		"children": "child",
		"mice":     "mouse",
		"geese":    "goose",
		"teeth":    "tooth",
		"feet":     "foot",
		"oxen":     "ox",
		"leaves":   "leaf",
		"lives":    "life",
		"knives":   "knife",
		"wives":    "wife",
		"wolves":   "wolf",
		"halves":   "half",
		"selves":   "self",
		"elves":    "elf",
		"loaves":   "loaf",
		"potatoes": "potato",
		"tomatoes": "tomato",
		"cacti":    "cactus",
		"foci":     "focus",
		"radii":    "radius",
		"data":     "datum",
		"media":    "medium",
		"indices":  "index",
		"matrices": "matrix",
		"vertices": "vertex",
		"axes":     "axis",
		"crises":   "crisis",
	}
	if singular, ok := irregulars[lower]; ok {
		// Preserve original casing of the first letter.
		if len(s) > 0 && unicode.IsUpper(rune(s[0])) {
			return strings.ToUpper(singular[:1]) + singular[1:]
		}
		return singular
	}

	// Already singular — doesn't end in s at all.
	if !strings.HasSuffix(lower, "s") {
		return s
	}

	// -ies → -y  (categories → category)
	if strings.HasSuffix(lower, "ies") && len(lower) > 3 {
		return s[:len(s)-3] + "y"
	}

	// -ves → -f or -fe
	if strings.HasSuffix(lower, "ves") {
		// Most common: -ves → -f (wolves → wolf)
		return s[:len(s)-3] + "f"
	}

	// -ses, -xes, -zes, -ches, -shes → strip -es
	if strings.HasSuffix(lower, "ses") ||
		strings.HasSuffix(lower, "xes") ||
		strings.HasSuffix(lower, "zes") ||
		strings.HasSuffix(lower, "ches") ||
		strings.HasSuffix(lower, "shes") {
		return s[:len(s)-2]
	}

	// -oes → -o  (potatoes handled above via irregulars; catches others)
	if strings.HasSuffix(lower, "oes") {
		return s[:len(s)-2]
	}

	// -ss words are already singular (class, grass) — leave alone.
	if strings.HasSuffix(lower, "ss") {
		return s
	}

	// Default: strip trailing -s.
	return s[:len(s)-1]
}

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
