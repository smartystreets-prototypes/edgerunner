package edgerunner

type SignalingTask struct {
	inner      Task
	initialize chan<- error
	shutdown   chan<- struct{}
	closer     chan struct{}
}

func NewSignalingTask(inner Task, initialize chan<- error, shutdown chan<- struct{}) *SignalingTask {
	return &SignalingTask{
		inner:      inner,
		initialize: initialize,
		shutdown:   shutdown,
		closer:     make(chan struct{}, 2),
	}
}

func (this *SignalingTask) Init() error {
	err := this.inner.Init()
	this.initialize <- err
	return err
}

func (this *SignalingTask) Listen() {
	this.inner.Listen()
	if this.isOpen() {
		this.shutdown <- struct{}{}
	}
}
func (this *SignalingTask) isOpen() bool {
	select {
	case _, open := <-this.closer:
		return open
	default:
		return true
	}
}

func (this *SignalingTask) Close() error {
	close(this.closer)
	return this.inner.Close()
}
