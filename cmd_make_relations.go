package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type relationEdit struct {
	StructName   string
	BelongsTo    []string
	HasMany      []string
	HasBelongsTo map[string]bool
	HasHasMany   map[string]bool
}

type relationLog struct {
	SourceModel string
	TargetModel string
	FKField     string
	FKType      string
	Kind        string // "belongs-to" | "has-many"
	Line        string
}

var (
	makeRelationsPath          string
	makeRelationsDryRun        bool
	makeRelationsVerbose       bool
	makeRelationsWithBelongsTo bool
	makeRelationsModels        []string
)

var makeRelationsCmd = &cobra.Command{
	Use:   "make:relations",
	Short: "Infer and add GORM relation fields between models",
	Long: bold(
		"make:relations",
	) + ` scans model structs and infers relationships from foreign key fields.

By default it generates only has-many fields on the target model:
  ` + colorGray + `PaymentMethods []PaymentMethod ` + "`gorm:\"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\" json:\"payment_methods,omitempty\"`" + colorReset + `

If you pass ` + colorGreen + `--with-belongs-to` + colorReset + `, it also generates belongs-to on the source model:
  ` + colorGray + `User *User ` + "`gorm:\"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\" json:\"user,omitempty\"`" + colorReset + `

It avoids duplicate fields and supports ` + colorGreen + `--dry-run` + colorReset + ` and ` + colorGreen + `--verbose` + colorReset + `.`,
	Args: cobra.NoArgs,
	RunE: runMakeRelations,
}

func init() {
	makeRelationsCmd.Flags().StringVar(
		&makeRelationsPath,
		"path",
		"internal/models",
		"Path to models directory",
	)
	makeRelationsCmd.Flags().BoolVar(
		&makeRelationsDryRun,
		"dry-run",
		false,
		"Preview inferred relations without writing files",
	)
	makeRelationsCmd.Flags().BoolVar(
		&makeRelationsVerbose,
		"verbose",
		false,
		"Show detailed relation inference logs",
	)
	makeRelationsCmd.Flags().BoolVar(
		&makeRelationsWithBelongsTo,
		"with-belongs-to",
		false,
		"Also generate belongs-to fields on source models",
	)
	makeRelationsCmd.Flags().StringSliceVar(
		&makeRelationsModels,
		"model",
		nil,
		"Only process specific source model(s), e.g. --model PaymentMethod or --model PaymentMethod,Order",
	)
}

type modelFile struct {
	Path       string
	StructName string
	Source     string
	FieldNames map[string]bool
	FKFields   []fkField
}

type fkField struct {
	Name       string
	TargetName string
	TypeName   string
}

func runMakeRelations(_ *cobra.Command, _ []string) error {
	fmt.Println()
	fmt.Printf(
		"  %sScanning models%s in %s\n",
		colorGray,
		colorReset,
		colorCyan+makeRelationsPath+colorReset,
	)
	fmt.Println()

	models, err := loadModelFiles(makeRelationsPath)
	if err != nil {
		return err
	}
	if len(models) == 0 {
		return fmt.Errorf("no model structs found in %s", makeRelationsPath)
	}

	modelByName := map[string]*modelFile{}
	for _, m := range models {
		modelByName[m.StructName] = m
	}

	selectedSources := map[string]bool{}
	if len(makeRelationsModels) > 0 {
		for _, raw := range makeRelationsModels {
			for _, part := range strings.Split(raw, ",") {
				name := strings.TrimSpace(part)
				if name == "" {
					continue
				}
				selectedSources[name] = true
			}
		}
		if len(selectedSources) == 0 {
			return fmt.Errorf(
				"--model was provided but no valid model names were parsed",
			)
		}

		var missing []string
		for name := range selectedSources {
			if _, ok := modelByName[name]; !ok {
				missing = append(missing, name)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			return fmt.Errorf(
				"model(s) not found: %s",
				strings.Join(missing, ", "),
			)
		}
	}

	editsByFile := map[string]*relationEdit{}
	getEdit := func(structName string, filePath string) *relationEdit {
		e, ok := editsByFile[filePath]
		if !ok {
			e = &relationEdit{
				StructName:   structName,
				BelongsTo:    []string{},
				HasMany:      []string{},
				HasBelongsTo: map[string]bool{},
				HasHasMany:   map[string]bool{},
			}
			editsByFile[filePath] = e
		}
		return e
	}

	logs := []relationLog{}

	for _, source := range models {
		if len(selectedSources) > 0 && !selectedSources[source.StructName] {
			continue
		}
		for _, fk := range source.FKFields {
			target, ok := modelByName[fk.TargetName]
			if !ok || target.StructName == source.StructName {
				continue
			}

			if makeRelationsWithBelongsTo {
				belongsFieldName := fk.TargetName
				if !source.FieldNames[belongsFieldName] {
					edit := getEdit(source.StructName, source.Path)
					if !edit.HasBelongsTo[belongsFieldName] {
						line := fmt.Sprintf(
							"\t%s *%s `gorm:\"foreignKey:%s;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\" json:\"%s,omitempty\"`",
							belongsFieldName,
							fk.TargetName,
							fk.Name,
							toSnakeCase(fk.TargetName),
						)
						edit.BelongsTo = append(edit.BelongsTo, line)
						edit.HasBelongsTo[belongsFieldName] = true

						logs = append(logs, relationLog{
							SourceModel: source.StructName,
							TargetModel: fk.TargetName,
							FKField:     fk.Name,
							FKType:      fk.TypeName,
							Kind:        "belongs-to",
							Line:        line,
						})
					}
				}
			}

			hasManyFieldName := betterPluralize(source.StructName)
			if !target.FieldNames[hasManyFieldName] {
				edit := getEdit(target.StructName, target.Path)
				if !edit.HasHasMany[hasManyFieldName] {
					line := fmt.Sprintf(
						"\t%s []%s `gorm:\"foreignKey:%s;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;\" json:\"%s,omitempty\"`",
						hasManyFieldName,
						source.StructName,
						fk.Name,
						toPlural(toSnakeCase(source.StructName)),
					)
					edit.HasMany = append(edit.HasMany, line)
					edit.HasHasMany[hasManyFieldName] = true

					logs = append(logs, relationLog{
						SourceModel: source.StructName,
						TargetModel: target.StructName,
						FKField:     fk.Name,
						FKType:      fk.TypeName,
						Kind:        "has-many",
						Line:        line,
					})
				}
			}
		}
	}

	if len(editsByFile) == 0 {
		fmt.Println(
			success("No relation fields to add. Everything looks up-to-date."),
		)
		fmt.Println()
		return nil
	}

	paths := make([]string, 0, len(editsByFile))
	for p := range editsByFile {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	totalAdded := 0
	for _, p := range paths {
		model := findModelByPath(models, p)
		if model == nil {
			continue
		}

		edit := editsByFile[p]
		newSrc, added, err := applyRelationEdit(model.Source, edit)
		if err != nil {
			return fmt.Errorf("failed updating %s: %w", p, err)
		}
		if added == 0 {
			continue
		}
		totalAdded += added

		if makeRelationsDryRun {
			fmt.Println(
				info(
					fmt.Sprintf("[dry-run] %s (+%d relation fields)", p, added),
				),
			)
			continue
		}

		if err := os.WriteFile(p, []byte(newSrc), 0o644); err != nil {
			return fmt.Errorf("failed writing %s: %w", p, err)
		}
		fmt.Println(
			success(fmt.Sprintf("Updated %s (+%d relation fields)", p, added)),
		)
	}

	if makeRelationsVerbose && len(logs) > 0 {
		fmt.Println()
		fmt.Println(
			"  " + badge(colorBgBlue, " VERBOSE ") + " inferred relations",
		)
		for _, lg := range logs {
			fmt.Printf(
				"    %s -> %s  %s  (%s %s)\n",
				colorCyan+lg.SourceModel+colorReset,
				colorCyan+lg.TargetModel+colorReset,
				colorGray+lg.FKField+colorReset,
				colorGray+lg.FKType+colorReset,
				colorGray+lg.Kind+colorReset,
			)
		}
	}

	if totalAdded == 0 {
		fmt.Println(
			success("No relation fields to add. Everything looks up-to-date."),
		)
		fmt.Println()
		return nil
	}

	fmt.Println()
	if makeRelationsDryRun {
		fmt.Println(
			done(
				fmt.Sprintf(
					"Dry-run complete. %d relation fields inferred.",
					totalAdded,
				),
			),
		)
	} else {
		fmt.Println(done(fmt.Sprintf("Added %d relation fields.", totalAdded)))
	}
	fmt.Println()

	return nil
}

func loadModelFiles(modelsDir string) ([]*modelFile, error) {
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf(
			"could not read models directory %s: %w",
			modelsDir,
			err,
		)
	}

	var out []*modelFile

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, "_test.go") {
			continue
		}

		fullPath := filepath.Join(modelsDir, name)
		src, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("could not read %s: %w", fullPath, err)
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, fullPath, src, parser.ParseComments)
		if err != nil {
			continue
		}

		for _, st := range extractStructs(file) {
			fields := map[string]bool{}
			fks := []fkField{}

			for _, sf := range st.Fields {
				fields[sf.Name] = true

				if !strings.HasSuffix(sf.Name, "ID") || len(sf.Name) <= 2 {
					continue
				}
				target := strings.TrimSuffix(sf.Name, "ID")
				if target == "" {
					continue
				}
				if !isSupportedFKType(sf.TypeName) {
					continue
				}

				fks = append(fks, fkField{
					Name:       sf.Name,
					TargetName: target,
					TypeName:   sf.TypeName,
				})
			}

			out = append(out, &modelFile{
				Path:       fullPath,
				StructName: st.Name,
				Source:     string(src),
				FieldNames: fields,
				FKFields:   fks,
			})
		}
	}

	return out, nil
}

type structField struct {
	Name     string
	TypeName string
}

type structMeta struct {
	Name   string
	Fields []structField
}

func extractStructs(file *ast.File) []structMeta {
	out := []structMeta{}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			fields := []structField{}
			for _, f := range st.Fields.List {
				typeName := exprToTypeName(f.Type)
				for _, n := range f.Names {
					fields = append(fields, structField{
						Name:     n.Name,
						TypeName: typeName,
					})
				}
			}

			out = append(out, structMeta{
				Name:   ts.Name.Name,
				Fields: fields,
			})
		}
	}

	return out
}

func exprToTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		left := exprToTypeName(t.X)
		if left == "" {
			return t.Sel.Name
		}
		return left + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToTypeName(t.X)
	case *ast.ArrayType:
		return "[]" + exprToTypeName(t.Elt)
	default:
		return ""
	}
}

func isSupportedFKType(typeName string) bool {
	normalized := strings.TrimPrefix(typeName, "*")
	switch normalized {
	case "string",
		"int",
		"int32",
		"int64",
		"uint",
		"uint32",
		"uint64",
		"uuid.UUID":
		return true
	default:
		return false
	}
}

func betterPluralize(s string) string {
	if s == "" {
		return s
	}

	parts := splitPascalWords(s)
	if len(parts) == 0 {
		return toPlural(s)
	}
	last := parts[len(parts)-1]
	parts[len(parts)-1] = toPlural(last)
	return strings.Join(parts, "")
}

func splitPascalWords(s string) []string {
	if s == "" {
		return nil
	}

	var parts []string
	runes := []rune(s)
	start := 0

	for i := 1; i < len(runes); i++ {
		prev := runes[i-1]
		cur := runes[i]
		nextLower := i+1 < len(runes) && isLower(runes[i+1])

		if isUpper(cur) && (isLower(prev) || (isUpper(prev) && nextLower)) {
			parts = append(parts, string(runes[start:i]))
			start = i
		}
	}
	parts = append(parts, string(runes[start:]))
	return parts
}

func isUpper(r rune) bool { return r >= 'A' && r <= 'Z' }
func isLower(r rune) bool { return r >= 'a' && r <= 'z' }

func findModelByPath(models []*modelFile, path string) *modelFile {
	for _, m := range models {
		if m.Path == path {
			return m
		}
	}
	return nil
}

func applyRelationEdit(src string, edit *relationEdit) (string, int, error) {
	structNeedle := "type " + edit.StructName + " struct {"
	start := strings.Index(src, structNeedle)
	if start < 0 {
		return src, 0, nil
	}

	openBrace := strings.Index(src[start:], "{")
	if openBrace < 0 {
		return src, 0, nil
	}
	openPos := start + openBrace

	closePos, err := findMatchingBrace(src, openPos)
	if err != nil {
		return src, 0, err
	}

	body := src[openPos+1 : closePos]
	insertLines := []string{}
	insertLines = append(insertLines, edit.BelongsTo...)
	insertLines = append(insertLines, edit.HasMany...)
	if len(insertLines) == 0 {
		return src, 0, nil
	}

	added := len(insertLines)
	var b strings.Builder
	b.WriteString(src[:openPos+1])

	trimBodyRight := strings.TrimRight(body, " \t\r\n")
	if strings.TrimSpace(trimBodyRight) != "" {
		b.WriteString(trimBodyRight)
		if !strings.HasSuffix(trimBodyRight, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("\n")
	}

	for _, ln := range insertLines {
		b.WriteString(ln)
		b.WriteString("\n")
	}
	b.WriteString(src[closePos:])

	return b.String(), added, nil
}

func findMatchingBrace(s string, openPos int) (int, error) {
	depth := 0
	for i := openPos; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return -1, fmt.Errorf("unmatched struct braces")
}
