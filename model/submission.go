package model

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

//hint: 如果你想直接返回结构体，可以考虑在这里加上`json`的tag
type Submission struct {
	ID        uint   `gorm:"not null;autoIncrement"`
	UserName  string `gorm:"type:varchar(255);"`
	Avatar    string //头像base64，也可以是一个头像链接
	CreatedAt int64  //提交时间
	Score     int    //评测成绩
	Sub1      int
	Sub2      int
	Sub3      int
}

//这里提供返回的submission的示例结构
type ReturnSub struct {
	UserName  string `json:"user"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"time"`
	Score     int    `json:"score"`
	Subs      [3]int `json:"subs"`
	Votes     uint   `json:"votes"`
}

type ReturnSubSlice []ReturnSub

func (a ReturnSubSlice) Len() int {
	return len(a)
}
func (a ReturnSubSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ReturnSubSlice) Less(i, j int) bool {
	return a[j].Score < a[i].Score
}

/*TODO: 添加相应的与数据库交互逻辑，补全参数和返回值，可以参考user.go的设计思路*/

func Judge(content string) (error, int, [3]int) {
	file, err := os.OpenFile("ground_truth.txt", os.O_RDONLY, 0666)
	if err != nil {
		return err, 0, [3]int{0, 0, 0}
	}
	defer file.Close()
	result := [3]int{0, 0, 0}
	reader := bufio.NewReader(file)
	answers := strings.Split(content, "\n")
	_, err = reader.ReadString('\n')
	if err != nil {
		return err, 0, [3]int{0, 0, 0}
	}
	nowline := 0
	for {
		if nowline >= len(answers) {
			err = errors.New("invalid input")
			return err, 0, [3]int{0, 0, 0}
		}
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		answer := strings.Split(answers[nowline], ",")
		correctAnswer := strings.Split(str, ",")
		if len(answer) != 3 {
			err = errors.New("invalid input")
			return err, 0, [3]int{0, 0, 0}
		}
		for i := 0; i < 3; i++ {
			if answer[i] == "1" && (correctAnswer[i+1] == "True" || correctAnswer[i+1] == "True\n") {
				result[i]++
			}
			if answer[i] == "0" && (correctAnswer[i+1] == "False" || correctAnswer[i+1] == "False\n") {
				result[i]++
			}
		}
		nowline++
	}
	print(result[0], " ", result[1], " ", result[2])
	return nil, (result[0] + result[1] + result[2]) / 30, result
}

func CreateSubmission(name string, avatar string, content string) error {
	err, _ := GetUserByName(name)
	if err != nil {
		err, _ = CreateUser(name)
		if err != nil {
			return err
		}
	}
	err, score, subs := Judge(content)
	if err != nil {
		return err
	}
	submission := Submission{
		UserName:  name,
		Avatar:    avatar,
		CreatedAt: time.Now().Unix(),
		Score:     score,
		Sub1:      subs[0],
		Sub2:      subs[1],
		Sub3:      subs[2],
	}
	print(submission.CreatedAt)
	tx := DB.Create(&submission)
	//我怎么知道插入到了哪一个table里面？
	return tx.Error
}

func GetUserSubmissions(username string) (error, []ReturnSub) {
	//返回某一用户的所有提交
	//在查询时可以使用.Order()来控制结果的顺序，详见https://gorm.io/zh_CN/docs/query.html#Order
	//当然，也可以查询后在这个函数里手动完成排序
	var TempSub []Submission
	var RetSub []ReturnSub
	tx := DB.Model(&Submission{}).Where("user_name=?", username).Order("created_at Desc").Find(&TempSub)

	for i := 0; i < len(TempSub); i++ {
		RetSub = append(RetSub, ReturnSub{
			UserName:  TempSub[i].UserName,
			Avatar:    TempSub[i].Avatar,
			CreatedAt: TempSub[i].CreatedAt,
			Score:     TempSub[i].Score,
			Subs:      [3]int{TempSub[i].Sub1, TempSub[i].Sub2, TempSub[i].Sub3},
		})
	}
	return tx.Error, RetSub
}

func GetLeaderBoard() []ReturnSub {
	//一个可行的思路，先全部选出submission，然后手动选出每个用户的最后一次提交
	var AllSub []Submission
	var TempSub []ReturnSub
	var RetSub []ReturnSub
	var user User
	DB.Model(&Submission{}).Where("1=1").Order("created_at Desc").Find(&AllSub)
	for i := 0; i < len(AllSub); i++ {
		//DB.First(&user, "user_name=?", AllSub[i].UserName)
		TempSub = append(TempSub, ReturnSub{
			UserName:  AllSub[i].UserName,
			Avatar:    AllSub[i].Avatar,
			CreatedAt: AllSub[i].CreatedAt,
			Score:     AllSub[i].Score,
			Subs:      [3]int{AllSub[i].Sub1, AllSub[i].Sub2, AllSub[i].Sub3},
		})
	}

	var users map[string]int
	users = make(map[string]int)
	for i := 0; i < len(AllSub); i++ {
		_, ok := users[TempSub[i].UserName]
		if ok {

		} else {
			users[AllSub[i].UserName] = 1

			_, user = GetUserByName(AllSub[i].UserName)
			TempSub[i].Votes = user.Votes
			RetSub = append(RetSub, TempSub[i])
		}
	}
	sort.Sort(ReturnSubSlice(RetSub))
	return RetSub
}
