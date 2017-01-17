package redis

//所有redis的协议内容
type RedisHandler interface {
	//key
	Set(key string,value []byte) error
	Del(keys []string) (int,error)
	//序列化给定 key ，并返回被序列化的值，
	Dump(key string) ([]byte,error)
	Exists(key string) (int,error)
	
	//设置过期时间
	Expire(key string,seconds int) (int,error)
	//设置过期时间，时间点，unix timestamp
	ExpireAt(key string,unixTimestamp string) (int,error)
	
	//查找匹配的键
	Keys(pattern string) ([]string,error)
	
	//将key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
	Migrate(host string,port int,key string, destination_db int, timeout int, method string/*[COPY] [REPLACE]*/) (int,error)
	
	//移动key到db库
	Move(key string, db int) (int,error)
	
	
	//OBJECT REFCOUNT <key> 返回给定 key 引用所储存的值的次数。此命令主要用于除错
	//OBJECT ENCODING <key> 返回给定 key 锁储存的值所使用的内部表示(representation)。
	//OBJECT IDLETIME <key> 返回给定 key 自储存以来的空闲时间(idle， 没有被读取也没有被写入)，以秒为单位。
	Object(key string,subCommand string) (int,error)
	
	//移除key的生存时间
	Persist(key string) (int,error)
}
