package suture

type generatedService struct {
	name          string
	serveFunc     func()
	stopFunc      func()
	isCompletable bool
}

// GenerateService converts a name, a serve function, a stop function, and a
// bool indicating whether the service is completable into a full Service. This
// is useful when you would otherwise be overloading the Service interface on a
// type.
func GenerateService(name string, serveFunc func(), stopFunc func(), isCompletable bool) Service {
	return &generatedService{
		name:          name,
		serveFunc:     serveFunc,
		stopFunc:      stopFunc,
		isCompletable: isCompletable,
	}
}

// Serve implements Service
func (gs *generatedService) Serve() {
	gs.serveFunc()
}

// Stop implements Service
func (gs *generatedService) Stop() {
	gs.stopFunc()
}

// String implements fmt.Stringer
func (gs *generatedService) String() string {
	return gs.name
}

// IsCompletable implements the optional interface described in Service
func (gs *generatedService) IsCompletable() bool {
	return gs.isCompletable
}
