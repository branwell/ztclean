package main

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	lorem1 = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
	lorem2 = "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."
	lorem3 = "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
)

func TestCleanMultipleSpeakers(t *testing.T) {
	in := `WEBVTT

1
01:19:02.330 --> 01:19:09.190
Kermit the Frog: ` + lorem1 + `

2
01:19:09.190 --> 01:19:59.500
Kermit the Frog: ` + lorem2 + `

3
01:20:10.500 --> 01:20:15.039
Fozzie Bear: ` + lorem3 + `
`
	want := "`Kermit the Frog:` " + lorem1 + " " + lorem2 + "\n\n**01:20**\n\n`Fozzie Bear:` " + lorem3 + "\n"
	if got := Clean(in); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCleanSingleSpeaker(t *testing.T) {
	in := `WEBVTT

1
01:19:02.330 --> 01:19:09.190
Kermit the Frog: ` + lorem1 + `

2
01:19:09.190 --> 01:19:59.500
Kermit the Frog: ` + lorem2 + `

3
01:20:10.500 --> 01:20:15.039
Kermit the Frog: ` + lorem3 + `
`
	want := lorem1 + " " + lorem2 + "\n\n**01:20**\n\n" + lorem3 + "\n"
	if got := Clean(in); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCleanLongTurn(t *testing.T) {
	in := "WEBVTT\r\n\r\n" +
		"1\r\n00:00:01.000 --> 00:00:05.000\r\nA: one\r\n\r\n" +
		"2\r\n00:04:59.000 --> 00:05:01.000\r\nA: two\r\n\r\n" +
		"3\r\n00:05:01.000 --> 00:05:03.000\r\nA: three\r\n\r\n" +
		"4\r\n00:06:00.000 --> 00:06:03.000\r\nB: four\r\n"
	want := "`A:` one two\n\n**00:05**\n\nthree\n\n**00:06**\n\n`B:` four\n"
	if got := Clean(in); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestCleanMultilineAndUnlabeledCue(t *testing.T) {
	in := "WEBVTT\n\n1\n00:00:01.000 --> 00:00:02.000\nA: first\nline\n\n" +
		"2\n00:00:02.000 --> 00:00:03.000\nstill A\n\n" +
		"3\n00:00:03.000 --> 00:00:04.000\nB: hi\n"
	want := "`A:` first line still A\n\n**00:00**\n\n`B:` hi\n"
	if got := Clean(in); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestToASCII(t *testing.T) {
	got := toASCII("“Café” — it’s naïve… ß 日本")
	want := `"Cafe" -- it's naive... ss `
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestOutputPath(t *testing.T) {
	cases := []struct{ in, out, want string }{
		{"meeting.vtt", "", "meeting.md"},
		{"dir/talk.2024.vtt", "", "dir/talk.2024.md"},
		{"noext", "", "noext.md"},
		{"meeting.vtt", "notes.md", "notes.md"},
		{"meeting.vtt", "notes.MD", "notes.MD"},
		{"meeting.vtt", "notes", "notes.md"},
		{"meeting.vtt", "notes.txt", "notes.txt.md"},
		{"meeting.vtt", "q3.review", "q3.review.md"},
	}
	for _, c := range cases {
		if got := outputPath(c.in, c.out); got != c.want {
			t.Errorf("outputPath(%q, %q) = %q, want %q", c.in, c.out, got, c.want)
		}
	}
}

func TestRunRefusesToOverwriteInput(t *testing.T) {
	in := filepath.Join(t.TempDir(), "notes.md")
	if err := os.WriteFile(in, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(in, ""); err == nil {
		t.Error("expected an error when output equals input")
	}
	if got, _ := os.ReadFile(in); string(got) != "original" {
		t.Errorf("input was modified: %q", got)
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "t.vtt")
	if err := os.WriteFile(in, []byte("WEBVTT\n\n1\n00:00:01.000 --> 00:00:02.000\nA: hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := run(in, ""); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dir, "t.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hi\n" {
		t.Errorf("got %q", got)
	}
}
