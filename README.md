# dun
> 网易易盾文本检测接口 for golang

### 安装方法

`go get github.com/jiangfengbing/dun`

### 使用方法

```go
import (
  "fmt"
  "github.com/jiangfengbing/dun"
)

const (
  secretID = "" // 产品秘钥ID ，由易盾反垃圾云服务分配
  secretKey = "" // 产品私钥
  businessID = "" // 业务ID，由易盾反垃圾云服务分配
)

func main() {
  checker := dun.NewChecker(secretID, secretKey, businessID)
  ret, err := checker.Check()
  if err != nil {
    fmt.Printf("Check error: %v", err)
    return
  }
  fmt.Printf("Check result: %d", ret)
}

```

