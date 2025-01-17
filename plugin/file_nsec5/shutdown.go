package file_nsec5

// OnShutdown shuts down any running go-routines for this zone.
func (z *Zone) OnShutdown() error {
	if 0 < z.ReloadInterval {
		z.reloadShutdown <- true
	}
	return nil
}
