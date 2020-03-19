package reqHandle

import (
	"encoding/json"
	"errors"
	"fmt"
	"giftRanking/service"
	"giftRanking/util"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	parameter_error="Url Param 'userId' is missing"
)

var p = util.NewSimplePoll(10)

/*送礼*/
func SendGift(writer http.ResponseWriter, request *http.Request) {
	s, _ := ioutil.ReadAll(request.Body)
	var gift service.SendGift
	err := json.Unmarshal(s, &gift)
	// 根据请求body创建一个json解析器实例
	if err!=nil{
		fmt.Println(err)
	}
	p.Add(parseTask( &gift))
}

func parseTask( gift *service.SendGift) func() {
	return func() {
		service.SendGiftService(gift)
	}
}


/*送礼排行榜*/
func GetSortGift(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userId, ok := query["userId"]
	if !ok || len(userId[0]) < 1 {
		errors.New(parameter_error)
	}
	var s []byte = make([]byte, 5, 10)
	map1:=service.GetSortGift(userId[0])
	s, _ = json.Marshal(map1)
	writer.Write(s)
}

/*资金流水*/
func GiftList(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	userId, ok := query["userId"]
	pageNo, ok := query["pageNo"]
	pageSize, ok := query["pageSize"]
	var pn int64
	var ps int64
	if !ok || len(userId[0]) < 1 {
		errors.New(parameter_error)
	}
	if !ok || len(pageNo[0]) < 1 {
		pn=1
	}else{
		pn,_ = strconv.ParseInt(pageNo[0], 10, 64)
	}
	if !ok || len(pageSize[0]) < 1 {
		ps=10
	}else{
		ps,_ = strconv.ParseInt(pageSize[0], 10, 64)
	}
	var s1 []byte = make([]byte, 10, 20)
	s1,_  = json.Marshal(service.GiftList(userId[0],pn,ps))
	writer.Write(s1)
}