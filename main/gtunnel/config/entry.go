package config

//通道配置
type TunnelConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Local   string `json:"local"`

	//通道标签，这样就可以直接使用标签打开一组端口
	Tags    []string `json:"tags"`
	Active  bool `json:"active"` //是否是活动状态
}

//server ssh config
type ServerConfig struct {
	Name              string
	Description       string

	Address           string
	User              string
	Password          string
	//私钥地址
	Privatekey        string
	//链接超时时间
	DialTimeoutSecond int
	//最大重试次数
	MaxDataThroughput uint64
}

type ServerTunnelConfig struct {
	Server  *ServerConfig `json:"server"`
	Tunnels []*TunnelConfig `json:"tunnels"`
}

type Config struct {
	//客户端连接时的连接地址
	Bind   string `json:"bind"`
	//web管理台地址设置
	Web    string `json:"web"`
	//web管理台资源文件位置
	WebUI  string `json:"webui"`
	Groups []*ServerTunnelConfig `json:"server_tunnel"`
}
