package remoting

type Status string

const (
	Ready   Status = "ready"
	Running Status = "running"
	Stoped  Status = "stoped"
)

func (status Status) IsStart() bool {
	return status == Running
}
