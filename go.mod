module github.com/pastDexter/fakturago

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/imdario/mergo v0.3.11
	github.com/johnfercher/maroto v0.28.0
	github.com/jung-kurt/gofpdf v1.4.2
	github.com/nicksnyder/go-i18n/v2 v2.1.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
	golang.org/x/text v0.3.4
)

replace github.com/johnfercher/maroto => github.com/pastDexter/maroto v0.28.1-0.20201120091517-b440edc06f9a
