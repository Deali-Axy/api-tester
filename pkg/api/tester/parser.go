package tester

import (
	"api-tester/pkg/utilities/goext"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

type Parser struct {
	apiInfos []ApiInfo
}

// LoadFromFile 解析 OpenAPI 文档并生成测试配置
func (c *Parser) LoadFromFile(filePath string) ([]ApiInfo, error) {
	// 解析 OpenAPI 文档
	swagger, err := openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	var apiInfos []ApiInfo

	// 遍历每个路径
	for path, pathItem := range swagger.Paths.Map() {
		// 遍历每个 HTTP 方法
		for method, operation := range pathItem.Operations() {
			// 提取请求参数
			queryParams := make(map[string]string)
			for _, param := range operation.Parameters {
				if param.Value.In == "query" {
					queryParams[param.Value.Name] = goext.If(param.Value.Schema != nil, param.Value.Schema.Value.Type, "")
				}
			}

			// 提取请求体
			var bodyParams map[string]string
			if operation.RequestBody != nil {
				content := operation.RequestBody.Value.Content
				for _, mediaType := range content {
					bodyParams = make(map[string]string)
					for name, schema := range mediaType.Schema.Value.Properties {
						bodyParams[name] = schema.Value.Type
					}
				}
			}

			// 构造 ApiInfo 结构体
			apiInfo := ApiInfo{
				ApiName:     operation.OperationID,
				ApiPath:     path,
				Method:      method,
				QueryParams: queryParams,
				Body:        bodyParams,
			}

			apiInfos = append(apiInfos, apiInfo)
		}
	}

	c.apiInfos = apiInfos

	return apiInfos, nil
}

// SaveConfig 保存测试配置到 JSON 文件
func (c *Parser) SaveConfig(filePath string) error {
	// 将测试配置转换为 JSON 格式
	configJSON, err := json.MarshalIndent(c.apiInfos, "", "    ")
	if err != nil {
		return err
	}

	// 写入 JSON 文件
	err = ioutil.WriteFile(filePath, configJSON, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Config saved successfully!")
	return nil
}
