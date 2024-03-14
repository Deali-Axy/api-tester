package tester

import (
	"api-tester/pkg/utilities/goext"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

// Tester 定义工具结构
type Tester struct {
	RestyClient *resty.Client
	Logger      *zap.SugaredLogger
	BaseURL     string
	AuthToken   string
}

// Request 发送请求
func (c *Tester) Request(apiName string, apiPath string, method string,
	queryParams map[string]string, body interface{},
) (*Report, error) {
	report := &Report{
		ApiName:  apiName,
		ApiPath:  apiPath,
		IsPassed: false,
		Elapsed:  time.Duration(-1),
	}

	// 开始计时
	start := time.Now()

	// 构建请求
	req := c.RestyClient.R()

	if len(c.AuthToken) > 0 {
		req.SetHeader("Authorization", "token "+c.AuthToken)
	}

	// 设置查询参数
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 发送请求
	var resp *resty.Response
	var err error
	switch method {
	case http.MethodGet:
		resp, err = req.Get(c.BaseURL + apiPath)
	case http.MethodPost:
		resp, err = req.Post(c.BaseURL + apiPath)
	case http.MethodPut:
		resp, err = req.Put(c.BaseURL + apiPath)
	case http.MethodDelete:
		resp, err = req.Delete(c.BaseURL + apiPath)
	default:
		return report, fmt.Errorf("不支持的 HTTP 方法: %s", method)
	}
	if err != nil {
		return report, fmt.Errorf("请求失败: %v", err)
	}

	// 停止计时并计算响应时间
	elapsed := time.Since(start)

	report.Elapsed = elapsed
	report.Response = resp

	// 打印响应时间
	c.Logger.Infof("接口 %s 的响应时间: %s", apiPath, elapsed)

	// 解析响应
	var data ApiResponse
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return report, fmt.Errorf("解析响应失败: %v\n", err)
	}

	// 打印响应信息
	c.Logger.Infof("响应code: %d, 响应信息: %s", data.Code, data.Message)

	report.IsPassed = goext.If(resp.StatusCode() >= 400, false, true)
	report.Data = &data

	return report, nil
}

// TestApis 测试接口
func (c *Tester) TestApis(apiInfos []ApiInfo, concurrencyLimit int) ([]*Report, error) {
	if concurrencyLimit < 0 {
		return nil, fmt.Errorf("concurrencyLimit can not be negative")
	}
	if concurrencyLimit == 1 {
		return c.testApisSync(&apiInfos)
	}

	return c.testApisParallel(&apiInfos, concurrencyLimit)
}

// testApisSync 测试接口 - 同步
func (c *Tester) testApisSync(apiInfos *[]ApiInfo) ([]*Report, error) {
	// 通过切片的方式提高数组操作性能，cap 为 []ApiInfo 长度
	reports := make([]*Report, 0, len(*apiInfos))

	// 遍历接口信息
	for _, info := range *apiInfos {
		// 调用 Tester 的 Request 方法测试接口
		report, err := c.Request(info.Name, info.Path, info.Method, info.Params, info.Body)
		if err != nil {
			c.Logger.Errorf("接口 %s 测试失败：%v", info.Name, err)
		}

		// 将测试报告保存到数组中
		reports = append(reports, report)
	}

	// 按照 Elapsed 属性排序，从大到小
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Elapsed > reports[j].Elapsed
	})

	return reports, nil
}

// testApisParallel 测试接口 - 并行
func (c *Tester) testApisParallel(apiInfos *[]ApiInfo, concurrencyLimit int) ([]*Report, error) {
	// 创建等待组，以等待所有请求完成
	var wg sync.WaitGroup

	// 创建管道用于接收测试结果
	results := make(chan *Report, len(*apiInfos))

	// 创建一个带缓冲的通道，用于限制并发数量
	semaphore := make(chan struct{}, concurrencyLimit)

	// 遍历接口信息并启动并发请求
	for _, info := range *apiInfos {
		// 控制并发数量，当通道已满时会阻塞
		// 尝试向信号量通道中写入信号量，如果通道已满则会阻塞直到有空位
		semaphore <- struct{}{}

		wg.Add(1)
		go func(apiInfo ApiInfo) {
			// 在函数退出时通知等待组完成，并释放信号量
			defer func() {
				wg.Done()
				<-semaphore
			}()

			// 调用 Tester 的 Request 方法测试接口
			report, err := c.Request(apiInfo.Name, apiInfo.Path, apiInfo.Method, apiInfo.Params, apiInfo.Body)
			if err != nil {
				c.Logger.Errorf("接口 %s 测试失败：%v", apiInfo.Name, err)
			}

			// 将测试报告发送到管道
			results <- report
		}(info)
	}

	// 等待所有请求完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集测试结果
	reports := make([]*Report, 0, len(*apiInfos))
	for result := range results {
		reports = append(reports, result)
	}

	// 按照 Elapsed 属性排序，从大到小
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Elapsed > reports[j].Elapsed
	})

	return reports, nil
}
