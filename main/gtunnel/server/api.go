package server

import (
    "fmt"
    "github.com/ihaiker/gokit/protocol/redis"
    "github.com/ihaiker/gokit/main/gtunnel/config"
    "github.com/emirpasic/gods/maps/hashmap"
    "errors"
    "strings"
    "os"
    "flag"
    "github.com/ihaiker/gokit/commons/logs"
)

type ServerApiHandler struct {
    cfg          *config.Config
    activeTunnel map[string]*GTunnel
    redisServer  *redis.Server
}

func convert(args []string) []interface{} {
    strs := make([]interface{}, len(args))
    for idx, v := range args {
        strs[idx] = v + " "
    }
    return strs
}

//列出所有服务下的通道及其现在的活动状态
func (self *ServerApiHandler) List(str... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("list ", fmt.Sprint(convert(str)...)))

    tunnels := make([][]byte, 0)
    for _, group := range self.cfg.Groups {
        tunnels = append(tunnels, []byte(fmt.Sprintf("%s(%s)", group.Server.Name, group.Server.Address)))
        tunnels = append(tunnels, []byte(fmt.Sprintf("	%s", group.Server.Description)))
        for _, s := range group.Tunnels {
            if s.Active {
                tunnels = append(tunnels, []byte(fmt.Sprintf(" [+] %-20s  %20s -> %-20s  %s", s.Name, s.Address, s.Local, strings.Join(s.Tags, ","))))
            } else {
                tunnels = append(tunnels, []byte(fmt.Sprintf("     %-20s  %20s -> %-20s  %s", s.Name, s.Address, s.Local, strings.Join(s.Tags, ","))))
            }
        }
    }
    return tunnels, nil
}

//显示所有服务
func (self *ServerApiHandler) Servers() ([][]byte, error) {
    logs.Logger("cmd").Debug("servers")

    servers := make([][]byte, 0)
    for _, group := range self.cfg.Groups {
        servers = append(servers, []byte(group.Server.Name))
    }
    return servers, nil
}

//显示所有服务的Tag
func (self *ServerApiHandler) Tags(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("tags ", fmt.Sprint(convert(args)...)))

    var cmd = flag.NewFlagSet("tags", flag.ContinueOnError)
    server := cmd.String("s", "", "限制指定服务下的所有Tag")
    if err := cmd.Parse(args); err != nil {
        return nil, err
    }

    tags := hashmap.New()
    for _, group := range self.cfg.Groups {
        if *server != "" {
            if *server != group.Server.Name {
                continue
            }
        }
        for _, s := range group.Tunnels {
            for _, tag := range s.Tags {
                if size, found := tags.Get(tag); found {
                    tags.Put(tag, (size.(int) + 1))
                } else {
                    tags.Put(tag, 1)
                }
            }
        }
    }
    out := make([][]byte, 0)
    for _, key := range tags.Keys() {
        size, _ := tags.Get(key)
        line := fmt.Sprintf("%-10s tunnel:%d", key, size)
        out = append(out, []byte(line))
    }
    return out, nil
}

//列出-s服务下的tag服务及其状态
func (self *ServerApiHandler) Tag(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("tag ", fmt.Sprint(convert(args)...)))

    var cmd = flag.NewFlagSet("tag", flag.ContinueOnError)
    server := cmd.String("s", "", "the server")
    if err := cmd.Parse(args); err != nil {
        return nil, err
    }
    if cmd.NArg() == 0 {
        return nil, errors.New("not found the tagName. Useage: tag [-s serverName] tagName")
    }
    giveName := cmd.Arg(0)
    tunnels := make([][]byte, 0)
    for _, group := range self.cfg.Groups {
        if *server != "" {
            if *server != group.Server.Name {
                continue
            }
        }
        for _, tunnel := range group.Tunnels {
            for _, tag := range tunnel.Tags {
                if tag == giveName {
                    if tunnel.Active {
                        tunnels = append(tunnels, []byte(fmt.Sprintf(" [+] %s,%s  %s->%s  %s",
                            group.Server.Name, tunnel.Name, tunnel.Address, tunnel.Local, strings.Join(tunnel.Tags, ","))))
                    } else {
                        tunnels = append(tunnels, []byte(fmt.Sprintf("     %s %s  %s->%s  %s",
                            group.Server.Name, tunnel.Name, tunnel.Address, tunnel.Local, strings.Join(tunnel.Tags, ","))))
                    }
                }
            }
        }
    }
    return tunnels, nil
}

//启动服务
//gtunnel api start -s <serverName>  [[-t <tunnelName> ] | [-g <tagName>]]
func (self *ServerApiHandler) Start(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("start ", fmt.Sprint(convert(args)...)))

    var cmd = flag.NewFlagSet("start", flag.ContinueOnError)
    serverName := cmd.String("s", "", "the server name")
    tunnelName := cmd.String("t", "", "the tunnel name")
    tagName := cmd.String("g", "", "the tag name")
    if err := cmd.Parse(args); err != nil {
        return nil, err
    } else if *serverName == "" && *tunnelName == "" && *tagName == "" {
        return nil, errors.New("Useage: start [-s <serverName> | [-t <tunnelName> ] | [-g <tagName>]] ")
    } else {
        lines := make([][]byte, 0)
        for _, group := range self.cfg.Groups {
            for _, tunnel := range group.Tunnels {
                for _, tag := range tunnel.Tags {
                    if *serverName != "" && group.Server.Name != *serverName {
                        continue
                    }
                    if *tunnelName != "" && tunnel.Name != *tunnelName {
                        continue
                    }
                    if *tagName != "" && tag != *tagName {
                        continue
                    }
                    key := group.Server.Name + ":" + tunnel.Name
                    if _, has := self.activeTunnel[key]; has {
                        lines = append(lines, []byte(key + " " + " tunnel is active."))
                    } else {
                        tunnel, err := StartTunnel(group.Server, tunnel);
                        if err == nil {
                            self.activeTunnel[key] = tunnel
                            lines = append(lines, []byte(key + "  create tunnel " + tunnel.tunnelConfig.Address + "->" + tunnel.tunnelConfig.Local))
                        } else {
                            lines = append(lines, []byte(key + " " + err.Error()))
                        }
                    }
                }
            }
        }
        return lines, nil
    }
}

//关闭服务
func (self *ServerApiHandler) Stop(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("stop ", fmt.Sprint(convert(args)...)))

    var cmd = flag.NewFlagSet("start", flag.ContinueOnError)
    serverName := cmd.String("s", "", "the server name")
    tunnelName := cmd.String("t", "", "the tunnel name")
    tagName := cmd.String("g", "", "the tag name")
    if err := cmd.Parse(args); err != nil {
        return nil, err
    } else if *serverName == "" && *tunnelName == "" && *tagName == "" {
        return nil, errors.New("Useage: start [-s <serverName> | [-t <tunnelName> ] | [-g <tagName>]] ")
    } else {
        lines := make([][]byte, 0)
        for _, group := range self.cfg.Groups {
            for _, tunnel := range group.Tunnels {
                for _, tag := range tunnel.Tags {
                    if *serverName != "" && group.Server.Name != *serverName {
                        continue
                    }
                    if *tunnelName != "" && tunnel.Name != *tunnelName {
                        continue
                    }
                    if *tagName != "" && tag != *tagName {
                        continue
                    }
                    key := group.Server.Name + ":" + tunnel.Name
                    if tunnel, has := self.activeTunnel[key]; has {
                        if err := tunnel.Stop(); err != nil {
                            lines = append(lines, []byte("stop tunnel " + key + " error. " + err.Error()))
                        } else {
                            lines = append(lines, []byte("stop tunnel " + key))
                            delete(self.activeTunnel, key)
                        }
                    }
                }
            }
        }
        return lines, nil
    }
}

//现在赢启动的服务
func (self *ServerApiHandler) Active(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug("active")

    actives := make([][]byte, 0)
    for name, lfs := range self.activeTunnel {
        actives = append(actives, []byte(fmt.Sprintf("%s %s->%s", name, lfs.tunnelConfig.Address, lfs.tunnelConfig.Local)))
    }
    return actives, nil
}

//现在赢启动的服务
func (self *ServerApiHandler) Server(args... string) ([][]byte, error) {
    logs.Logger("cmd").Debug(fmt.Sprint("server ", fmt.Sprint(convert(args)...)))
    if len(args) == 0 {
        return nil, errors.New("can't found the server name. Usage: gtunnel api server <serverName>.")
    }
    tunnels := make([][]byte, 0)
    for _, group := range self.cfg.Groups {
        for _, arg := range args {
            if arg != group.Server.Name {
                continue
            }
            tunnels = append(tunnels, []byte(fmt.Sprintf("%s(%s)", group.Server.Name, group.Server.Address)))
            tunnels = append(tunnels, []byte(fmt.Sprintf("	%s", group.Server.Description)))
            for _, s := range group.Tunnels {
                if s.Active {
                    tunnels = append(tunnels, []byte(fmt.Sprintf(" [+] %-20s  %20s -> %-20s  %s", s.Name, s.Address, s.Local, strings.Join(s.Tags, ","))))
                } else {
                    tunnels = append(tunnels, []byte(fmt.Sprintf("     %-20s  %20s -> %-20s  %s", s.Name, s.Address, s.Local, strings.Join(s.Tags, ","))))
                }
            }
        }
    }
    return tunnels, nil
}

func (self *ServerApiHandler) StopAll() ([][]byte, error) {
    logs.Logger("cmd").Debug("stopAll")

    lines := make([][]byte, 0)
    for _, group := range self.cfg.Groups {
        for _, tunnel := range group.Tunnels {
            key := group.Server.Name + ":" + tunnel.Name
            if tunnel, has := self.activeTunnel[key]; has {
                if err := tunnel.Stop(); err != nil {
                    lines = append(lines, []byte("stop tunnel " + key + " error :" + err.Error()))
                } else {
                    lines = append(lines, []byte("stop tunnel " + key))
                    delete(self.activeTunnel, key)
                }
            }
        }
    }
    return lines, nil
}

func (self *ServerApiHandler) Shutdown() ([][]byte, error) {
    logs.Logger("cmd").Debug("shutdown")

    return [][]byte{}, self.redisServer.Stop()
}

//运行客户端命令监听接口
func RunClient(cfg *config.Config, activeTunnel map[string]*GTunnel) error {
    defer func() {
        if msg := recover(); msg != nil {
            fmt.Printf("Panic: %v\n", msg)
        }
    }()
    bind := strings.SplitN(cfg.Bind, ":", 2)
    srv := redis.NewServer()
    srv.Proto = bind[0]
    srv.Address = bind[1]
    if srv.Proto == "unix" {
        os.Remove(srv.Address)
    }
    srv.RegisterHandler(&ServerApiHandler{cfg:cfg, activeTunnel:activeTunnel, redisServer:srv})
    if err := srv.ListenAndServe(); err != nil {
        return err
    }
    return nil
}