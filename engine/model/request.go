package model

type RunContainerRequestBody struct {
	TemplateName string                 `json:"templateName" binding:"required"` // 模版名
	Param        map[string]interface{} `json:"param" binding:"required"`        // 函数参数
	Synchronized bool                   `json:"synchronized"`                    // 同步等待结果
}
