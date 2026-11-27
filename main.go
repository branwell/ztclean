// Command ztclean converts a Zoom .vtt transcript into a compact Markdown file.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	in := flag.String("f", "", "input .vtt transcript (required)")
	out := flag.String("o", "", "output .md file (default: input name with .md extension)")
	noStamps := flag.Bool("t", false, "omit hh:mm timestamps (paragraph breaks are kept)")
	flag.Parse()

	if *in == "" || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*in, *out, !*noStamps); err != nil {
		fmt.Fprintln(os.Stderr, "ztclean:", err)
		os.Exit(1)
	}
}

func run(in, out string, stamps bool) error {
	out = outputPath(in, out)
	if filepath.Clean(out) == filepath.Clean(in) {
		return fmt.Errorf("output file %q would overwrite the input", out)
	}
	data, err := os.ReadFile(in)
	if err != nil {
		return fmt.Errorf("read transcript: %w", err)
	}
	if err := os.WriteFile(out, []byte(Clean(string(data), stamps)), 0o644); err != nil {
		return fmt.Errorf("write markdown: %w", err)
	}
	return nil
}

// outputPath returns the file to write. With no -o, it is the input with
// its extension (normally .vtt) replaced by .md. A given -o keeps a .md
// extension (any case) as is; otherwise .md is appended.
func outputPath(in, out string) string {
	if out == "" {
		return strings.TrimSuffix(in, filepath.Ext(in)) + ".md"
	}
	if strings.EqualFold(filepath.Ext(out), ".md") {
		return out
	}
	return out + ".md"
}
