package main

import "strings"

// asciiMap holds ASCII replacements for common non-ASCII characters.
// Runes not listed here (and not ASCII) are dropped.
var asciiMap = map[rune]string{
	// Quotes and apostrophes
	'‘': "'", '’': "'", '‚': "'", '‛': "'", '′': "'",
	'“': `"`, '”': `"`, '„': `"`, '‟': `"`, '″': `"`,
	'«': `"`, '»': `"`, '‹': "'", '›': "'",
	// Dashes and hyphens
	'‐': "-", '‑': "-", '‒': "-", '–': "-", '—': "--",
	'―': "--", '−': "-",
	// Spaces
	' ': " ", ' ': " ", ' ': " ", ' ': " ", ' ': " ",
	// Punctuation and symbols
	'…': "...", '•': "*", '·': "*", '×': "x", '÷': "/",
	'©': "(c)", '®': "(R)", '™': "(TM)", '°': " degrees",
	'€': "EUR", '£': "GBP", '¥': "JPY", '½': "1/2", '¼': "1/4", '¾': "3/4",
	'¿': "?", '¡': "!",
	// Latin letters
	'À': "A", 'Á': "A", 'Â': "A", 'Ã': "A", 'Ä': "A", 'Å': "A", 'Ā': "A", 'Ą': "A",
	'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a", 'ā': "a", 'ą': "a",
	'Æ': "AE", 'æ': "ae", 'Ç': "C", 'Ć': "C", 'Č': "C", 'ç': "c", 'ć': "c", 'č': "c",
	'Ď': "D", 'Đ': "D", 'Ð': "D", 'ď': "d", 'đ': "d", 'ð': "d",
	'È': "E", 'É': "E", 'Ê': "E", 'Ë': "E", 'Ē': "E", 'Ę': "E", 'Ě': "E",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e", 'ē': "e", 'ę': "e", 'ě': "e",
	'Ğ': "G", 'ğ': "g",
	'Ì': "I", 'Í': "I", 'Î': "I", 'Ï': "I", 'Ī': "I", 'İ': "I",
	'ì': "i", 'í': "i", 'î': "i", 'ï': "i", 'ī': "i", 'ı': "i",
	'Ł': "L", 'ł': "l", 'Ñ': "N", 'Ń': "N", 'Ň': "N", 'ñ': "n", 'ń': "n", 'ň': "n",
	'Ò': "O", 'Ó': "O", 'Ô': "O", 'Õ': "O", 'Ö': "O", 'Ø': "O", 'Ō': "O", 'Ő': "O",
	'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o", 'ō': "o", 'ő': "o",
	'Œ': "OE", 'œ': "oe", 'Ř': "R", 'ř': "r",
	'Ś': "S", 'Š': "S", 'Ş': "S", 'ś': "s", 'š': "s", 'ş': "s", 'ß': "ss",
	'Ť': "T", 'ť': "t", 'Þ': "Th", 'þ': "th",
	'Ù': "U", 'Ú': "U", 'Û': "U", 'Ü': "U", 'Ū': "U", 'Ů': "U", 'Ű': "U",
	'ù': "u", 'ú': "u", 'û': "u", 'ü': "u", 'ū': "u", 'ů': "u", 'ű': "u",
	'Ý': "Y", 'Ÿ': "Y", 'ý': "y", 'ÿ': "y",
	'Ź': "Z", 'Ż': "Z", 'Ž': "Z", 'ź': "z", 'ż': "z", 'ž': "z",
}

// toASCII replaces known non-ASCII characters with ASCII equivalents
// and drops any other non-ASCII characters.
func toASCII(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r < 0x80:
			b.WriteRune(r)
		default:
			b.WriteString(asciiMap[r])
		}
	}
	return b.String()
}
