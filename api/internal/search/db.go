package search

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/goutil/strutil"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/dao"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/model"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/util"
)

// DB mysql 或者sqlite3
type dBSearch struct{}

var defaultDBSearch = &dBSearch{}

// GetDBSearch 获取DB搜索实例
func GetDBSearch() Search {
	return defaultDBSearch
}
func (in *dBSearch) getHTTPClient() *resty.Client {
	panic("implement me")
}

// SearchAnswer 搜索答案
func (in *dBSearch) SearchAnswer(req model.SearchRequest) (answer [][]string, err error) {
	answer = make([][]string, 0)
	questionText := util.GetQuestionText(req.Question)
	questionHash := strutil.ShortMd5(questionText)
	tiku := dao.Tiku
	find, err := tiku.Where(tiku.QuestionHash.Eq(questionHash)).Find()
	if err != nil {
		return nil, err
	}
	if len(find) == 0 {
		find2, err := tiku.Where(tiku.QuestionText.Like("%" + questionText + "%")).Find()
		if err != nil {
			return nil, err
		}
		find = find2
	}
	for i := range find {
		var answers []string // 最后所有的答案的二维数组
		err := json.Unmarshal([]byte(find[i].Answer), &answers)
		if err != nil {
			continue
		}
		if len(answers) > 0 {
			answer = append(answer, answers)
		}
	}
	return answer, nil
}
