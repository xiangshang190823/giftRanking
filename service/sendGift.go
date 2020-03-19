package service

import (
	"bytes"
	"errors"
	"fmt"
	"giftRanking/model/coinHistory/coinHistory_crud"
	"giftRanking/model/user/user_crud"
	"giftRanking/util"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)
const (
	gift="gift"
	underline="_"
	dist_lock = "distlock"
)

type SendGift struct {
	UserId    string `json:"userId"`
	GetUserId string `json:"getUserId"`
	Coin      float32 `json:"coin"`
}
type a []UserCoin
type Ret struct{
	Code int
	Param string
	Msg string
	TotalCount int64
	NowPageNo int64
	NowPageSize int64
	Data a
}
type UserCoin struct {
	UserName string
	CostCoin float32
	UserId string
}
func SendGiftService(send *SendGift) (ret Ret){
	// 从池里获取连接
	rc := util.RedisClient.Get()
	pidstr := strconv.Itoa(os.Getpid())
	//分布式锁
	lockName:=getString(dist_lock,send.UserId,send.GetUserId)
	_, err := redis.String(rc.Do("GET", lockName))
	if err == nil {
		//fmt.Println(err)
		return
	}
	rc.Do("SETNX",lockName,pidstr)
	rc.Do("EXPIRE", lockName, 10)
	//1.查询送礼用户的信息，并查看这个金额是否大于他自己的金额
	user_crud.InItMongodb()
	error,user:=user_crud.G_logMgr.SelectMongodbByUserId(send.UserId)
	if error!=nil{
		log.Fatalln("user does not exist:", error)
		errors.New("user does not exist")
	}
	err,toUser:=user_crud.G_logMgr.SelectMongodbByUserId(send.GetUserId)
	if err!=nil{
		log.Fatalln("toUser does not exist:", err)
		errors.New("toUser does not exist")
	}

	//1.1.如果大于，则扣减需要送礼的用户的金币，否则提示失败
	if send.Coin>user.UserCoin{
		log.Fatalln("您的金币不够哟!请充值后再赠送", err)
		errors.New("您的金币不够哟!请充值后再赠送")
	}
	//1.2用户赠送，删减金币数
	err1:=user_crud.G_logMgr.UpdateCoin(*user,-send.Coin)
	if err1!=nil{
		errors.New("用户删除金币失败")
	}
	//2.增加金币数
	err2:=user_crud.G_logMgr.UpdateCoin(*toUser,send.Coin)
	if err2!=nil{
		user_crud.G_logMgr.UpdateCoin(*user,send.Coin)
	}

	//2.5生成对应的redis key，规则是gift_主播Id
	key:=strings.Join([]string{gift,toUser.UserId}, underline)

	//3.增加流水记录
	coinHistory_crud.InItMongodb()
	coinHistory:=new(coinHistory_crud.CoinHistory)
	coinHistory.UserId=user.UserId
	coinHistory.Coin=send.Coin
	coinHistory.GetUserId=toUser.UserId
	//历史记录ID:来源+用户ID+日期时间戳
	coinHistory.HistoryId=getString(gift,user.UserId,strconv.FormatInt(time.Now().Unix(),10))
	coinHistory.Time=time.Now()
	err3:=coinHistory_crud.G_logMgr.SaveCoinHistory(coinHistory)

	if err3!=nil{
		user_crud.G_logMgr.UpdateCoin(*toUser,-send.Coin)
		user_crud.G_logMgr.UpdateCoin(*user,send.Coin)
	}
	rc.Do("Del",dist_lock+send.UserId+send.GetUserId)
	//4.增加redis zset数据
	sign:=getString(user.UserName,underline,user.UserId)

	if _,e:=rc.Do("ZREVRANK",key,user.UserId); e != nil{
		rc.Do("zadd",key,send.Coin,sign)
	}else{
		rc.Do("ZIncrBy",key,send.Coin,sign)
	}
	// 用完后将连接放回连接池
	defer rc.Close()
	ret.Code=0
	ret.Msg="success"
	return ret
}

func GetSortGift(userId string)( ret Ret){
	var usercoins *UserCoin
	var buffer bytes.Buffer
	buffer.WriteString(gift)
	buffer.WriteString(underline)
	buffer.WriteString(userId)
	u:=buffer.String()
	rc := util.RedisClient.Get()
	user_map, err :=redis.StringMap(rc.Do("zrevrange", u, 0, -1, "withscores"))
	if err != nil {
		errors.New("redis get failed")
	}
	usercoins = &UserCoin{}
	var i int64 =0
	for user := range user_map {
		fmt.Printf("user name: %v %v\n", user,user_map[user] )
		usercoins.UserName=strings.Split(user,underline)[0]
		usercoins.UserId=strings.Split(user,underline)[1]
		v1, _ := strconv.ParseFloat(user_map[user], 32)
		usercoins.CostCoin= float32(v1)
		ret.Data = append(ret.Data, *usercoins)
		i++
	}
	ret.TotalCount=i

	sort.Stable(ret.Data)
	return ret
}

func GiftList(userId string,pageNo int64,pageSize int64) interface{}{
	coinHistory_crud.InItMongodb()
	return coinHistory_crud.G_logMgr.SelectMongodb(userId,pageNo,pageSize)
}


func (s a) Len() int { return len(s) }

func (s a) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s a) Less(i, j int) bool { return s[i].CostCoin > s[j].CostCoin }

func getString(a string,b string,c string)  string{
	var buffer bytes.Buffer
	buffer.WriteString(a)
	buffer.WriteString(b)
	buffer.WriteString(c)
	return buffer.String()
}