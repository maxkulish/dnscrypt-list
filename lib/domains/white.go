package domains

// WhiteList keeps the map of domains
type WhiteList struct {
	Store map[string]string
}

// Add adds domain to the map
func (w *WhiteList) Add(key, value string) {
	if _, ok := w.Store[key]; !ok {
		w.Store[key] = value
	}
}
