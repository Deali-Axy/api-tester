package tester

import (
	"api-tester/pkg/utilities/goext"
	"api-tester/pkg/utilities/goext/strext"
	"encoding/json"
	"go.uber.org/zap"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

type Parser struct {
	Logger   *zap.SugaredLogger
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
					if param.Value.Schema != nil {
						queryParams[param.Value.Name] = param.Value.Schema.Value.Type
					} else {
						queryParams[param.Value.Name] = ""
					}
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
				Id: operation.OperationID,
				Name: goext.If(
					strext.IsNullOrWhiteSpace(operation.Summary),
					operation.OperationID, operation.Summary,
				),
				Description: operation.Description,
				Path:        path,
				Method:      method,
				Params:      queryParams,
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
	if err = os.WriteFile(filePath, configJSON, os.ModePerm); err != nil {
		return err
	}

	c.Logger.Infoln("Config saved successfully!")
	return nil
}
