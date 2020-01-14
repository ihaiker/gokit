package remoting

type Status string

const (
	Ready   Status = "ready"
	Running Status = "running"
	Stop    Status = "stop"
)

func (status Status) IsStart() bool {
	return status == Running
}
