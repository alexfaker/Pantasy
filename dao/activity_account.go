package dao

type Account struct {
	ID              int64  `gorm:"id"`
	UserID          string `gorm:"user_id"`           //用户ID
	CashDevice      int    `gorm:"cash_device"`       //红包金额(单位分)
	CashInviteTotal int    `gorm:"cash_invite_total"` //邀请赠送金额(单位分)
	CashExtract     int    `gorm:"cash_extract"`      //已经提取金额(单位分)
	Status          int    `gorm:"status"`            //状态：0-正常 1-重复注册
	UpdatedAt       int64  `gorm:"updated_at"`        //完成时间
	DeletedAt       int64  `gorm:"deleted_at"`        //删除时间
	CreatedAt       int64  `gorm:"created_at"`        //创建时间
}

func (t *Account) TableName() string {
	return "activity_account"
}

func (t *Account) Insert() error {
	return nil
}
