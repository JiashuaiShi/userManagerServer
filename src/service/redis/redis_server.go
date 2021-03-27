package redis

import (
	"github.com/go-redis/redis"
	_ "gopkg.in/yaml.v2"
	log "zego.com/userManageServer/src/logger"
	"zego.com/userManageServer/src/models"
)

var (
	client *redis.Client
	zSet   = "zSet"
)

// 初始化redis
func InitRedis(conf models.RedisConfig) (err error) {
	return initClient(conf)
}

// 根据redis配置初始化一个客户端
func initClient(conf models.RedisConfig) (err error) {
	log.Info(nil, "initClient start!")
	client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,     // redis地址
		Password: conf.Password, // redis密码，没有则留空
		DB:       conf.Db,       // 默认数据库，默认是0
	})

	// 通过 *redis.client.Ping() 来检查是否成功连接到了redis服务器 Output: PONG <nil>
	_, err = client.Ping().Result()
	return
}

// 根据userId从redis读取数据
func GetData(key string) (value string, err error) {
	var luaScript = redis.NewScript(`
		-- key 表示uid， value 表示userInfo
		local key = KEYS[1]
		local value = redis.call("GET", key)
		
		-- 查询失败返回false
		if value == false then
			return false
		end
		return value
		`)

	res, err := luaScript.Run(client, []string{key}).Result()
	logFile := log.Field{
		"key": key,
	}

	// 存在err或者lua返回false
	if err != nil || res == nil {
		log.Error(logFile, "run GetData luaScript failed, error:%+v", err)
		return
	}

	value = res.(string)
	//logFile["value"] = value
	log.Info(nil, "GetData success!")

	return value, err
}

// 根据userId向redis设置数据
func SetData(key string, value string, score int64) (err error) {
	logFile := log.Field{
		"key":   key,
		"value": value,
		"score": score,
	}

	var luaScript = redis.NewScript(`
		local zset =   KEYS[1] -- 有序集合
		local score =  ARGV[1] -- loginTime
		local member = ARGV[2] -- uid
		local key =    KEYS[2] -- uid
		local value =  ARGV[3] -- userInfo

		-- 向有序集合插入uid，uid作为member， loginTime作为score 
		local errno = redis.call("ZADD", zset, score, member)
		if errno == false then
			return errno
		end

		-- 同时记录key-value，key表示uid， value表示userInfo
		return redis.call("SET", key, value)
	`)

	member := key
	res, err := luaScript.Run(client, []string{zSet, key}, score, member, value).Result()

	// 存在err或者lua返回false
	if err != nil || res == nil {
		log.Error(logFile, "SetData luaScript return false. error: %+v", err)
		return err
	}

	log.Info(nil, "redis setData success!")
	return nil
}

// 根据userId从redis读取全部数据
func GetDataList(start int, stop int) (values []string, total int32, err error) {
	var luaScript = redis.NewScript(`
		local zset =  KEYS[1] -- 有序集合
		local start = ARGV[1] -- 查询起始值
		local stop =  ARGV[2] -- 查询结束值
		
		-- 从有序集合中获取所有记录的总数，用于返回前段计算分页总页码数
		local total = redis.call("ZCARD", zset)
		if total == false then
			return false
		end

		-- 从有序集合中获取start到stop范围的记录的uid，用于根据uid查询userInfo
		local tabKey = redis.call("ZRANGE", zset, start, stop)
		if tabKey == false then
			return false
		end

		-- 根据uid查询userInfo，用于前端展示每一页记录表格中的数据
		local tabVal = {}
		for i, v in pairs(tabKey) do
			local val = redis.call("GET", v)
			if val == false then
				return false
			end
			table.insert(tabVal, val)
		end

		-- 返回redis存储的全部记录的个数，以及userInfo的数组
		return {total, tabVal}
	`)

	result, err := luaScript.Run(client, []string{zSet}, start, stop).Result()

	// 存在err或者lua返回false
	if err != nil || result == nil {
		logField := log.Field{
			"start": start,
			"stop":  stop,
		}
		log.Error(logField, "GetData failed. error:%+v", err)
		return
	}

	// lua返回类型先转化[]interface{}
	res := result.([]interface{})

	// redis总记录数
	total = int32(res[0].(int64))

	// usersInfo数组
	vals := res[1].([]interface{})

	// 转化为[]string返回
	for _, val := range vals {
		values = append(values, val.(string))
	}

	log.Info(nil, "redis getDataList success!")
	return
}

// 删除特定用户
func DelData(key string) (err error) {
	var luaScript = redis.NewScript(`	
		local zset =   KEYS[1]  -- 有序集合
		local member = KEYS[2]  -- uid
		local key =    KEYS[2]  -- uid

		-- 先从有序集合删除uid记录，然后删除key-value中的uid记录，如果失败返回false
		if redis.call("ZREM", zset, member) ~= false then
			return redis.call("DEL", key)
		end
		return false
	`)

	res, err := luaScript.Run(client, []string{zSet, key}).Result()

	// 存在err或者lua返回false
	if err != nil || res == nil {
		logField := log.Field{
			"key": key,
			"res": res,
		}

		log.Error(logField, "DelData luaScript return false, error:%+v", err)
		return
	}

	log.Info(nil, "redis delData success!")
	return
}
