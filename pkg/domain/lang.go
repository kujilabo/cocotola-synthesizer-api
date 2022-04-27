package domain

import (
	"golang.org/x/xerrors"

	libD "github.com/kujilabo/cocotola-synthesizer-api/pkg_lib/domain"
)

const Lang2Len = 2

// const Lang3Len = 3
const Lang5Len = 5

type Lang2 interface {
	String() string
	// ToLang3() Lang3
	ToLang5() Lang5
}

type lang2 struct {
	value string
}

func NewLang2(lang string) (Lang2, error) {
	if len(lang) != Lang2Len {
		return nil, xerrors.Errorf("invalid parameter. Lang2: %s. err: %w", lang, libD.ErrInvalidArgument)
	}

	return &lang2{
		value: lang,
	}, nil
}

func (l *lang2) String() string {
	return l.value
}

// func (l *lang2) ToLang3() Lang3 {
// 	switch l.value {
// 	case "en":
// 		return Lang3ENG
// 	case "ja":
// 		return Lang3JPN
// 	default:
// 		return Lang3Unknown
// 	}
// }

func (l *lang2) ToLang5() Lang5 {
	switch l.value {
	case "en":
		return Lang5ENUS
	case "ja":
		return Lang5JAJP
	default:
		return Lang5Unknown
	}
}

// type Lang3 interface {
// 	String() string
// }

// type lang3 struct {
// 	value string
// }

// func NewLang3(lang string) (Lang3, error) {
// 	if len(lang) != Lang3Len {
// 		return nil, fmt.Errorf("invalid parameter. Lang3: %s", lang)
// 	}

// 	return &lang3{
// 		value: lang,
// 	}, nil
// }

// func (l *lang3) String() string {
// 	return l.value
// }

type Lang5 interface {
	String() string
	ToLang2() Lang2
}

type lang5 struct {
	value string
}

func NewLang5(lang string) (Lang5, error) {
	if len(lang) != Lang5Len {
		return nil, xerrors.Errorf("invalid parameter. Lang5: %s", lang)
	}

	return &lang5{
		value: lang,
	}, nil
}

func (l *lang5) String() string {
	return l.value
}

func (l *lang5) ToLang2() Lang2 {
	switch l.value {
	case "eng":
		return Lang2EN
	case "jpn":
		return Lang2JA
	default:
		return Lang2Unknown
	}
}
