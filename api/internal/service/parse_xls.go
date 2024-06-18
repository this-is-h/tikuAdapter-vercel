package service

import (
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/model"
	"github.com/this-is-h/tikuAdapter-vercel/api/pkg/util"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

// XLSXOptions xlsx选项
type XLSXOptions struct {
	SheetName string   `json:"sheet"`
	Question  string   `json:"q"`
	Answer    string   `json:"a"`
	Option    []string `json:"o"`
}

// ParseXls 解析xlsx文件
func ParseXls(file *excelize.File, opt XLSXOptions) []model.Question {
	rows, err := file.GetRows(opt.SheetName)
	if err != nil {
		return nil
	}
	var tikus []model.Question
	for i := 0; i < len(rows); i++ {
		index := strconv.Itoa(i)
		question, _ := file.GetCellValue(opt.SheetName, opt.Question+index)
		answer, _ := file.GetCellValue(opt.SheetName, opt.Answer+index)
		var options = make([]string, 0)
		for _, v := range opt.Option {
			option, _ := file.GetCellValue(opt.SheetName, v+index)
			options = append(options, option)
		}

		var as = make([]string, 0)
		if util.IsAlpha(answer) && len(options) > 0 {
			for _, a := range answer {
				if len(options) > int(a-65) {
					as = append(as, options[int(a-65)])
				}
			}
		} else {
			if util.IsAlpha(answer) {
				as = strings.Split(answer, "")
			} else {
				as = strings.Split(answer, "#")
			}
		}

		tiku := model.Question{
			Question: util.FormatString(question),
			Answer:   as,
			Options:  options,
		}
		if len(tiku.Answer) > 1 {
			tiku.Type = 1
		}
		if len(tiku.Answer) == 1 && (util.IsFalse(tiku.Answer[0]) || util.IsTrue(tiku.Answer[0])) {
			tiku.Type = 3
		}

		if len(tiku.Answer) == 1 && len(tiku.Options) == 0 {
			tiku.Type = 4
		}

		if answer != "" && tiku.Question != "" && len(tiku.Answer) > 0 {
			tikus = append(tikus, tiku)
		}
	}
	return tikus
}
