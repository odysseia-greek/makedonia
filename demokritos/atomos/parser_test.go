package atomos

import (
	"encoding/json"
	"testing"
)

type Case struct {
	Input string `json:"input"`
	Want  struct {
		Lemma        string `json:"lemma"`
		Article      string `json:"article"`
		Gender       string `json:"gender"`
		PartOfSpeech string `json:"partOfSpeech"`
		Genitive     string `json:"genitive"`
		Declension   string `json:"declension"`
	} `json:"want"`
}

const casesJSON = `
[
  {
    "input": "τέχνη –ης, ἡ",
    "want": {
      "lemma": "τέχνη",
      "article": "ἡ",
      "gender": "fem",
      "partOfSpeech": "noun",
      "genitive": "-ης",
      "declension": "first"
    }
  },
  {
    "input": "φυλακή -ῆς, ἡ",
    "want": {
      "lemma": "φυλακή",
      "article": "ἡ",
      "gender": "fem",
      "partOfSpeech": "noun",
      "genitive": "-ῆς",
      "declension": "first"
    }
  },
  {
    "input": "λόγος –ου, ὁ",
    "want": {
      "lemma": "λόγος",
      "article": "ὁ",
      "gender": "masc",
      "partOfSpeech": "noun",
      "genitive": "-ου",
      "declension": "second"
    }
  },
  {
    "input": "ὁ λόγος",
    "want": {
      "lemma": "λόγος",
      "article": "ὁ",
      "gender": "masc",
      "partOfSpeech": "noun",
      "genitive": "-ου",
      "declension": "second"
    }
  },
  {
    "input": "τηρέω",
    "want": {
      "lemma": "τηρέω",
      "article": "",
      "gender": "",
      "partOfSpeech": "verb",
      "genitive": "",
      "declension": ""
    }
  }
]
`

func TestSimpleParse_JSON(t *testing.T) {
	var cases []Case
	if err := json.Unmarshal([]byte(casesJSON), &cases); err != nil {
		t.Fatalf("bad test json: %v", err)
	}
	for i, tc := range cases {
		got := SimpleParse(tc.Input)
		if got.Lemma != tc.Want.Lemma {
			t.Errorf("[%d] lemma: got=%q want=%q (input=%q)", i, got.Lemma, tc.Want.Lemma, tc.Input)
		}
		if got.Article != tc.Want.Article {
			t.Errorf("[%d] article: got=%q want=%q (input=%q)", i, got.Article, tc.Want.Article, tc.Input)
		}
		if got.Gender != tc.Want.Gender {
			t.Errorf("[%d] gender: got=%q want=%q (input=%q)", i, got.Gender, tc.Want.Gender, tc.Input)
		}
		if got.PartOfSpeech != tc.Want.PartOfSpeech {
			t.Errorf("[%d] pos: got=%q want=%q (input=%q)", i, got.PartOfSpeech, tc.Want.PartOfSpeech, tc.Input)
		}
		if got.Genitive != tc.Want.Genitive {
			t.Errorf("[%d] gen: got=%q want=%q (input=%q)", i, got.Genitive, tc.Want.Genitive, tc.Input)
		}
		if got.Declension != tc.Want.Declension {
			t.Errorf("[%d] decl: got=%q want=%q (input=%q)", i, got.Declension, tc.Want.Declension, tc.Input)
		}
	}
}
