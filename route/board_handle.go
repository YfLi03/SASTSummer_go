package route

import (
	"github.com/gin-gonic/gin"
	"leadboard/model"
	"net/http"
)

//TODO:在这里完成handle function，返回所有的leader board内容
func HandleGetBoard(g *gin.Context) {
	g.JSON(http.StatusAccepted, model.GetLeaderBoard())
}

//TODO:在这里完成返回一个用户提交历史的Handle function
func HandleUserHistory(g *gin.Context) {
	name := g.Param("name")
	err, ret := model.GetUserSubmissions(name)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
		})
	} else {
		g.JSON(http.StatusAccepted, gin.H{
			"code": 0,
			"data": ret,
		})
	}
}

//TODO:在这里完成接受提交内容，进行评判的handle function
func HandleSubmit(g *gin.Context) {
	type SubmitForm struct {
		UserName string `json:"user"`
		Avatar   string `json:"avatar"`
		Content  string `json:"content"`
	}
	var form SubmitForm
	if err := g.ShouldBindJSON(&form); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "参数不全啊",
		})
	}

	if len(form.UserName) > 255 {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "用户名太长了",
		})
	}

	if len(form.Avatar) > 10240 {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -2,
			"msg":  "图像太大了",
		})
	}

	err := model.CreateSubmission(form.UserName, form.Avatar, form.Content)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{
			"code": -3,
			"msg":  "非法内容呜呜",
		})
	} else {
		type Data struct {
			Leaderboard []model.ReturnSub `json:"leaderboard"`
		}
		data := Data{
			Leaderboard: model.GetLeaderBoard(),
		}
		g.JSON(http.StatusAccepted, gin.H{
			"code": 0,
			"msg":  "提交成功",
			"data": data,
		})
	}

}
