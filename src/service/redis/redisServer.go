package redis

import (
	"github.com/go-redis/redis"
	_ "gopkg.in/yaml.v2"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/models"
)

var (
	// 声明一个全局的Client变量
	client *redis.Client
	uidSet = "uidSet"
)

func InitRedis(conf models.RedisConfig) (err error) {
	return initClient(conf)
}

// 根据redis配置初始化一个客户端
func initClient(conf models.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,     // redis地址
		Password: conf.Password, // redis密码，没有则留空
		DB:       conf.Db,       // 默认数据库，默认是0
	})

	// 通过 *redis.client.Ping() 来检查是否成功连接到了redis服务器
	// Output: PONG <nil>
	_, err = client.Ping().Result()
	return
}

// 根据userId从redis读取数据
func GetData(key string) (values []string, err error) {
	var luaScript = redis.NewScript(`
		local value = redis.call("GET", KEYS[1])
		if value == false then
			return false
		end
		return value
		`)

	res, err := luaScript.Run(client, []string{key}).Result()
	if err != nil || res == nil {
		log.Error(
			log.Field{
				"key": key,
				"res": res,
			},
			"run GetData luaScript failed, error:%+v", err)
		return
	}
	values = append(values, res.(string))
	return values, err
}

// 根据userId向redis设置数据
func SetData(key string, value string, score int64) (err error) {
	var luaScript = redis.NewScript(`
		local zset = KEYS[1]
		local score = ARGV[1]
		local member = ARGV[2]
		local value = ARGV[3]
		local key = KEYS[2]

		local errno = redis.call("ZADD", zset, score, member)
		if errno == false then
			return errno
		end

		return redis.call("SET", key, value)
	`)
	member := key
	res, err := luaScript.Run(client, []string{uidSet, key}, score, member, value).Result()

	// 先打印错误日志，然后返回error？还是错误只在一处打印就够了呢？
	if err != nil || res == nil {
		log.Error(log.Field{
			"key":    key,
			"score":  score,
			"member": member,
			"value":  value,
			"res":    res,
		}, "SetData luaScript return false. error: %+v", err)
		return err
	}
	return err
}

// 根据userId从redis读取全部数据
func GetAllData(minv int, maxv int) (values []string, total int32, err error) {
	var luaScript = redis.NewScript(`
		local zset = KEYS[1]
		local minv = ARGV[1]
		local maxv = ARGV[2]

		local total = redis.call("ZCARD", zset)
		if total == false then
			return false
		end

		local tabKey = redis.call("ZRANGE", zset, minv, maxv)
		if tabKey == false then
			return false
		end

		local tabVal = {}
		for i, v in pairs(tabKey) do
			local val = redis.call("GET", v)
			if val == false then
				return false
			end
			table.insert(tabVal, val)
		end

		return {total, tabVal}
	`)

	result, err := luaScript.Run(client, []string{uidSet}, minv, maxv).Result()
	if err != nil || result == nil {
		log.Error(log.Field{
			"minv": minv,
			"maxv": maxv,
		}, "GetData failed. error:%+v", err)
		return
	}

	res := result.([]interface{})
	vals := res[1].([]interface{})
	total = int32(res[0].(int64))

	for _, val := range vals {
		values = append(values, val.(string))
	}

	return
}

// 删除特定用户
func DelData(key string) (err error) {
	//删除set中member, 再删除key
	var luaScript = redis.NewScript(`	
		local zset = KEYS[1]
		local member = KEYS[2]
		local key = KEYS[2]

		if redis.call("ZREM", zset, member) ~= false then
			return redis.call("DEL", key)
		end
		return false
	`)
	res, err := luaScript.Run(client, []string{uidSet, key}).Result()
	if err != nil || res == nil {
		log.Error(log.Field{
			"key": key,
			"res": res,
		}, "DelData luaScript return false, error:%+v", err)
		return
	}

	return
}
