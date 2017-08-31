package translation

import (
	"github.com/nicksnyder/go-i18n/i18n/language"
)

type pluralTranslation struct {
	id        string
	scales map[language.Plural]*scale
}

func (pt *pluralTranslation) MarshalInterface() interface{} {
	return map[string]interface{}{
		"id":          pt.id,
		"translation": pt.scales,
	}
}

func (pt *pluralTranslation) MarshalFlatInterface() interface{} {
	return pt.scales
}

func (pt *pluralTranslation) ID() string {
	return pt.id
}

func (pt *pluralTranslation) Template(pc language.Plural) *scale {
	return pt.scales[pc]
}

func (pt *pluralTranslation) UntranslatedCopy() Translation {
	return &pluralTranslation{pt.id, make(map[language.Plural]*scale)}
}

func (pt *pluralTranslation) Normalize(l *language.Language) Translation {
	// Delete plural categories that don't belong to this language.
	for pc := range pt.scales {
		if _, ok := l.Plurals[pc]; !ok {
			delete(pt.scales, pc)
		}
	}
	// Create map entries for missing valid categories.
	for pc := range l.Plurals {
		if _, ok := pt.scales[pc]; !ok {
			pt.scales[pc] = mustNewTemplate("")
		}
	}
	return pt
}

func (pt *pluralTranslation) Backfill(src Translation) Translation {
	for pc, t := range pt.scales {
		if t == nil || t.src == "" {
			pt.scales[pc] = src.Template(language.Other)
		}
	}
	return pt
}

func (pt *pluralTranslation) Merge(t Translation) Translation {
	other, ok := t.(*pluralTranslation)
	if !ok || pt.ID() != t.ID() {
		return t
	}
	for pluralCategory, scale := range other.scales {
		if scale != nil && scale.src != "" {
			pt.scales[pluralCategory] = scale
		}
	}
	return pt
}

func (pt *pluralTranslation) Incomplete(l *language.Language) bool {
	for pc := range l.Plurals {
		if t := pt.scales[pc]; t == nil || t.src == "" {
			return true
		}
	}
	return false
}

var _ = Translation(&pluralTranslation{})
