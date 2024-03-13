package tester

import (
	"encoding/json"
	"os"
)

// ReadConfig 读取 JSON 配置文件并解析
func ReadConfig(filename string) ([]ApiInfo, error) {
	var apiInfos []ApiInfo

	// 读取 JSON 文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 数据
	err = json.Unmarshal(data, &apiInfos)
	if err != nil {
		return nil, err
	}

	return apiInfos, nil
}
