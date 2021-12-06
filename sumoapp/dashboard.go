package sumoapp

import (
	"fmt"

	"github.com/imdario/mergo"
)

func copyDashboardList(dList map[string]*dashboard) map[string]*dashboard {
	newList := make(map[string]*dashboard)

	for dName, d := range dList {
		newList[dName] = d.Copy()
	}

	return newList
}

func (d *dashboard) Merge(dash *dashboard) error {
	newDash := dash.Copy()

	if err := mergo.Merge(newDash, d, mergo.WithOverride); err != nil {
		return err
	}

	if err := mergo.Merge(d, newDash, mergo.WithOverride); err != nil {
		return err
	}

	return nil
}

func (d *dashboard) Copy() *dashboard {
	return &dashboard{
		Type:             d.Type,
		Name:             d.Name,
		Description:      d.Description,
		Title:            d.Title,
		Theme:            d.Theme,
		TopologyLabelMap: d.TopologyLabelMap,
		RefreshInterval:  d.RefreshInterval,
		TimeRange:        d.TimeRange, //TODO: TimeRange is a pointer and needs its own Copy() function. This just copies the pointer
		Layout:           d.Layout,
		Panels:           d.Panels,
		Variables:        d.Variables,
		RootPanel:        d.RootPanel,
		IncludeVariables: d.IncludeVariables,
	}
}

func (d *dashboard) Populate(stream *appStream) error {
	//Populate the dashboard with the panels and variables
	//It is very important the the panels and variables have been
	//loaded before calling this function

	for _, layoutPanel := range d.Layout.LayoutStructures {
		if _, ok := stream.Panels[layoutPanel.Key]; !ok {
			err := fmt.Errorf("Could not find panel '%s'. Referenced in dashboard '%s' layout", layoutPanel.Key, d.Name)
			return err
		}

		p := stream.Panels[layoutPanel.Key]

		d.Panels = append(d.Panels, p)
	}

	//If the AppendLayoutStructures is defined, the panels referenced there should be added to the dashboards
	for _, layoutPanel := range d.Layout.AppendLayoutStructures {
		//Add the layout structure the panels layout
		d.Layout.LayoutStructures = append(d.Layout.LayoutStructures, layoutPanel)

		if _, ok := stream.Panels[layoutPanel.Key]; !ok {
			err := fmt.Errorf("Could not find panel '%s'. Referenced in dashboard '%s' layout", layoutPanel.Key, d.Name)
			return err
		}

		p := stream.Panels[layoutPanel.Key]
		d.Panels = append(d.Panels, p)
	}

	for _, variableName := range d.IncludeVariables {
		v, ok := stream.Variables[variableName]
		if !ok {
			err := fmt.Errorf("Could not find variable '%s'. Referenced in dashboard '%s'", variableName, d.Name)
			return err
		}

		d.Layout.AppendLayoutStructures = make([]layoutStructure, 0)

		d.Variables = append(d.Variables, v)
	}

	return nil
}
