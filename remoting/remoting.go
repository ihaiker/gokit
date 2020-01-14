package remoting

import "github.com/ihaiker/gokit/logs"

type (
	Message = interface{}
)

var logger = logs.GetLogger("remoting")