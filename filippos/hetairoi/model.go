package hetairoi

type Meaning struct {
	Language   string   `json:"language"`
	Definition string   `json:"definition"`
	Notes      []string `json:"notes,omitempty"`
	Example    string   `json:"example,omitempty"`
}

type Definition struct {
	Grade    int       `json:"grade"`
	Meanings []Meaning `json:"meanings"`
}

type Noun struct {
	Declension string `json:"declension,omitempty"`
	Genitive   string `json:"genitive,omitempty"`
}

type ModernConnection struct {
	Term string `json:"term"`
	Note string `json:"note"`
}

type Verb struct {
	PrincipalParts []string `json:"principalParts"`
}

type LemmaSource struct {
	ID           string             `json:"id,omitempty"`         // if you store one
	Greek        string             `json:"greek"`                // "λόγος"
	Normalized   string             `json:"normalized,omitempty"` // "λογος"
	LinkedWord   string             `json:"linkedWord,omitempty"`
	PartOfSpeech string             `json:"partOfSpeech"`
	Article      string             `json:"article,omitempty"`
	Gender       string             `json:"gender,omitempty"`
	Noun         *Noun              `json:"noun,omitempty"`
	Verb         *Verb              `json:"verb,omitempty"`
	Definitions  []Definition       `json:"definitions,omitempty"`
	ModernConns  []ModernConnection `json:"modernConnections,omitempty"`
	// quick glosses from the top-level short fields you keep:
	English string `json:"english,omitempty"`
	Dutch   string `json:"dutch,omitempty"`
}
