package storage

import (
	"github.com/victorfernandesraton/bushido/bushido"
	"github.com/victorfernandesraton/bushido/storage/sqlite"
)

func DatabseFactory() (bushido.LocalStorage, error) {
	st, err := sqlite.New()
	return st, err
}
