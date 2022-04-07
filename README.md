# globalVariable

找到那些在 init 函数之外的地方修改过的全局变量以及对应修改的位置。

下载库并使用 github.com/google/pprof 做测试：
``` shell
go get -v  github.com/zdyj3170101136/globalVariable/cmd/globalVariable
git clone https://github.com/google/pprof
cd pprof && globalVariable /Users/jie.yang05/Downloads/pprof-master/...
```

输出：
``` shell
yangjie05-mac:pprof-master jie.yang05$ globalVariable  /Users/jie.yang05/Downloads/pprof-master/...
/Users/jie.yang05/Downloads/pprof-master/internal/graph/dotgraph.go:575:2: global variable
/Users/jie.yang05/Downloads/pprof-master/internal/graph/dotgraph.go:580:2: modify
/Users/jie.yang05/Downloads/pprof-master/internal/graph/dotgraph.go:576:2: global variable
/Users/jie.yang05/Downloads/pprof-master/internal/graph/dotgraph.go:581:2: modify
```


