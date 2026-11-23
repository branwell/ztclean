This is a small command-line tool to clean up zoom transcripts.

## Running the tool.

Command: `ztclean -f <filename>`
Optional specification of output file: `ztclean -f <filename> -o foo.md` otherwise default is output file is the same as the filename, replace .vtt with .md
(Standard Go flag syntax, so `-f=<filename>` also works.)

## Cleanup Rules

The "cleaning" process will do the following:
- Remove WEBVTT
- Translate all unicode characters to ASCII.
- Remove all the numeric headers (1, 2, 3, 4, 5...)
- remove native timestamps and replace them with hh:mm (see below)
- Output should be a markdown formatted file, so replace .vtt with .md
- If there was only 1 speaker
  - remove the speaker
  - put a hh:mm timestamp every 1 minute
  - compress the paragraphs roughly between timestamps
- If there was more than one speaker
  - wrap the speaker including the colon inside backticks, e.g. `` `Fozzie Bear:` ``
  - compress the sections when the speaker switches
  - put a hh:mm timestamp between the speaker switching
  - if one speaker talks for a long stretch, also put a hh:mm timestamp around every 5 minutes

## Examples

### Sample with multiple speakers

```md
WEBVTT

1
01:19:02.330 --> 01:19:09.190
Kermit the Frog: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

2
01:19:09.190 --> 01:19:59.500
Kermit the Frog: Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

3
01:20:10.500 --> 01:20:15.039
Fozzie Bear: Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

### Sample output 

```md
`Kermit the Frog:` Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

**01:20**

`Fozzie Bear:` Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

### sample with one speaker

```md
WEBVTT

1
01:19:02.330 --> 01:19:09.190
Kermit the Frog: Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

2
01:19:09.190 --> 01:19:59.500
Kermit the Frog: Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

3
01:20:10.500 --> 01:20:15.039
Kermit the Frog: Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

### Sample output (ignore speaker, ignore timestamps)

```md
Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

**01:20**

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
```

## Design decisions

- **Speaker label:** wrapped in backticks (`` `Name:` ``), matching the sample output.
- **Unicode to ASCII:** a hand-written lookup table (smart quotes, dashes, ellipsis, special spaces, common symbols, accented Latin letters). Any other non-ASCII character is dropped. No external library.
- **Command line:** Go standard `flag` package; `-f` is required, `-o` is optional.
- **Timestamps:** formatted as `**hh:mm**` on their own line, from the start time of the cue that begins the new paragraph. No timestamp before the first paragraph.
  - One speaker: a new paragraph and timestamp whenever the cue start minute changes.
  - Several speakers: a new paragraph and timestamp on each speaker change, and also at the first cue starting 5 or more minutes after the current paragraph began. The speaker label is not repeated on that continuation paragraph.
- **Speaker detection:** text before the first `": "` in a cue. A cue with no speaker prefix continues the previous speaker. Multi-line cue text is joined with spaces.
- **Existing output file:** overwritten.
- **Output extension:** an `-o` name ending in `.md` (any case) is used as is; otherwise `.md` is appended (`notes` becomes `notes.md`, `notes.txt` becomes `notes.txt.md`). Without `-o`, the input's extension is replaced with `.md`.
- **Input protection:** if the output path is the same as the input path (for example, input `notes.md` with no `-o`), the tool exits with an error instead of overwriting the input.
