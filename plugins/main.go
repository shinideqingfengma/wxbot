package fofa

import (
	"fmt"
	"github.com/imroc/req/v3"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	fofa_email = "962850765@qq.com"
	fofa_key   = "af3c11358202ae64d889dab2b8caa559"
)

func init() {
	engine := control.Register("fofa", &control.Options{
		Alias: "fofa查询",
		Help: "指令:\n" +
			"* fofa查询 [查询内容]",
	})

	engine.OnRegex(`^FoFa查询 ?(.*?)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		qbase := ctx.State["regex_matched"].([]string)[1]
		if data, err := searchFoFa(qbase); err == nil {
			if data == nil {
				ctx.ReplyText("没有查询到相关信息")
			} else {
				fileName := "fofa_result.xlsx"
				f := excelize.NewFile()
				sheetName := "Sheet1"
				f.NewSheet(sheetName)
				f.SetCellValue(sheetName, "A1", "IP地址")
				f.SetCellValue(sheetName, "B1", "端口")
				f.SetCellValue(sheetName, "C1", "协议")
				f.SetCellValue(sheetName, "D1", "主机")
				f.SetCellValue(sheetName, "E1", "域名")
				f.SetCellValue(sheetName, "F1", "操作系统")
				f.SetCellValue(sheetName, "G1", "服务器")
				f.SetCellValue(sheetName, "H1", "ICP备案")
				f.SetCellValue(sheetName, "I1", "标题")
				for i, item := range data {
					f.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), item["ip"])
					f.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), item["port"])
					f.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), item["protocol"])
					f.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), item["host"])
					f.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), item["domain"])
					f.SetCellValue(sheetName, "F"+strconv.Itoa(i+2), item["os"])
					f.SetCellValue(sheetName, "G"+strconv.Itoa(i+2), item["server"])
					f.SetCellValue(sheetName, "H"+strconv.Itoa(i+2), item["icp"])
					f.SetCellValue(sheetName, "I"+strconv.Itoa(i+2), item["title"])
				}
				err = f.SaveAs(fileName)
				if err != nil {
					fmt.Println("保存文件失败", err)
				}
				ctx.ReplyFile(fileName)
			}
		} else {
			ctx.ReplyText("查询失败，这一定不是bug🤔")
		}
	})
}

type fofaResponse struct {
	Data []map[string]string `json:"data"`
}

func searchFoFa(qbase string) ([]map[string]string, error) {
	var data fofaResponse
	client := req.C().
		SetQueryParams(map[string]string{
			"email":   fofa_email,
			"key":     fofa_key,
			"qbase64": qbase,
			"fields":  "ip,port,protocol,host,domain,os,server,icp,title",
			"size":    "1000",
		}).
		SetHeader("Accept", "application/json")

	if err := client.Get("https://fofa.info/api/v1/search/all").Do().Into(&data); err != nil {
		return nil, err
	}
	if len(data.Data) == 0 {
		return nil, nil
	}
	return data.Data, nil
}
