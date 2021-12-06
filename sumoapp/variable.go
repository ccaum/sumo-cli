package sumoapp

import (
	"github.com/imdario/mergo"
)

func (v *variable) Merge(varObj *variable) error {
	newVar := varObj.Copy()

	if err := mergo.Merge(newVar, v, mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(v, newVar, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

func (v *variable) Copy() *variable {
	return &variable{
		Id:               v.Id,
		Name:             v.Name,
		DisplayName:      v.DisplayName,
		DefaultValue:     v.DefaultValue,
		SourceDefinition: v.SourceDefinition,
		AllowMultiSelect: v.AllowMultiSelect,
		IncludeAllOption: v.IncludeAllOption,
		HideFromUI:       v.HideFromUI,
		ValueType:        v.ValueType,
	}
}
