package tcpKit

import (
    "net"
    "sync"
    "strings"
    atomicKit "github.com/ihaiker/gokit/commons/atomic"
    "errors"
    "github.com/ihaiker/gokit/commons"
    "io"
    "time"
    "github.com/ihaiker/gokit/commons/logs"
)

var (
    ErrConnClosing  = errors.New("use of closed network connection")
    ErrWriteTimeout = errors.New("write timeout")
)

type Connect struct {
    config  *Config
    connect *net.TCPConn

    once      sync.Once
    group     sync.WaitGroup
    closeFlag atomicKit.AtomicInt32

    closeChan chan struct{}
    sendChan  chan interface{}

    idleTimer *time.Timer
    idleFlag  *atomicKit.AtomicInt32

    logger *logs.LoggerEntry

    Protocol
    Handler
}

func (c *Connect) Do(done func(connect *Connect)) {
    defer func() {
        c.Close()
        done(c)
    }()
    go c.syncDo(c.writeLoop)
    go c.syncDo(c.readLoop)
    if c.config.IdleTime > 0 {
        go c.syncDo(c.heartbeatLoop)
    }
    c.waitForStop()
}

func (c *Connect) waitForStop() {
    for {
        select {
        case <-c.closeChan:
            return
        }
    }
}

func (c *Connect) Write(msg interface{}) (err error) {
    if c.IsClosed() {
        return ErrConnClosing
    }
    defer commonKit.Catch(err)
    if bytes, err := c.Protocol.Encode(msg); err != nil {
        return err
    } else if _, err := c.connect.Write(bytes); err != nil {
        if c.isCloseTCPConnect(err) {
            return ErrConnClosing
        }
        return err
    }
    return nil
}

func (c *Connect) AsyncWrite(msg interface{}, timeout time.Duration) (err error) {
    if c.IsClosed() {
        return ErrConnClosing
    }
    defer commonKit.Catch(err)

    select {
    case c.sendChan <- msg:
        return nil
    case <-time.After(timeout):
        return ErrWriteTimeout
    }
    return nil
}

func (c *Connect) PopUnSend(timeout time.Duration) (bytes interface{}) {
    select {
    case bytes = <-c.sendChan:
        return
    case <-time.After(timeout):
        return nil
    }
    return nil
}

func (c *Connect) writeLoop() {
    defer func() {
        c.logger.Debugf("close writer: %s", c.connect.RemoteAddr().String())
        c.signClose()
    }()
    for {
        select {
        case <-c.closeChan:
            return
        case msg := <-c.sendChan:
            if err := c.Write(msg); err != nil {
                if err == ErrConnClosing {
                    //TODO 已经读取出来的消息怎么转换问题
                    return
                } else {
                    c.Handler.OnEncodeError(c, msg, err)
                }
            }
        }
    }
}

func (c *Connect) readLoop() {
    defer func() {
        c.logger.Debugf("close reader: %s", c.connect.RemoteAddr().String())
        c.signClose()
    }()

    for {
        select {
        case <-c.closeChan:
            return
        default:
            if msg, err := c.Decode(c.connect); err != nil {
                if c.isCloseTCPConnect(err) {
                    return
                }
                c.Handler.OnDecodeError(c, err)
            } else {
                c.resetReadIdle()
                if c.config.AsynHandler {
                    go commonKit.Try(func() {
                        c.Handler.OnMessage(c, msg)
                    }, func(err interface{}) {
                        c.Handler.OnError(c, err.(error), msg)
                    })
                } else {
                    commonKit.Try(func() {
                        c.Handler.OnMessage(c, msg)
                    }, func(err interface{}) {
                        c.Handler.OnError(c, err.(error), msg)
                    })
                }
            }
        }
    }
}

func (c *Connect) heartbeatLoop() {
    defer func() {
        c.idleTimer.Stop()
        c.logger.Debugf("close heartbeat: %s", c.connect.RemoteAddr().String())
        c.signClose()
    }()
    c.idleTimer = time.NewTimer(time.Millisecond * time.Duration(c.config.IdleTime) )
    c.idleFlag = atomicKit.NewInt32()

    for {
        select {
        case <-c.closeChan:
            return
        case <-c.idleTimer.C:
            if c.idleFlag.GetAndIncrement() < int32(c.config.IdleTimeout) {
                c.idleTimer.Reset(time.Millisecond * time.Duration(c.config.IdleTime))
                c.Handler.OnIdle(c)
            } else {
                c.logger.Debugf("ReadTimerOut : %s", c.connect.RemoteAddr().String())
                c.signClose()
                return
            }
        }
    }
}

func (c *Connect) resetReadIdle() {
    if c.idleTimer != nil {
        c.idleTimer.Reset(time.Millisecond * time.Duration(c.config.IdleTime))
        c.idleFlag.Set(0)
    }
}

func (c *Connect) isCloseTCPConnect(err error) bool {
    if strings.Contains(err.Error(), ErrConnClosing.Error()) || strings.Contains(err.Error(), "connection reset by peer") {
        return true
    } else if err == io.EOF {
        return true
    }
    return false
}

func (c *Connect) IsClosed() bool {
    return c.closeFlag.Get() == 1
}

func (c *Connect) signClose() {
    if !c.IsClosed() {
        c.logger.Debugf("发送TCP关闭信息：%s", c.connect.RemoteAddr().String())
        c.closeFlag.Set(1)
        close(c.closeChan)
    }
}

func (c *Connect) Close() {
    c.once.Do(func() {
        c.logger.Debugf("关闭TCP连接[Start]：%s", c.connect.RemoteAddr().String())
        c.group.Add(1)
        c.signClose()
        c.Handler.OnClose(c)
        close(c.sendChan)
        c.connect.Close()
        c.group.Done()
        c.logger.Debugf("关闭TCP连接[End]：%s", c.connect.RemoteAddr().String())
    })
    c.group.Wait()
}

func (c *Connect) syncDo(fn func()) {
    c.group.Add(1)
    defer c.group.Done()
    fn()
}

func newConnect(s *Server, c *net.TCPConn) *Connect {
    conn := &Connect{
        config:    s.config,
        connect:   c,
        logger:    s.logger,
        closeChan: make(chan struct{}),
        sendChan:  make(chan interface{}, s.config.PacketSendChanLimit),
    }
    conn.Handler = s.maker.handler(c)
    conn.Protocol = s.maker.protocol(c)
    return conn
}
