package fakturago

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type Localizer interface {
	T(key string) string
}

type I18nLocalizer struct {
	i18n *i18n.Localizer
	lang string
}

func NewLocalizer(bundle *i18n.Bundle, lang string) Localizer {
	loc := i18n.NewLocalizer(bundle, lang)
	return &I18nLocalizer{loc, lang}
}

func (l *I18nLocalizer) T(key string) string {
	value, err := l.i18n.Localize(&i18n.LocalizeConfig{MessageID: key})
	if err != nil {
		log.WithFields(log.Fields{"key": key, "lang": l.lang}).Warnf("Missing translation")
		return fmt.Sprintf("[missing translation \"%s\"]", key)
	}
	return value
}

func loadLanguageBundle(path string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	err := walkFilesWithExt(path, ".yaml", func(path string) error {
		log.WithFields(log.Fields{"file": path}).Debug("Loading language...")
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			return errors.WithMessagef(err, "language %s", path)
		}
		return nil
	})
	if err != nil {
		log.Error("Loading failed! ", err.Error())
		return nil, err
	}

	return bundle, nil
}

func walkFilesWithExt(root, ext string, walkFn func(path string) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ext {
			return nil
		}
		return walkFn(path)
	})
}
