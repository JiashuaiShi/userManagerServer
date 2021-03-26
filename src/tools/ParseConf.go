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
)

func ParseConf(path string) (conf models.Config, err error) {
	//str, _ := os.Getwd()
	//fmt.Println(str)
	yamlBuffer, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("config path error")
		return conf, err
	}
	err = yaml.Unmarshal(yamlBuffer, &conf)
	return conf, err
}