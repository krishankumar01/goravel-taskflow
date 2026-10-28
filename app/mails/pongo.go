package mails

import (
	"errors"

	"github.com/flosch/pongo2/v6"
)

type Pongo struct {
	views *pongo2.TemplateSet
}

func NewPongo(dir string) *Pongo {
	return &Pongo{views: pongo2.NewSet("mail", pongo2.MustNewLocalFileSystemLoader(dir))}
}

func (receiver *Pongo) Render(name string, data any) (string, error) {
	with, ok := data.(map[string]any)
	if !ok {
		return "", errors.New("pongo: the mail's With must be a map")
	}

	view, err := receiver.views.FromCache(name)
	if err != nil {
		return "", err
	}

	return view.Execute(pongo2.Context(with))
}
