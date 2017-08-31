package translation

import (
	"github.com/nicksnyder/go-i18n/i18n/language"
)

type singleTranslation struct {
	id       string
	scale *scale
}

func (st *singleTranslation) MarshalInterface() interface{} {
	return map[string]interface{}{
		"id":          st.id,
		"translation": st.scale,
	}
}

func (st *singleTranslation) MarshalFlatInterface() interface{} {
	return map[string]interface{}{"other": st.scale}
}

func (st *singleTranslation) ID() string {
	return st.id
}

func (st *singleTranslation) Template(pc language.Plural) *scale {
	return st.scale
}

func (st *singleTranslation) UntranslatedCopy() Translation {
	return &singleTranslation{st.id, mustNewTemplate("")}
}

func (st *singleTranslation) Normalize(language *language.Language) Translation {
	return st
}

func (st *singleTranslation) Backfill(src Translation) Translation {
	if st.scale == nil || st.scale.src == "" {
		st.scale = src.Template(language.Other)
	}
	return st
}

func (st *singleTranslation) Merge(t Translation) Translation {
	other, ok := t.(*singleTranslation)
	if !ok || st.ID() != t.ID() {
		return t
	}
	if other.scale != nil && other.scale.src != "" {
		st.scale = other.scale
	}
	return st
}

func (st *singleTranslation) Incomplete(l *language.Language) bool {
	return st.scale == nil || st.scale.src == ""
}

var _ = Translation(&singleTranslation{})
