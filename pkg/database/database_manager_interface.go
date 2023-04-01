package database

import "passvault/pkg/types"

type DatabaseManagerInterface interface {
	Save(types.Entry) error
	Get(string) (*types.Entry, error)
	Contains(string) (bool, error)
	Update(entry types.Entry) error
}
