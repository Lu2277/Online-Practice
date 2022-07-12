# 基于Gin、Gorm实现的在线练习系统

## 数据表

> 1.问题表problem
>
>2.用户表user
>
> 3.提交表submit
>

## 用户注册

1.填写基本信息、发送验证码

```go
// CreateRandCode  生成随机验证码
func CreateRandCode() string {
rand.Seed(time.Now().UnixNano())
code := ""
for i := 0; i < 6; i++ {
randNum := rand.Intn(10)
code = code + strconv.Itoa(randNum)
}
return code
}
```

```go
// SendCode 发送验证码
func SendCode(toEmail, code string) error {
em := email.NewEmail()
// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
em.From = "LU <2591134973@qq.com>"
// 设置 receiver 接收方的邮箱
em.To = []string{toEmail}
// 设置主题
em.Subject = "发送验证码"
// 简单设置文件发送的内容，暂时设置成纯文本
em.HTML = []byte("本次验证码为：<b>" + code + "</b>")
//设置服务器相关的配置
err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "2591134973@qq.com", "uhpitgbjrkirdjed", "smtp.qq.com"))
if err != nil {
log.Fatal(err)
}
log.Println("send successfully ... ")
return nil
}
```

2.检查验证码是否正确

```go
    //判断验证码是否正确
trueCode, err := models.RDB.Get(c, email).Result()
if err != nil {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "获取验证码错误",
})
return
}
if trueCode != code {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "验证码错误",
})
return
}
```

3.检查用户是否已注册

```go
    //判断邮箱是否已经注册
var num int64
err = models.DB.Where("email= ?", email).Model(&models.User{}).Count(&num).Error
if err != nil {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "获取用户数量错误：" + err.Error(),
})
return
}
if num > 0 {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "该邮箱已经注册了",
})
return
}
```

4.用户数据插入到user数据库

```go
    //数据插入
data := &models.User{
Identity: helper.GetUUID(), //唯一标识码uuid
Name:     name,
Password: helper.GetMd5(password), //密码加密为md5格式
Phone:    phone,
Email:    email,
}
//插入到数据库
err = models.DB.Create(data).Error
```

5.生成用户的token（jwt鉴权）

## JWT鉴权

1.生成token(包含用户的identity、name和是否为管理员等信息)

```go
var myKey = []byte("project - key")
// GenerateToken 生成token
func GenerateToken(identity string, name string, isAdmin int) (string, error) {
UserClaim := &UserClaims{
Identity:       identity,
Name:           name,
IsAdmin:        isAdmin,
StandardClaims: jwt.StandardClaims{},
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
signedString, err := token.SignedString(myKey)
if err != nil {
return "", err
}
//fmt.Println(signedString)
return signedString, nil
}
```

2.解析token

```go
// AnalyseToken 解析token
func AnalyseToken(signedString string) (*UserClaims, error) {
userClaim := new(UserClaims)
claims, err := jwt.ParseWithClaims(signedString, userClaim, func(token *jwt.Token) (interface{}, error) {
return myKey, nil
})
if err != nil {
return nil, err
}
if !claims.Valid {
return nil, fmt.Errorf("anlyse Token Error:%v", err)
}
//fmt.Println(userClaim)
return userClaim, nil
}
```

## 用户登录

1.登陆验证

```go
func Login(c *gin.Context) {
username := c.PostForm("username")
password := c.PostForm("password")
if username == "" || password == "" {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "用户名、密码不能为空",
})
return
}
//获取密码的md5格式
password = helper.GetMd5(password)
//print(username, password)
data := new(models.User)
err := models.DB.Where("name= ? AND password = ?", username, password).First(&data).Error
if err != nil {
if err == gorm.ErrRecordNotFound {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "用户名或者密码错误",
})
return
} else {
c.JSON(http.StatusOK, gin.H{
"code": -1,
"msg":  "Get User Error:" + err.Error(),
})
return
}
}
```

2.生成登录用户的token

```go
    token, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
```

## 整合 Swagger

参考文档： https://github.com/swaggo/gin-swagger
接口访问地址：http://localhost:8080/swagger/index.html

```text
// GetProblem
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
```

