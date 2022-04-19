package model

type User struct {
	ID            uint      `json:"id"`
	Username      string    `json:"username"` // 用户名
	Phone         string    `json:"phone"`    // 手机号码
	//Nickname      string    `json:"nickname"` // 昵称
	//Email         string    `json:"email"`    // 邮箱
	//Password      string    `json:"password"` // 密码
	//IsDel         uint      `json:"is_del"`
	//RememberToken string    `json:"remember_token"`
	//ApiToken      string    `json:"api_token"`
	//CreatedAt     time.Time `json:"created_at"`
	//UpdatedAt     time.Time `json:"updated_at"`
	//CreatedTime   int `gorm:"autoCreateTime"`
}

//TableName 重写表名
func (User) TableName() string {
	return "union_users"
}

func NewUser() *User {
	u  := &User{}

	return u
}

func (u *User) Info() *User {
	//where := "id = 20"
	//user := Demo{}

	/*us2 := &User{}
	Find(us2, "id=10", "")*/

	/*us := User{
		Username: "tcl",
	}
	Updates(us, "id=13")*/

	//GetDB().Where("id = ?", "10").First(&User{})
	//GetDB().Where("id = ?", "5").First(&User{})

	Find(u, "id=5", "")

	return u
}

