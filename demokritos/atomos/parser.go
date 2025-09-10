package atomos

import (
	"strings"
	"unicode"
)

type Parsed struct {
	Lemma        string
	Article      string
	Gender       string
	PartOfSpeech string
	Genitive     string
	Declension   string
}

var articleGender = map[string]string{
	"ὁ":  "masc",
	"ἡ":  "fem",
	"τό": "neut",
	"το": "neut",
}

// replace en/em-dashes with hyphen, collapse spaces, trim
func normalize(s string) string {
	r := strings.NewReplacer("—", "-", "–", "-", "  ", " ")
	return strings.TrimSpace(r.Replace(s))
}

// SimpleParse implements the minimal rules:
//   - "lemma -GEN, ARTICLE"  -> lemma/gen/article
//   - "ARTICLE lemma"        -> article/lemma
//   - otherwise lemma only
func SimpleParse(s string) Parsed {
	s = normalize(s)
	p := Parsed{}

	// Case 1: "ARTICLE lemma" (e.g., "ὁ λόγος")
	if hasLeadingArticle(s) {
		parts := strings.Fields(s)
		p.Article = parts[0]
		if g, ok := articleGender[p.Article]; ok {
			p.Gender = g
			p.PartOfSpeech = "noun"
		}
		if len(parts) > 1 {
			p.Lemma = strings.Join(parts[1:], " ")
		}
		// best-effort declension if looks like second (-ος)
		if strings.HasSuffix(p.Lemma, "ος") {
			p.Declension = "second"
			p.Genitive = "-ου"
		}
		return p
	}
	// Case 2: "lemma -GEN, ARTICLE"
	if strings.Contains(s, ",") && strings.Contains(s, "-") {
		leftRight := strings.SplitN(s, ",", 2)
		left := strings.TrimSpace(leftRight[0])
		right := strings.TrimSpace(leftRight[1])

		// right is the article, usually exactly "ὁ/ἡ/τό"
		p.Article = right
		if g, ok := articleGender[p.Article]; ok {
			p.Gender = g
			p.PartOfSpeech = "noun"
		}

		// left is "lemma -GEN"
		lparts := strings.SplitN(left, "-", 2)
		p.Lemma = strings.TrimSpace(lparts[0])
		if len(lparts) > 1 {
			p.Genitive = strings.TrimSpace("-" + strings.TrimSpace(lparts[1]))
		}

		// Declension heuristics (tiny & safe) using accent-stripped checks
		baseLemma := stripAccents(p.Lemma)
		baseGen := stripAccents(p.Genitive) // e.g., "-ης", "-ας", "-ου"

		// First declension (very common fem in -η/-α with gen -ης/-ας)
		if endsWithBase(baseLemma, "η") || endsWithBase(baseLemma, "α") || baseGen == "-ης" || baseGen == "-ας" {
			p.Declension = "first"
			// If genitive was missing, set a sensible default based on lemma ending
			if p.Genitive == "" {
				if endsWithBase(baseLemma, "η") {
					p.Genitive = "-ης"
				} else if endsWithBase(baseLemma, "α") {
					p.Genitive = "-ας"
				}
			}
			return p
		}

		// Second declension (very common masc/neut in -ος with gen -ου)
		if endsWithBase(baseLemma, "ος") || baseGen == "-ου" {
			p.Declension = "second"
			if p.Genitive == "" {
				p.Genitive = "-ου"
			}
			return p
		}

		// Unknown declension: leave empty but keep noun fields we have
		return p
	}

	if isLikelyVerb(s) {
		p.PartOfSpeech = "verb"
	}

	// Fallback: lemma only
	p.Lemma = s
	return p
}

func hasLeadingArticle(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "ὁ ") || strings.HasPrefix(s, "ἡ ") || strings.HasPrefix(s, "τό ") || strings.HasPrefix(s, "το ")
}

// Optional: quick proper-noun helper (unused in simple parser but handy)
func looksProperNoun(lemma string) bool {
	rs := []rune(lemma)
	return len(rs) > 0 && unicode.IsUpper(rs[0])
}

func stripAccents(s string) string {
	r := strings.NewReplacer(
		"ά", "α", "ὰ", "α", "ᾶ", "α",
		"έ", "ε", "ὲ", "ε",
		"ή", "η", "ὴ", "η", "ῆ", "η",
		"ί", "ι", "ὶ", "ι", "ῖ", "ι",
		"ό", "ο", "ὸ", "ο",
		"ύ", "υ", "ὺ", "υ", "ῦ", "υ",
		"ώ", "ω", "ὼ", "ω", "ῶ", "ω",
		"ῆς", "ης", "ᾶς", "ας", // common genitive singular endings
	)
	return r.Replace(s)
}

func endsWithBase(s, suffix string) bool {
	return strings.HasSuffix(stripAccents(s), suffix)
}

func isLikelyVerb(lemma string) bool {
	n := len([]rune(lemma))
	if n == 0 {
		return false
	}
	if strings.HasSuffix(lemma, "ω") || strings.HasSuffix(lemma, "ομαι") || strings.HasSuffix(lemma, "μι") {
		return true
	}
	// Contracted verbs like τιμάω/ποιέω/τηρέω:
	if strings.HasSuffix(lemma, "άω") || strings.HasSuffix(lemma, "έω") || strings.HasSuffix(lemma, "όω") {
		return true
	}
	return false
}
