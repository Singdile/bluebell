# BlueBell Score 异步更新方案重构指南

## 一、方案概述

### 设计原则

| 原则 | 说明 |
|------|------|
| Redis 存全量索引 | 所有帖子 ID 都在 Redis 中（用于排序加速） |
| MySQL 是数据源 | Redis 为空时从 MySQL 查询 |
| 数据全量统一 | Redis 和 MySQL 数据一致 |
| 单次查询单数据源 | 一次分页查询只用 Redis 或 MySQL，不混合 |
| Score 异步更新 | 投票时先写 Redis，再异步批量写 MySQL |
| 投票数据只存 Redis | 投票详情不持久化到 MySQL，Redis 为空时投票数据为空 |

### Redis 的作用

```
Redis 的作用不是"缓存数据"，而是"索引加速"

┌─────────────────────────────────────────────────────────────┐
│                      查询流程                                │
├─────────────────────────────────────────────────────────────┤
│  1. Redis 排序 ──→ 获取有序的 post_id 列表                   │
│     └── 快：ZSet 跳表，O(log n)                             │
│                                                              │
│  2. MySQL 查询 ──→ 根据 post_id 获取详情                     │
│     └── 快：主键查询，O(m)                                   │
│                                                              │
│  总体：O(log n + m)，比 MySQL 排序 O(n) 快很多               │
└─────────────────────────────────────────────────────────────┘
```

### 投票数据存储策略

```
投票数据的特点：
├── 投票是短期行为（7 天后不能投票）
├── 投票数据变化频繁（实时更新）
├── 投票数据不需要长期持久化
└── Redis ZSet 天然适合存储投票关系

存储策略：
├── 投票详情只存 Redis
│   └── KeyPostVotedZsetPrefix + post_id = {user_id: vote_direction}
├── Score（投票结果）存 MySQL（持久化）
│   └── post 表的 score 字段
└── Redis 为空时，投票数据为空是合理的
    └── 因为投票是短期行为，可以重新积累
```

### 核心逻辑

```
┌─────────────────────────────────────────────────────────────┐
│                      查询流程                                │
├─────────────────────────────────────────────────────────────┤
│  用户请求第 N 页                                             │
│      │                                                       │
│      ▼                                                       │
│  1. 获取总数（固定）                                         │
│      ├── 全局：Redis.ZCard(KeyPostTimeZset)                  │
│      ├── 社区：Redis.SCard(KeyPostCommunitySetPrefix:ID)     │
│      └── fallback：MySQL.COUNT()                            │
│      │                                                       │
│      ▼                                                       │
│  2. 计算总页数（固定）                                       │
│      totalpage = total / pagesize                            │
│      │                                                       │
│      ▼                                                       │
│  3. 分页查询数据                                             │
│      ├── Redis.ZRange() ──→ 有数据 → 返回                   │
│      └── MySQL.LIMIT OFFSET ──→ fallback                    │
│      │                                                       │
│      ▼                                                       │
│  4. 返回响应                                                 │
│      Total: 固定                                             │
│      TotalPages: 固定                                        │
│      Page: 当前页                                            │
│      List: 当前页数据                                        │
│      VoteP/VoteN: Redis 有数据时返回，否则为空               │
└─────────────────────────────────────────────────────────────┘
```

---

## 二、代码检查结果

### 文件状态汇总

| 序号 | 文件 | 状态 | 问题数 |
|------|------|------|--------|
| 1 | `dao/redis/post.go` | ⚠️ 有问题 | 1 个 |
| 2 | `dao/mysql/post.go` | ✅ 正常 | 0 个 |
| 3 | `dao/mysql/score_queue.go` | ✅ 已修复 | 0 个 |
| 4 | `logic/post.go` | ❌ 多个错误 | 7 个 |
| 5 | `dao/redis/vote.go` | ✅ 正常 | 0 个 |
| 6 | `main.go` | ✅ 正常 | 0 个 |

---

## 三、问题详情与修复方案

### 问题 1：`dao/redis/post.go` - 函数缺少返回值

**位置**：第 28-38 行

**问题**：`GetCommunityPostCount` 函数缺少 `return` 语句

```go
// 当前代码（错误）
func GetCommunityPostCount(ctx *gin.Context, communityid int64) (int64,error) {
    key := KeyPostCommunitySetPrefix + strconv.FormatInt(communityid,10)
    count,err := rdb.SCard(ctx,key).Result()
    if err != nil {
        zap.L().Error("Failed to get total post count from redis for community",zap.Error(err))
        return 0,err
    }
    // ❌ 缺少 return count, nil
}
```

**修复**：

```go
// 修复后
func GetCommunityPostCount(ctx *gin.Context, communityid int64) (int64,error) {
    key := KeyPostCommunitySetPrefix + strconv.FormatInt(communityid,10)
    count,err := rdb.SCard(ctx,key).Result()
    if err != nil {
        zap.L().Error("Failed to get total post count from redis for community",zap.Error(err))
        return 0,err
    }
    return count, nil  // ✅ 添加返回值
}
```

---

### 问题 2：`logic/post.go` - 函数调用参数不匹配

**位置**：第 178 行

**问题**：`getCommunityPostCount()` 缺少 `ctx` 参数

```go
// 当前代码（错误）
total, err := getCommunityPostCount(query.CommunityID)  // ❌ 缺少 ctx

// 修复后
total, err := getCommunityPostCount(ctx, query.CommunityID)  // ✅ 添加 ctx
```

---

### 问题 3：`logic/post.go` - 返回值类型错误

**位置**：第 180 行

**问题**：`return err` 应该返回 `nil, err`

```go
// 当前代码（错误）
if err != nil {
    return err  // ❌ 函数返回 (*models.PostListResponse, error)
}

// 修复后
if err != nil {
    return nil, err  // ✅ 返回 (nil, err)
}
```

---

### 问题 4：`logic/post.go` - 函数调用参数不匹配

**位置**：第 216 行

**问题**：`mysql.GetPostListByIDs()` 返回值不匹配

```go
// 当前代码（错误）
list, total, err := mysql.GetPostListByIDs(postids, query.CommunityID)
// ❌ GetPostListByIDs 只返回 ([]*models.PostListItem, error)

// 修复后
list, err := mysql.GetPostListByIDs(postids, query.CommunityID)  // ✅ 只接收两个返回值
```

---

### 问题 5：`logic/post.go` - 变量重复定义

**位置**：第 184 行和第 228 行

**问题**：`totalpage` 变量重复定义

```go
// 当前代码（错误）
totalpage := total / query.Pagesize  // 第 184 行
// ...
totalpage := int64(len(postids)) / query.Pagesize  // 第 228 行，重复定义

// 修复后：删除第 228 行的重复定义
totalpage := total / query.Pagesize  // ✅ 只保留第 184 行
if total % query.Pagesize != 0 {
    totalpage++
}
```

---

### 问题 6：`logic/post.go` - Total 计算错误

**位置**：第 236 行

**问题**：`Total` 应该使用前面获取的 `total`，而不是 `len(postids)`

```go
// 当前代码（错误）
Total: int64(len(postids)),  // ❌ 应该用 total

// 修复后
Total: total,  // ✅ 使用前面获取的 total
```

---

### 问题 7：`logic/post.go` - 函数调用参数不匹配

**位置**：第 267 行

**问题**：`redis.GetTotalPostCount()` 缺少 `ctx` 参数

```go
// 当前代码（错误）
count, err := redis.GetTotalPostCount()  // ❌ 缺少 ctx

// 修复后
count, err := redis.GetTotalPostCount(ctx)  // ✅ 添加 ctx
```

---

## 四、完整修复代码

### 4.1 `dao/redis/post.go` 修复

只需在第 38 行添加返回值：

```go
// GetCommunityPostCount 获取社区包含的帖子的总数
func GetCommunityPostCount(ctx *gin.Context, communityid int64) (int64,error) {
    // 构建社区保存的post_id 的无序集合 key
    key := KeyPostCommunitySetPrefix + strconv.FormatInt(communityid,10)

    // 统计无序集合的数量
    count,err := rdb.SCard(ctx,key).Result()
    if err != nil {
        zap.L().Error("Failed to get total post count from redis for community",zap.Error(err))
        return 0,err
    }

    return count, nil  // ✅ 添加返回值
}
```

---

### 4.2 `logic/post.go` 完整修复

```go
package logic

import (
    "bluebell/dao/mysql"
    "bluebell/dao/redis"
    "bluebell/models"
    "bluebell/pkg/snowflake"
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// CreatePost 创建post
func CreatePost(ctx *gin.Context, post *models.Post) (err error) {
    // 创建post_id
    post_id := snowflake.GenID()
    post.ID = post_id

    // 初始化CreateTime
    post.CreateTime = time.Now()

    // 初始化 Score（等于创建时间戳）
    post.Score = post.CreateTime.Unix()

    //保存到数据库
    err = mysql.InsertPost(post)
    if err != nil {
        return err
    }
    //保存到reids
    err = redis.CreatePost(ctx, post)
    if err != nil {
        return err
    }
    //返回
    return err
}

// GetPostByID 根据ID查询post
func GetPostDetailByID(post_id int64) (*models.PostDetail, error) {
    // 查询数据库，获取对应的信息
    return mysql.GetPostDetailByID(post_id)
}

// getPostListByOrder 获取帖子列表（按排序）
func getPostListByOrder(ctx *gin.Context, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
    // 先获取帖子的总数
    total, err := getTotalPostCount(ctx)
    if err != nil {
        return nil, err
    }

    // 计算总页数
    totalpage := total / postquery.Pagesize
    if total%postquery.Pagesize != 0 {
        totalpage++
    }

    // 分页查询数据
    postids, err := redis.GetPostIDs(ctx, postquery)

    // redis 正常有数据 ——> 返回redis结果
    if err == nil && len(postids) > 0 {
        responsedata, err := buildResponseFromRedis(ctx, postids, postquery)
        if err != nil {
            return nil, err
        }
        responsedata.Total = total
        responsedata.TotalPages = totalpage
        return responsedata, nil
    }

    // redis 查询失败 ——> 查询mysql
    responsedata, err := getPostListFromMySQL(ctx, postquery)
    if err != nil {
        return nil, err
    }
    responsedata.Total = total
    responsedata.TotalPages = totalpage
    return responsedata, nil
}

// buildResponseFromRedis redis查询到postids,再到mysql中查找对应的
func buildResponseFromRedis(ctx *gin.Context, postids []string, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
    //按照post_id到mysql数据库中查询post
    list, err := mysql.GetPostListByIDs(postids, 0)
    if err != nil {
        return nil, err
    }

    //查询每篇帖子的投票数
    votePList, voteNList, err := redis.GetPostVoteData(ctx, postids)
    if err != nil {
        zap.L().Warn("Failed to get vote data", zap.Error(err))
        votePList = []string{}
        voteNList = []string{}
    }

    //构造成响应结构体返回
    responsedata := &models.PostListResponse{
        Page:     postquery.Page,
        PageSize: postquery.Pagesize,
        List:     list,
        VoteP:    votePList,
        VoteN:    voteNList,
    }

    return responsedata, nil
}

// getPostListFromMySQL 从mysql中分页查询postlist
// 注意：投票数据只存 Redis，MySQL fallback 时投票数据为空
func getPostListFromMySQL(ctx *gin.Context, postquery *models.ParamPostQuery) (*models.PostListResponse, error) {
    // MySQL 按序分页查询获取数据
    list, err := mysql.GetPostListByOrder(postquery.Page, postquery.Pagesize, postquery.Order)
    if err != nil {
        return nil, err
    }

    // 构造返回
    // 投票数据只存 Redis，MySQL fallback 时投票数据为空
    responsedata := &models.PostListResponse{
        Page:     postquery.Page,
        PageSize: postquery.Pagesize,
        List:     list,
        VoteP:    []string{},
        VoteN:    []string{},
    }

    return responsedata, nil
}

// getPostListByCommunity 获取社区帖子列表
func getPostListByCommunity(ctx *gin.Context, query *models.ParamPostQuery) (*models.PostListResponse, error) {
    // 获取社区帖子总数
    total, err := getCommunityPostCount(ctx, query.CommunityID)  // ✅ 修复：添加 ctx
    if err != nil {
        return nil, err  // ✅ 修复：返回 (nil, err)
    }

    // 计算总页数
    totalpage := total / query.Pagesize
    if total % query.Pagesize != 0 {
        totalpage++
    }

    // 该社区没有帖子
    if total == 0 {
        return &models.PostListResponse{
            Total:      0,
            Page:       query.Page,
            PageSize:   query.Pagesize,
            TotalPages: 0,
            List:       []*models.PostListItem{},
            VoteP:      []string{},
            VoteN:      []string{},
        }, nil
    }

    // 查询社区按照对应的order的post_id 列表
    postids, err := redis.GetPostIDs(ctx, query)
    if err != nil {
        return nil, err
    }

    if len(postids) == 0 {
        return &models.PostListResponse{
            Total:      0,
            Page:       query.Page,
            PageSize:   query.Pagesize,
            TotalPages: 0,
            List:       []*models.PostListItem{},
            VoteP:      []string{},
            VoteN:      []string{},
        }, nil
    }

    // 按照post_id到mysql数据库中查询数据
    list, err := mysql.GetPostListByIDs(postids, query.CommunityID)  // ✅ 修复：只接收两个返回值
    if err != nil {
        return nil, err
    }

    // 根据post_id 列表 查询每篇帖子的投票数
    votePList, voteNList, err := redis.GetPostVoteData(ctx, postids)
    if err != nil {
        zap.L().Warn("Failed to get vote data", zap.Error(err))
        votePList = []string{}
        voteNList = []string{}
    }

    // 构造响应结构体返回
    responsedata := &models.PostListResponse{
        Total:      total,  // ✅ 修复：使用 total 而不是 len(postids)
        Page:       query.Page,
        PageSize:   query.Pagesize,
        TotalPages: totalpage,
        List:       list,
        VoteP:      votePList,
        VoteN:      voteNList,
    }
    return responsedata, nil
}

// GetPostList 统一查询帖子列表的接口
func GetPostList(ctx *gin.Context, query *models.ParamPostQuery) (*models.PostListResponse, error) {
    if query.CommunityID == 0 {
        return getPostListByOrder(ctx, query)
    } else {
        return getPostListByCommunity(ctx, query)
    }
}

// VotePost 用户为帖子投票
func VotePost(ctx *gin.Context, user_id int64, vote *models.ParamVote) error {
    return redis.VotePost(ctx, fmt.Sprintf("%v", user_id), fmt.Sprintf("%v", vote.PostID), float64(vote.Direction))
}

// getTotalPostCount 获取帖子总数(先redis， fallback mysql)
func getTotalPostCount(ctx *gin.Context) (int64, error) {
    //优先从redis中获取
    count, err := redis.GetTotalPostCount(ctx)  // ✅ 修复：添加 ctx 参数
    if err == nil && count > 0 {
        return count, nil
    }

    // reids 获取失败，fallback 到 mysql
    return mysql.GetTotalPostCount()
}

// getCommunityPostCount 获取社区帖子总数 (先redis， fallback mysql)
func getCommunityPostCount(ctx *gin.Context, communityid int64) (int64, error) {
    // 优先从redis中获取
    count, err := redis.GetCommunityPostCount(ctx, communityid)
    if err == nil && count > 0 {
        return count, nil
    }

    // redis 获取失败， fallbcak 到 mysql
    return mysql.GetCommunityPostCount(communityid)
}
```

---

## 五、问题汇总表

| 序号 | 文件 | 行号 | 问题类型 | 问题描述 |
|------|------|------|---------|---------|
| 1 | `dao/redis/post.go` | 38 | 缺少返回值 | `GetCommunityPostCount` 缺少 `return count, nil` |
| 2 | `logic/post.go` | 178 | 参数缺失 | `getCommunityPostCount` 缺少 `ctx` 参数 |
| 3 | `logic/post.go` | 180 | 返回值错误 | `return err` 应为 `return nil, err` |
| 4 | `logic/post.go` | 216 | 返回值不匹配 | `GetPostListByIDs` 只返回两个值 |
| 5 | `logic/post.go` | 228 | 变量重复 | `totalpage` 重复定义 |
| 6 | `logic/post.go` | 236 | 逻辑错误 | `Total` 应使用 `total` 而非 `len(postids)` |
| 7 | `logic/post.go` | 267 | 参数缺失 | `redis.GetTotalPostCount()` 缺少 `ctx` |

---

## 六、修改顺序建议

建议按以下顺序修改：

| 步骤 | 文件 | 修改内容 |
|------|------|---------|
| 1 | `dao/redis/post.go` | 添加 `return count, nil` |
| 2 | `logic/post.go` | 修复所有函数调用和返回值 |
| 3 | 编译测试 | `go build ./...` |

---

## 七、投票数据存储说明

### 为什么投票数据只存 Redis？

```
投票数据的特点：
├── 时间限制：7 天后不能投票
│   └── 投票数据是"短期有效"的
├── 变化频繁：用户随时可以改票
│   └── Redis 的原子操作更适合
├── 数据量大：每篇帖子可能有几百个投票
│   └── MySQL 存储压力大
└── Score 已经持久化：
    └── 投票的"结果"（Score）已经异步写入 MySQL
    └── 投票的"详情"（谁投了什么）不需要长期保存
```

### fallback 时投票数据为空是合理的

```
当 Redis 为空时：
├── 帖子列表从 MySQL 查询（有 score）
├── 投票数据为空（vote_p=[], vote_n=[]）
└── 这是合理的，因为：
    ├── Redis 为空意味着系统重启或索引重建
    ├── 投票数据可以重新积累
    └── 用户可以重新投票（7 天内）
```

---

## 八、验证步骤

### 1. 编译检查

```bash
cd ~/Work/bluebell
go build ./...
```

### 2. 启动服务

```bash
go run main.go
```

### 3. 测试接口

```bash
# 测试全局帖子列表
curl "http://localhost:8080/api/v1/posts?page=1&size=10&order=time"

# 测试社区帖子列表
curl "http://localhost:8080/api/v1/posts?page=1&size=10&order=time&community_id=1"

# 测试投票
curl -X POST "http://localhost:8080/api/v1/vote" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"post_id": 123, "direction": 1}'
```

---

## 九、参考文件路径

| 文件 | 路径 |
|------|------|
| post.go (redis) | `~/Work/bluebell/dao/redis/post.go` |
| post.go (mysql) | `~/Work/bluebell/dao/mysql/post.go` |
| score_queue.go | `~/Work/bluebell/dao/mysql/score_queue.go` |
| post.go (logic) | `~/Work/bluebell/logic/post.go` |
| vote.go (redis) | `~/Work/bluebell/dao/redis/vote.go` |
| main.go | `~/Work/bluebell/main.go` |
