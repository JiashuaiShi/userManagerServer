/**
 * @Author: ryder
 * @Description:
 * @File: redisConf
 * @Version: 1.0.0
 * @Date: 2021/3/25 6:40 下午
 */
package models

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
