package routers

import (
	"devops/api"
	"github.com/gin-gonic/gin"
)

func KubernetesGroupInit(r *gin.Engine) {
	kubernetesGroup := r.Group("/k8s")
	{
		kubernetesGroup.GET("/namespace", api.Namespace{}.GetNamesapce)

	}
}
