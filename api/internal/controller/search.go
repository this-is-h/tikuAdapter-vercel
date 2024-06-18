package controller

import (
	"encoding/json" // 用于 JSON 编码和解码
	// "fmt"
	"net/http" // 标准 HTTP 包
	"strings"
	"sync"

	"github.com/this-is-h/tikuAdapter-vercel/api/internal/middleware"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/registry/manager"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/search"
	// "github.com/this-is-h/tikuAdapter-vercel/api/pkg/global"
	// "github.com/this-is-h/tikuAdapter-vercel/api/pkg/logger"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/model"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/util"
)

// 搜题接口处理函数
func Search(w http.ResponseWriter, r *http.Request) {
	var req model.SearchRequest
	decoder := json.NewDecoder(r.Body) // 创建 JSON 解码器
	err := decoder.Decode(&req) // 将请求体中的 JSON 解码到 req 结构体中
	if err != nil {
		http.Error(w, http.StatusText(400), 400) // 参数错误时返回 400 状态码和错误信息
		return
	}

	var result [][]string // 存储所有答案的二维数组
	// var localAnswer [][]string // 存储本地答案的二维数组
	query := r.URL.Query() // 获取 URL 查询参数
	use := query.Get("use") // 获取 use 参数的值


	// // 如果使用本地题库
	// if strings.Contains(use, "local") {
	// 	localAnswer, err = search.GetDBSearch().SearchAnswer(req) // 查询本地答案
	// 	if err != nil {
	// 		logger.SysError(fmt.Sprintf("查询本地答案出错：%s", err.Error())) // 查询本地答案出错时记录日志
	// 	}
	// 	result = append(result, localAnswer...) // 将本地答案追加到结果集中
	// }

	// 如果本地没有结果，再查询第三方
	// if len(result) == 0 {
		var clients = []search.Search{
			&search.BuguakeClient{
				Enable: strings.Contains(use, "buguake") || use == "", // 判断是否启用 BuguakeClient
			},
			&search.IcodefClient{
				Token:  query.Get("icodefToken"),
				Enable: strings.Contains(use, "icodef") || use == "", // 判断是否启用 IcodefClient
			},
			&search.WannengClient{
				Token:  query.Get("wannengToken"),
				Enable: strings.Contains(use, "wanneng") || use == "", // 判断是否启用 WannengClient
			},
			&search.EnncyClient{
				Token:  query.Get("enncyToken"),
				Enable: strings.Contains(use, "enncy"), // 判断是否启用 EnncyClient
			},
			&search.AidianClient{
				Enable: strings.Contains(use, "aidian"), // 判断是否启用 AidianClient
				YToken: query.Get("aidianYToken"),
			},
			&search.LemonClient{
				Enable: strings.Contains(use, "lemon"), // 判断是否启用 LemonClient
				Token:  query.Get("lemonToken"),
			},
		}

		cfg := manager.GetManager().GetConfig() // 获取配置管理器的配置
		for _, api := range cfg.API {
			if strings.Contains(use, api.Name) {
				clients = append(clients, api) // 动态添加第三方 API 客户端
			}
		}

		var wg sync.WaitGroup // 创建 WaitGroup 同步多个 goroutine
		var mu sync.Mutex // 创建互斥锁保护共享数据

		for i := range clients {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				res, err := clients[idx].SearchAnswer(req) // 查询第三方答案
				if err == nil && len(res) > 0 {
					mu.Lock()
					defer mu.Unlock()
					result = append(result, res...) // 将第三方答案追加到结果集中
				}
			}(i)
		}
		wg.Wait() // 等待所有 goroutine 完成
	// }

	resp := util.FillAnswerResponse(result, &req) // 填充答案响应

	if query.Get("collect") != "" {
		middleware.CollectAnswer(resp) // 收集答案
	}

	w.Header().Set("Content-Type", "application/json") // 设置响应头
	json.NewEncoder(w).Encode(resp) // 返回 JSON 响应
}
