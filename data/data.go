package data

// 更新账户最小信息Update Minimum Information
type UpdateMI struct {
	UpdateTimeStamp string
	InstID          string
	Value           float64
}

func NewUpdateMI(time string, instID string, value float64) UpdateMI {
	return UpdateMI{
		UpdateTimeStamp: time,
		InstID:          instID,
		Value:           value,
	}

}
