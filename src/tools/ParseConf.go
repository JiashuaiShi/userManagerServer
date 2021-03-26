/**
 * @Author: ryder
 * @Description:Parse ymal file to conf struct
 * @File: ParseConf
 * @Version: 1.0.0
 * @Date: 2021/3/25 6:49 下午
 */

package tools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"zego.com/userManageServer/src/models"
	_ "zego.com/userManageServer/src/models"
)

func GetConf(path string) (conf models.RedisConfig) {
	yamlBuffer, err := ioutil.ReadFile("../conf/conf.ymal")
	if err != nil {
		fmt.Println(err.Error())
	}

	if err = yaml.Unmarshal(yamlBuffer, conf); err != nil {
		fmt.Println(err.Error())
	}

	return conf
}
