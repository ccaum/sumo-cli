package sumoapp

import "github.com/imdario/mergo"

func (s *savedSearch) Merge(search *savedSearch) error {
	if err := mergo.Merge(s, search, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}
