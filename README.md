## 介绍
异步任务系统，处理异步结果，支持多级任务、顺序任务和任务状态，并记录每个环节入参和出参。
##  功能
- 支持父级任务和子级任务
- 支持任务的顺序执行和并发执行
- 任务有状态，每个环节有入参和出参
- 支持自动重试，可指定自动重试次数和间隔时间
- 支持手工跳过和手工重试，重试可以编辑参数
- 去中心化，去单点
- 支持环境参数添加
- 支持幂等性，支持某页面的任务只执行一次
- 支持令牌机制，同令牌是互斥的
- 支持任务预占，可以指定某些任务在某些节点上执行
- 支持任务同步等待

## 安装方法：
```go
go get github.com/sqc157400661/jobx
```
## 使用方法：
### 初始化：
```go
// 一个服务进程，请只实例化一个JobFlow对象
// 必选参数1：uniqueID 能唯一标识一个JobFlow对象，请不要重复
// 必选参数2：engine  *xorm.Engine对象
jobFlow, err := NewJobFlow("uniqueID", engine)

// 可选参数
options.Desc("jobflow desc info")
options.LoopInterval(5 * time.Second) // control the time interval of cyclic data fetching
options.PoolLen(10) // control the number of tasks executed at the same time

// 注册执行主体提供者
// demo 是主体提供者标识/名称
// taskProvider 是实现了TaskProvider接口的实例
// 单个添加
jobFlow.AddProvider(taskProvider,"taskProviderName")
// 批量注册 
jobFlow.Register(taskProvider1,taskProvider2,taskProvider3)

```
### （1）添加一个简单的任务：
```go
// 添加一个job，并添加一个任务节点
// job:  必选参数1：job名称， 必选参数2：job拥有者
// task: 必选参数1：名称，    必选参数2：执行主体名称
err := NewJober("jober1", "sqc").
		AddPipeline("task_1", "demo").
		Exec()

// 可选参数
options.JobDesc("job desc info")      // job&task的描述
options.JobInput(input)               // job&task的入参
options.JobEnv(env)                   // job&task的ENV
options.BizId("uuid")                 // job的业务ID，不可重复
options.Pause(1)                      // job&task是否可暂停，1代码可暂停，默认可暂停
options.Pause(0)                      // job&task是否可暂停，0代码不可暂停
options.NoPause()                     // 等同于options.Pause(0)
options.RetryNum(3)                   // task自动重试的次数，默认3次
options.RetryGapSecond(5)             // task重试的时间间隔梯度，默认是10秒
options.AddTokens([]string{"a","b"})  // 指定令牌，同类任务令牌是互斥的，只有当一个任务完成或则废弃后令牌释放后，下一个任务才可以执行
options.PreLock("127.0.0.1")          // 任务预占逻辑，指定任务执行者
options.Sync(true,2)                  // 指定该是否等待任务的执行结果,以及等待的超时时间
```

### （2）添加一个多个任务节点的顺序任务：
```go
// 添加一个job，并添加多个任务节点，任务节点执行顺序执行
// JobInput可以添加入参，属于job和task的可选参数
// JobEnv 可以添加环境标识，属于job和task的可选参数
err := NewJober("jober1", "sqc",JobEnv(map[string]interface{}{"env":"test"})).
	AddPipeline("task_1", "demo").
	AddPipeline("task_2", "demo").
	AddPipeline("task_3", "demo",JobInput(map[string]interface{}{"all_config": "yes",})).
Exec()
```

### （3）添加多个子任务：
```go
// 添加多个子任务，子任务直接是可以并发执行的，互不影响
jobTwo := NewJober("jober3", "sqc")
jobTwo.AddJob(
    NewJober("jober3_1", "sqc").
        AddPipeline("task_1", "demo").
        AddPipeline("task_2", "demo").
        AddPipeline("task_3", "demo")
    ).Exec()
jobTwo.AddJob(
    NewJober("jober3_2", "sqc").
        AddPipeline("task_1", "demo").
        AddPipeline("task_2", "demo").
        AddPipeline("task_3", "demo")
    ).Exec()
```


## 文档

## 参考
1. 项目结构：https://github.com/golang-standards/project-layout/blob/master/README_zh.md

