// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/12/15 21:27'
//
// Usage:
//

package requests_handler

import (
	"encoding/json"
	"fmt"
	"github.com/asmcos/requests"
	"video_client/utils/api_error"
)

//
//  @Description: 响应内容检查,并提取数据
//  @param resp:
//  @return respData:
//  @return err:
//
func ResponseCheck(resp *requests.Response) (respData map[string]interface{}, err error) {
	if resp.R.StatusCode != 200 {
		return respData, api_error.New(500, fmt.Sprintf("接口失败,http状态码为:%v", resp.R.StatusCode))
	}

	body := resp.Text()

	err = json.Unmarshal([]byte(body), &respData)
	if err != nil {
		return respData, api_error.New(500, "响应body转为json失败")
	}

	respCode := int(respData["respCode"].(float64))
	respMsg := respData["respMsg"].(string)
	if respCode != 200 {
		return respData, api_error.New(500, fmt.Sprintf("响应状态码为%v,错误内容为%v", respCode, respMsg))
	}
	return respData, nil
}
