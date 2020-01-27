package support

type Supporter struct {
	progress int  // The progress of the analysis
	stop     bool // If the analysis is stoped
}

// Returns the progress of the analysis
func NewSupporter() *Supporter {
	return &Supporter{
		progress: 0,
		stop:     false,
	}
}

// Returns the progress of the analysis
func (sup *Supporter) Progress() int {
	return sup.progress
}

// Increments the progress of the analysis
func (sup *Supporter) IncrementProgress() {
	sup.progress++
}

// Tells the supported to stop the analysis
// It will just finish the current computations
func (sup *Supporter) Cancel() {
	sup.stop = true
}

// Tells if hasbeen canceled or not
func (sup *Supporter) Canceled() bool {
	return sup.stop
}
