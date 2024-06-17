package state

import (
	"github.com/Mikeloangel/squasher/cmd/shortener/interfaces"
	"github.com/Mikeloangel/squasher/config"
)

type State struct {
	Links interfaces.Storager
	Conf  *config.Config
}

func NewState(links interfaces.Storager, conf *config.Config) State {
	return State{
		Links: links,
		Conf:  conf,
	}
}
