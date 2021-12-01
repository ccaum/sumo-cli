package sumoapp

func (t *timeBoundary) IsEmpty() bool {
	if t.Type != "" || t.RelativeTime != "" {
		return false
	}

	return true
}

func (t *timerange) IsEmpty() bool {
	if t.Type != "" {
		return false
	}

	if !t.From.IsEmpty() {
		return false
	}

	if !t.To.IsEmpty() {
		return false
	}

	return true
}

func (t *timerange) EmptyToNil() {
	if t.To.IsEmpty() {
		t.To = nil
	}

	if t.From.IsEmpty() {
		t.From = nil
	}
}
