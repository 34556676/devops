package api

import (
	"devops/common"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Namespace struct {
}

// 获取K8S命名空间
func (n Namespace) GetNamesapce(c *gin.Context) {

	clientset, err := common.InitClient()
	if err != nil {
		fmt.Println(err)
	}
	namespaceList, _ := clientset.CoreV1().Namespaces().List(c, v1.ListOptions{})

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "success",
		"data": namespaceList.Items,
	})
}
