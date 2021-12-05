package sumoapp

import (
	"github.com/imdario/mergo"
)

func (v *variable) Merge(varObj *variable) error {
	if err := mergo.Merge(v, varObj, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}
