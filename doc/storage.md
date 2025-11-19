# 云存储集成文档

ClarkGo CMS 提供了灵活的云存储解决方案，支持本地文件系统、阿里云 OSS 和 AWS S3 三种存储方式。

## 目录

- [功能特性](#功能特性)
- [存储驱动](#存储驱动)
- [配置方法](#配置方法)
- [使用示例](#使用示例)
- [API 参考](#api-参考)
- [最佳实践](#最佳实践)

## 功能特性

### 统一接口
- 所有存储驱动实现相同的 `Storage` 接口
- 零代码改动即可切换存储方式
- 支持运行时动态配置

### 三种驱动
1. **Local** - 本地文件系统存储
   - 适合开发环境
   - 无需额外配置
   - 即开即用

2. **OSS** - 阿里云对象存储
   - 高可用、高并发
   - 支持 CDN 加速
   - 内网上传节省流量

3. **S3** - AWS 对象存储
   - 全球可用
   - 支持 CloudFront CDN
   - 与 AWS 生态无缝集成

### 核心功能
- ✅ 文件上传
- ✅ 文件删除
- ✅ 文件检查
- ✅ URL 生成
- ✅ 文件大小获取
- ✅ 签名 URL（私有访问）

## 存储驱动

### Storage 接口

所有存储驱动都实现了以下接口：

```go
type Storage interface {
    Save(file io.Reader, path string) error
    Delete(path string) error
    Exists(path string) bool
    URL(path string) string
    Size(path string) (int64, error)
}
```

### 本地存储 (Local)

**优点：**
- 无需配置
- 无外部依赖
- 适合开发测试

**缺点：**
- 扩展性受限
- 无法跨服务器共享
- 需要自行管理磁盘空间

**使用场景：**
- 开发环境
- 小型项目
- 单服务器部署

### 阿里云 OSS

**优点：**
- 中国区速度快
- 价格实惠
- 支持内网传输

**缺点：**
- 国际访问速度较慢
- 需要备案（公网访问）

**使用场景：**
- 中国区业务
- 高并发场景
- 需要 CDN 加速

### AWS S3

**优点：**
- 全球节点覆盖
- 稳定性极高
- 生态完善

**缺点：**
- 中国区访问受限
- 价格较高

**使用场景：**
- 国际业务
- 与 AWS 服务集成
- 需要全球访问

## 配置方法

### 1. 环境变量配置

复制 `.env.example` 到 `.env`，根据需要配置：

#### 本地存储配置

```env
STORAGE_DRIVER=local
LOCAL_STORAGE_PATH=./storage/uploads
LOCAL_STORAGE_URL=/uploads
```

#### 阿里云 OSS 配置

```env
STORAGE_DRIVER=oss
OSS_ENDPOINT=oss-cn-hangzhou.aliyuncs.com
OSS_ACCESS_KEY_ID=your_access_key_id
OSS_ACCESS_KEY_SECRET=your_access_key_secret
OSS_BUCKET_NAME=your_bucket_name
OSS_BASE_URL=https://your-cdn-domain.com
```

**OSS Endpoint 列表：**
- 杭州：`oss-cn-hangzhou.aliyuncs.com`
- 北京：`oss-cn-beijing.aliyuncs.com`
- 上海：`oss-cn-shanghai.aliyuncs.com`
- 深圳：`oss-cn-shenzhen.aliyuncs.com`
- 香港：`oss-cn-hongkong.aliyuncs.com`

#### AWS S3 配置

```env
STORAGE_DRIVER=s3
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=your_access_key_id
S3_SECRET_ACCESS_KEY=your_secret_access_key
S3_BUCKET_NAME=your_bucket_name
S3_BASE_URL=https://your-cloudfront-domain.com
```

**S3 Region 列表：**
- 美东 (弗吉尼亚)：`us-east-1`
- 美西 (俄勒冈)：`us-west-2`
- 亚太 (新加坡)：`ap-southeast-1`
- 亚太 (东京)：`ap-northeast-1`
- 欧洲 (法兰克福)：`eu-central-1`

### 2. 代码中使用

```go
package main

import (
    "github.com/chenyu/clarkgo/config"
)

func main() {
    // 获取存储实例（根据环境变量自动选择）
    storage, err := config.GetStorage()
    if err != nil {
        panic(err)
    }

    // 使用存储（无论是什么驱动，代码完全一致）
    err = storage.Save(file, "images/avatar.jpg")
}
```

## 使用示例

### 基础上传

```go
package controllers

import (
    "github.com/chenyu/clarkgo/config"
    "github.com/cloudwego/hertz/pkg/app"
)

func UploadFile(c *app.RequestContext) {
    // 获取上传文件
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, map[string]interface{}{
            "error": "file is required",
        })
        return
    }

    // 打开文件
    src, err := file.Open()
    if err != nil {
        c.JSON(500, map[string]interface{}{
            "error": "failed to open file",
        })
        return
    }
    defer src.Close()

    // 获取存储实例
    storage, err := config.GetStorage()
    if err != nil {
        c.JSON(500, map[string]interface{}{
            "error": "storage not available",
        })
        return
    }

    // 保存文件
    path := "uploads/" + file.Filename
    if err := storage.Save(src, path); err != nil {
        c.JSON(500, map[string]interface{}{
            "error": "failed to save file",
        })
        return
    }

    // 返回文件URL
    c.JSON(200, map[string]interface{}{
        "url": storage.URL(path),
    })
}
```

### 删除文件

```go
func DeleteFile(c *app.RequestContext) {
    path := c.Query("path")
    
    storage, _ := config.GetStorage()
    
    if !storage.Exists(path) {
        c.JSON(404, map[string]interface{}{
            "error": "file not found",
        })
        return
    }
    
    if err := storage.Delete(path); err != nil {
        c.JSON(500, map[string]interface{}{
            "error": "failed to delete file",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "message": "deleted successfully",
    })
}
```

### 获取文件信息

```go
func FileInfo(c *app.RequestContext) {
    path := c.Query("path")
    
    storage, _ := config.GetStorage()
    
    if !storage.Exists(path) {
        c.JSON(404, map[string]interface{}{
            "error": "file not found",
        })
        return
    }
    
    size, _ := storage.Size(path)
    
    c.JSON(200, map[string]interface{}{
        "path": path,
        "url":  storage.URL(path),
        "size": size,
    })
}
```

### 私有文件访问（签名 URL）

```go
// 注意：Storage 接口未定义 SignURL，需要类型断言

func GetPrivateFileURL(c *app.RequestContext) {
    path := c.Query("path")
    
    storage, _ := config.GetStorage()
    
    // 检查是否支持签名URL
    type Signer interface {
        SignURL(path string, expireSeconds int64) (string, error)
    }
    
    signer, ok := storage.(Signer)
    if !ok {
        // 不支持签名（本地存储）
        c.JSON(200, map[string]interface{}{
            "url": storage.URL(path),
        })
        return
    }
    
    // 生成1小时有效的签名URL
    signedURL, err := signer.SignURL(path, 3600)
    if err != nil {
        c.JSON(500, map[string]interface{}{
            "error": "failed to generate signed URL",
        })
        return
    }
    
    c.JSON(200, map[string]interface{}{
        "url":        signedURL,
        "expires_in": 3600,
    })
}
```

## API 参考

### Storage 接口方法

#### Save

```go
Save(file io.Reader, path string) error
```

上传文件到存储。

**参数：**
- `file`: 文件内容（实现了 io.Reader 接口）
- `path`: 存储路径（相对路径）

**返回：**
- `error`: 错误信息

#### Delete

```go
Delete(path string) error
```

删除文件。

**参数：**
- `path`: 文件路径

**返回：**
- `error`: 错误信息

#### Exists

```go
Exists(path string) bool
```

检查文件是否存在。

**参数：**
- `path`: 文件路径

**返回：**
- `bool`: true 表示存在

#### URL

```go
URL(path string) string
```

获取文件访问 URL。

**参数：**
- `path`: 文件路径

**返回：**
- `string`: 完整访问 URL

#### Size

```go
Size(path string) (int64, error)
```

获取文件大小。

**参数：**
- `path`: 文件路径

**返回：**
- `int64`: 文件大小（字节）
- `error`: 错误信息

### 额外方法（OSS/S3）

#### SignURL

```go
SignURL(path string, expireSeconds int64) (string, error)
```

生成签名 URL（用于私有文件临时访问）。

**参数：**
- `path`: 文件路径
- `expireSeconds`: 过期时间（秒）

**返回：**
- `string`: 签名 URL
- `error`: 错误信息

## 最佳实践

### 1. 开发与生产环境分离

```env
# 开发环境 (.env.dev)
STORAGE_DRIVER=local

# 生产环境 (.env.prod)
STORAGE_DRIVER=oss
```

### 2. CDN 加速

配置 `BASE_URL` 使用 CDN 域名：

```env
# OSS
OSS_BASE_URL=https://cdn.example.com

# S3
S3_BASE_URL=https://d123456.cloudfront.net
```

### 3. 私有与公有文件分离

```go
// 公有文件 - 使用 public/ 前缀
storage.Save(file, "public/images/avatar.jpg")

// 私有文件 - 使用 private/ 前缀
storage.Save(file, "private/documents/invoice.pdf")

// 配置 Bucket 策略：
// - public/* 允许公开访问
// - private/* 仅签名访问
```

### 4. 文件路径规划

```go
// 按类型分类
"images/products/123.jpg"
"documents/contracts/456.pdf"
"videos/tutorials/789.mp4"

// 按日期分类
"uploads/2024/01/15/avatar.jpg"

// 按用户分类
"users/1001/avatar.jpg"
"users/1001/documents/resume.pdf"
```

### 5. 错误处理

```go
storage, err := config.GetStorage()
if err != nil {
    // 记录日志
    log.Error("Storage initialization failed:", err)
    
    // 降级处理
    return fallbackHandler(c)
}
```

### 6. 性能优化

```go
// 1. 并发上传多个文件
var wg sync.WaitGroup
for _, file := range files {
    wg.Add(1)
    go func(f *multipart.FileHeader) {
        defer wg.Done()
        storage.Save(f, path)
    }(file)
}
wg.Wait()

// 2. 使用流式上传（大文件）
reader := bufio.NewReader(file)
storage.Save(reader, path)

// 3. 检查文件是否已存在（避免重复上传）
if storage.Exists(path) {
    return storage.URL(path), nil
}
```

### 7. 安全建议

```go
// 1. 文件类型验证
allowedTypes := []string{".jpg", ".png", ".pdf"}
ext := filepath.Ext(filename)
if !contains(allowedTypes, ext) {
    return errors.New("invalid file type")
}

// 2. 文件大小限制
maxSize := int64(10 * 1024 * 1024) // 10MB
if size > maxSize {
    return errors.New("file too large")
}

// 3. 文件名清理
safeName := sanitizeFilename(originalName)
path := "uploads/" + safeName

// 4. 使用随机文件名
uuid := generateUUID()
path := fmt.Sprintf("uploads/%s%s", uuid, ext)
```

## 迁移指南

### 从本地存储迁移到 OSS/S3

1. **准备工作：**
   ```bash
   # 安装阿里云 OSS 工具
   wget http://gosspublic.alicdn.com/ossutil/1.7.13/ossutil64
   chmod +x ossutil64
   
   # 或安装 AWS CLI
   pip install awscli
   ```

2. **批量上传：**
   ```bash
   # OSS
   ./ossutil64 cp -r ./storage/uploads oss://your-bucket/
   
   # S3
   aws s3 sync ./storage/uploads s3://your-bucket/
   ```

3. **更新配置：**
   ```env
   STORAGE_DRIVER=oss  # 或 s3
   ```

4. **验证迁移：**
   ```bash
   # 检查文件数量
   ./ossutil64 ls oss://your-bucket/ -r | wc -l
   
   # 或
   aws s3 ls s3://your-bucket/ --recursive | wc -l
   ```

## 故障排查

### 1. OSS 403 错误

**原因：** AccessKey 权限不足

**解决：**
- 检查 RAM 用户是否有 OSS 读写权限
- 确认 Bucket 策略配置正确

### 2. S3 连接超时

**原因：** 网络问题或 Region 错误

**解决：**
- 检查 `S3_REGION` 是否正确
- 使用 VPC Endpoint（AWS 内网）
- 配置代理

### 3. 文件 URL 无法访问

**原因：** Bucket 未公开或 CDN 未配置

**解决：**
- OSS: 设置 Bucket 为公共读
- S3: 配置 Bucket Policy 允许 GetObject
- 检查 `BASE_URL` 配置是否正确

### 4. 上传失败

**原因：** 网络波动或文件过大

**解决：**
```go
// 增加超时时间
client.Config.Timeout = 300 * time.Second

// 使用分片上传（大文件）
// OSS: PutObjectFromFile with partSize
// S3: CreateMultipartUpload
```

## 总结

ClarkGo 的云存储集成提供了：

- ✅ **灵活性** - 三种驱动随意切换
- ✅ **简单性** - 统一接口零学习成本
- ✅ **可靠性** - 成熟 SDK 保障稳定
- ✅ **扩展性** - 易于添加新驱动

无论是开发环境的快速迭代，还是生产环境的高并发场景，都能找到合适的存储方案。
