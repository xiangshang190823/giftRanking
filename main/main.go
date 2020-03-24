package main

import (
	"gift/reqHandle"
	"log"
	"net/http"
	_"net/http/pprof"
)

func main() {
	//设置访问的路由
	http.HandleFunc("/sendGift", reqHandle.SendGift) //送礼接口
	http.HandleFunc("/getSortGift", reqHandle.GetSortGift) //送礼排行榜
	http.HandleFunc("/giftList", reqHandle.GiftList) //送礼流水
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}





