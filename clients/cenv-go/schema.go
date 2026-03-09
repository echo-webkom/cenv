package cenv

type Schema struct {
	Entries []Entry `toml:"entries"`
}

type Entry struct {
	Key            string     `toml:"key"`
	Hint           *string    `toml:"hint,omitempty"`
	Required       bool       `toml:"required"`
	Default        *string    `toml:"default,omitempty"`
	LegalValues    []string   `toml:"legal_values,omitempty"`
	RequiredLength *int       `toml:"required_length,omitempty"`
	RegexMatch     *string    `toml:"regex_match,omitempty"`
	Kind           *EntryKind `toml:"kind,omitempty"`
}

type EntryKind struct {
	Type string `toml:"type"`

	MinInt   *int64   `toml:"min_int,omitempty"`
	MaxInt   *int64   `toml:"max_int,omitempty"`
	MinFloat *float64 `toml:"min_float,omitempty"`
	MaxFloat *float64 `toml:"max_float,omitempty"`
}
