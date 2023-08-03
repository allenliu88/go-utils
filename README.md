# Go Utils

Utilities for golang development.

## Concurrency fail-fast

### Overview

通过`Response Channel`及`Error Channel`相互配合，实现Fail-Fast逻辑，具体参考[utils/utils.go](./utils/utils.go)。同时，需要注释掉示例程序[utils/sample.go](./utils/sample.go)中关于`sync.WaitGroup`的逻辑（`DoJob`中的前3行代码）。

### Implemetation

一般实现范式如下：

- 定义Cancel Context，返回`cancel`函数，并`defer cancel()`，后续的Goroutine中使用该Cancel Context，以此来实现当其中一个请求失败后，理解Cancel剩余的请求
- 定义RespChan、ErrChan
- For循环创建并发Goroutine，同时，记录For循环总数；Goroutine内响应Cancel Context，并在出问题时将错误写入ErrChan，否则将结果写入RespChan（可以写入结果，或者Done标识等）
- For-Select循环如上步骤的循环总数，分别选择RespChan及ErrChan，RespChan可以按需处理（是当次循环的一个结束条件），ErrChan判断获取的值是否`nil`，如果非空，则调用`cancel()`函数，然后快速失败`return`第一个Error即可
- 注意：可以有一些变种，例如，如果Goroutine中不论是成功、失败都将结果写入到ErrChan，也即成功时写入`nil`，则只需要一个ErrChan即可，不需要RespChan；

### Run

```shell
go run main.go

WARN[0000] doing the job 0 of request [1ab557f4-e8ae-40c5-a7f1-75d02e70ba35] 
INFO[0000] successfully done for the job 0 of request [1ab557f4-e8ae-40c5-a7f1-75d02e70ba35] 
WARN[0000] doing the job 1 of request [1ab557f4-e8ae-40c5-a7f1-75d02e70ba35] 
WARN[0000] doing the job 0 of request [5b93132e-5630-4a85-9cdd-944014a38418] 
WARN[0000] doing the job 0 of request [267b3ff5-022e-4036-b094-51ada36f7c09] 
Result: [], Error: error occurred for doing the job 1 of request [1ab557f4-e8ae-40c5-a7f1-75d02e70ba35]
INFO[0002] successfully done for the job 0 of request [5b93132e-5630-4a85-9cdd-944014a38418] 
ERRO[0002] cancel occurred for the job 1 of request [5b93132e-5630-4a85-9cdd-944014a38418] 
INFO[0003] successfully done for the job 0 of request [267b3ff5-022e-4036-b094-51ada36f7c09] 
ERRO[0003] cancel occurred for the job 1 of request [267b3ff5-022e-4036-b094-51ada36f7c09] 
```

结果解析：

- 其中，请求`1ab557f4-e8ae-40c5-a7f1-75d02e70ba35`进入到了`1`任务，然后发生异常，提前返回了，同时触发了`cancel()`操作，因此，结果中首次错误也是该请求的任务`1`
- 另外，注意另外两个请求`5b93132e-5630-4a85-9cdd-944014a38418`以及`267b3ff5-022e-4036-b094-51ada36f7c09`还没有来得及进入执行任务`1`（没有`WARN[000] doing the job 1 of request...`的日志）的分支，因此，走到了`<-ctx.Done()`分支，因此产生了`ERRO`错误日志
- 因此，达到了Fail-Fast的目的。同时，通过`context.WithValue`传递`sync.WaitGroup`主要是为了演示全部`Cancel`的流程，阻止`main`函数提前结束；实际上如果去掉[utils/sample.go](./utils/sample.go)中`sync.WaitGroup`部分的逻辑，是不会有此处的`ERRO`日志的，主要是因为提前返回后，`main`函数提前结束了，因此，协程也会直接结束掉

## License

Copyright &copy; 2023 Allen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
