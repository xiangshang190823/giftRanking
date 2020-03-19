package user_crud

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type User struct {
	UserName string `bson:"user_name"`
	UserId string `bson:"user_id"`
	Status int  `bson:"status"` //0 主播  1 用户
	UserCoin float32 `bson:"user_coin"`
}

type LogMgr struct {
	client *mongo.Client
	collection *mongo.Collection
}

var (
	G_logMgr *LogMgr
	Client *mongo.Client
)

func InItMongodb()  {

	var(
		ctx context.Context
		opts *options.ClientOptions
		client *mongo.Client
		err error
		collection *mongo.Collection
	)
	// 连接数据库
	ctx, _ = context.WithTimeout(context.Background(), time.Duration(2000) * time.Millisecond)  // ctx
	opts = options.Client().ApplyURI("mongodb://localhost:27017").SetMaxPoolSize(10)  // opts
	if client, err = mongo.Connect(ctx,opts); err != nil{
		log.Fatal(err)
		return
	}

	//链接数据库和表
	collection = client.Database("test").Collection("user")

	//单例
	G_logMgr = &LogMgr{
		client:client,
		collection:collection,
	}
}

//保存数据
func (logMgr *LogMgr)SaveMongodb() (err error) {
	var(
		insetRest *mongo.InsertOneResult
		id interface{}
		users []interface{}
	)
	//userId := bson.NewObjectId().String()
	user := User{"主播", "1", 0,0}
	if insetRest, err = logMgr.collection.InsertOne(context.TODO(), &user); err != nil{
		fmt.Println(err)
		return
	}
	id = insetRest.InsertedID
	fmt.Println(id)

	users = append(users, &User{UserName:"用户1", UserId:"10001",Status:1,UserCoin:5102},
	&User{UserName:"用户2", UserId:"10002",Status:1,UserCoin:84263},
		&User{UserName:"用户3", UserId:"10003",Status:1,UserCoin:1520},
		&User{UserName:"用户4", UserId:"10004",Status:1,UserCoin:358})
	if _, err = logMgr.collection.InsertMany(context.TODO(), users); err != nil{
		log.Fatal(err)
		return
	}
	return
}


//查询数据
func (logMgr *LogMgr)SelectMongodbByUserId(userId string)(err error,user *User)  {
	if err = logMgr.collection.FindOne( context.TODO(), bson.M{"user_id":userId}).Decode(&user); err != nil{
		errors.New("数据库查询失败")
		log.Fatal(err)
	}
	return
}


//更新数据
func (logMgr *LogMgr)UpdateCoin(user User,coin float32) (err error)  {
	var(
		ctx context.Context
	)
	if singleResult := logMgr.collection.FindOneAndUpdate(ctx, bson.M{"user_id": user.UserId,"user_coin":user.UserCoin},
	bson.M{"$inc": bson.M{"user_coin": coin}}); singleResult.Err() != nil{
		errors.New("数据库更新失败")
		log.Fatal(err)
		return
	}
	return
}


//删除数据
func (logMgr *LogMgr)DeleteMongo()(err error)  {
	var(
		ctx context.Context
	)
	if _, err = logMgr.collection.DeleteMany(ctx, bson.M{"age":bson.M{"$gte":10}}); err != nil{
		log.Fatal(err)
		return
	}
	return
}
