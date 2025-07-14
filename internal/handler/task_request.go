package handler

type ProcessRequest struct {
	Limit       int     `query:"limit" validate:"required,min=1,max=100"`
	MinDelay    int     `query:"minDelay" validate:"min=0"`
	MaxDelay    int     `query:"maxDelay" validate:"gtefield=MinDelay"`
	SuccessRate float64 `query:"successRate" validate:"required,gte=0,lte=1"`
}

type TaskQueryParams struct {
	Page     int    `query:"page" validate:"gte=1"`
	PageSize int    `query:"pageSize" validate:"gte=1,lte=100"`
	Status   string `query:"status" validate:"omitempty,oneof=NEW PROCESSING PROCESSED"`
}
