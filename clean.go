package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// longTurn is how long one speaker may talk in a multi-speaker
// transcript before a timestamp is inserted mid-turn.
const longTurn = 5 * time.Minute

type cue struct {
	start   time.Duration
	speaker string
	text    string
}

type paragraph struct {
	start   time.Duration
	speaker string
	label   bool // show the speaker label
	stamp   bool // precede with an hh:mm timestamp
	text    []string
}

// Clean converts the contents of a Zoom .vtt transcript into Markdown.
// If stamps is false, the hh:mm timestamp lines are left out; paragraph
// breaks are unchanged.
func Clean(vtt string, stamps bool) string {
	cues := parseCues(toASCII(vtt))
	paras := group(cues)
	return render(paras, stamps)
}

// parseCues extracts the cues from a WebVTT document, dropping the header,
// numeric cue identifiers and end times.
func parseCues(vtt string) []cue {
	vtt = strings.ReplaceAll(vtt, "\r\n", "\n")
	var cues []cue
	var cur *cue
	for _, line := range strings.Split(vtt, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			cur = nil
			continue
		}
		if start, ok := parseTiming(line); ok {
			cues = append(cues, cue{start: start})
			cur = &cues[len(cues)-1]
			continue
		}
		if cur == nil {
			// Header (WEBVTT), cue identifiers, or stray lines outside a cue.
			continue
		}
		if cur.text == "" {
			cur.speaker, cur.text = splitSpeaker(line)
		} else {
			cur.text += " " + line
		}
	}
	// Drop empty cues and let unlabeled cues continue the previous speaker.
	out := cues[:0]
	prev := ""
	for _, c := range cues {
		if c.text == "" {
			continue
		}
		if c.speaker == "" {
			c.speaker = prev
		}
		prev = c.speaker
		out = append(out, c)
	}
	return out
}

// parseTiming parses a "start --> end" line and returns the start time.
func parseTiming(line string) (time.Duration, bool) {
	start, _, ok := strings.Cut(line, "-->")
	if !ok {
		return 0, false
	}
	return parseTimestamp(strings.TrimSpace(start))
}

// parseTimestamp parses "hh:mm:ss.mmm" or "mm:ss.mmm".
func parseTimestamp(s string) (time.Duration, bool) {
	parts := strings.Split(s, ":")
	if len(parts) == 2 {
		parts = append([]string{"0"}, parts...)
	}
	if len(parts) != 3 {
		return 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	sec, err3 := strconv.ParseFloat(parts[2], 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, false
	}
	d := time.Duration(h)*time.Hour + time.Duration(m)*time.Minute +
		time.Duration(sec*float64(time.Second))
	return d, true
}

// splitSpeaker splits "Name: text" into its speaker and text.
func splitSpeaker(line string) (speaker, text string) {
	name, rest, ok := strings.Cut(line, ": ")
	if !ok || strings.TrimSpace(name) == "" {
		return "", line
	}
	return strings.TrimSpace(name), strings.TrimSpace(rest)
}

// group merges cues into paragraphs. With one speaker, a new paragraph
// starts each minute. With several, a new paragraph starts on each
// speaker change, or when one speaker has talked for longTurn.
func group(cues []cue) []paragraph {
	speakers := map[string]bool{}
	for _, c := range cues {
		if c.speaker != "" {
			speakers[c.speaker] = true
		}
	}
	multi := len(speakers) > 1

	var paras []paragraph
	for _, c := range cues {
		n := len(paras)
		switch {
		case n == 0:
			paras = append(paras, paragraph{start: c.start, speaker: c.speaker, label: multi})
		case multi && c.speaker != paras[n-1].speaker:
			paras = append(paras, paragraph{start: c.start, speaker: c.speaker, label: true, stamp: true})
		case multi && c.start-paras[n-1].start >= longTurn:
			paras = append(paras, paragraph{start: c.start, speaker: c.speaker, stamp: true})
		case !multi && minute(c.start) != minute(paras[n-1].start):
			paras = append(paras, paragraph{start: c.start, stamp: true})
		}
		p := &paras[len(paras)-1]
		p.text = append(p.text, c.text)
	}
	return paras
}

func minute(d time.Duration) int64 { return int64(d / time.Minute) }

// hhmm formats a duration as hh:mm.
func hhmm(d time.Duration) string {
	m := minute(d)
	return fmt.Sprintf("%02d:%02d", m/60, m%60)
}

func render(paras []paragraph, stamps bool) string {
	var blocks []string
	for _, p := range paras {
		if stamps && p.stamp {
			blocks = append(blocks, "**"+hhmm(p.start)+"**")
		}
		body := strings.Join(p.text, " ")
		if p.label {
			body = "`" + p.speaker + ":` " + body
		}
		blocks = append(blocks, body)
	}
	if len(blocks) == 0 {
		return ""
	}
	return strings.Join(blocks, "\n\n") + "\n"
}
