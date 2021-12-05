package sumoapp

import (
	"github.com/imdario/mergo"
)

func copyPanelList(pList map[string]panel) map[string]*panel {
	newList := make(map[string]*panel)

	for pName, p := range pList {
		newList[pName] = p.Copy()
	}

	return newList
}

func (p *panel) Merge(pan *panel) error {
	newPanel := pan.Copy()

	if err := mergo.Merge(newPanel, p, mergo.WithOverride); err != nil {
		return err
	}

	//Easiest way to update self with the newly merged
	//panel value. There might be a more efficient way to do this
	if err := mergo.Merge(p, newPanel, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

func (p *panel) Copy() *panel {
	pan := &panel{}
	mergo.Merge(pan, p)
	return pan
}
