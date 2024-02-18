


============================================== 参考 https://zhuanlan.zhihu.com/p/664770047

纯Go语言工具go-ora来实现可配置化Oracle数据库巡检
香生享IT

​
目录
收起
1 Go相关包介绍
database包
go-ora工具
yaml.v3
go-pretty
2 程序实现
功能概述
安装准备
主要代码
3 程序使用
创建项目
编译运行
改换下代码
总结
引言
在数据库管理领域，DBA（数据库管理员）的日常工作涉及到对数据库进行巡检，以确保其稳定运行和高效性能。为了简化和自动化这个过程，本文介绍了如何使用Go语言和一些开源包来实现可配置化的数据库巡检。通过使用go-ora开源包，DBA可以轻松地进行数据库巡检，并根据实际需求进行配置和定制。
数据库的稳定性和高效性能对于任何企业来说都至关重要。作为DBA，巡检数据库是日常工作的重要组成部分，以确保数据库的正常运行和最佳性能。然而，传统的巡检方法通常需要手动编写大量的SQL查询和脚本来检查各种指标和配置，这既耗时又繁琐。

为了简化这个过程，我们可以利用Go语言的强大功能和一些优秀的开源包来实现可配置化的数据库巡检。首先，我们使用go-ora开源包，它提供了与Oracle数据库的交互和查询功能，使我们能够轻松地连接、查询和分析数据库。其次，我们使用yaml.v3包来解析和读取配置文件，通过配置文件可以定义要检查的指标和规则。此外，我们还使用go-pretty包来美化和格式化输出结果，使巡检报告更加易读和优雅。

在实现可配置化的数据库巡检过程中，还可以使用database/sql包来与不同的数据库进行交互，这使得我们可以轻松地切换和扩展支持的数据库类型。通过使用这些开源包，DBA可以根据实际需求配置和定制数据库巡检流程。无论是检查重要性能指标、验证配置是否符合最佳实践，还是监控数据库安全性，都可以通过简单的配置来实现。这种可配置化的巡检方法不仅提高了效率，还减少了出错的可能性，让DBA能够更专注于核心任务，提升数据库的稳定性和性能。

1 Go相关包介绍
database包
database/sql包

包sql提供了sql(或类sql)数据库的通用接口。sql包必须与数据库驱动程序一起使用。参见https://golang.org/s/sqldrivers 获取驱动程序列表。如下列出一些代码中用到常用类型或接口：

func Open：Open打开由其数据库驱动程序名称和特定于驱动程序的数据源名称指定的数据库，通常数据源名称至少由数据库名称和连接信息组成。大多数用户将通过特定于驱动程序的连接辅助函数打开数据库，该函数返回一个*DB。Open可以只验证它的参数，而不创建到数据库的连接。要验证数据源名称是否有效，请调用Ping。返回的DB对于多个例程的并发使用是安全的，并维护自己的空闲连接池。因此，应该只调用Open函数一次。很少需要关闭DB。
func Open(driverName, dataSourceName string) (*DB, error)
Type DB：DB是一个数据库句柄，表示一个由零个或多个底层连接组成的池。对于多个例程的并发使用是安全的。sql包自动创建和释放连接;它还维护一个空闲连接的空闲池。如果数据库具有每个连接状态的概念，则可以在事务(Tx)或连接(Conn)中可靠地观察到这种状态。一旦DB.Begin调用时，返回的Tx绑定到单个连接。一旦在事务上调用Commit或Rollback，该事务的连接将返回到DB的空闲连接池。池的大小可以通过SetMaxIdleConns来控制。
func (*DB) Close：Close关闭数据库并阻止新的查询开始。等待服务器上已开始处理的所有查询完成然后关闭。
func (db *DB) Close() error
func (*DB) Ping：Ping验证到数据库的连接是否仍然存在，如果需要，则建立连接。
func (db *DB) Ping() error
func (*DB) Query：Query执行返回行的查询，通常是SELECT。参数用于查询中的任何占位符参数。
func (db *DB) Query(query string, args ...any) (*Rows, error)
type Rows：Rows是查询的结果，其游标从结果集的第一行之前开始。使用“Next”从一行移动到另一行。
func (*Rows) Columns：Columns返回列名。如果rows已经关闭，Columns返回一个错误。
func (rs *Rows) Columns() ([]string, error)
func (*Rows) Next：Next准备下一个结果行，以便使用Scan方法读取。成功时返回true，如果没有下一个结果行或在准备结果时发生错误则返回false。应该参考Err来区分这两种情况。 每次调用Scan，甚至是第一次调用，都必须先调用Next。
func (rs *Rows) Next() bool
func (*Rows) Scan ：Scan将当前行的列复制到dest所指向的值中，dest中的值数必须与Rows中的列数相同。
func (rs *Rows) Scan(dest ...any) error
func (*Rows) Close：Close关闭Rows，防止进一步枚举。如果调用Next并返回false，并且没有进一步的结果集，则自动关闭行，并且检查Err的结果。Close是幂等的，不影响Err的结果。
func (rs *Rows) Close() error
database/sql/driver包

包Driver定义由数据库驱动程序实现的接口，如包sql所使用的接口。大多数代码应该使用包sql。

随着时间的推移，驱动程序接口也在不断发展。驱动程序应该实现Connector和DriverContext接口。Connector.Connect和Driver.Open方法不应该返回ErrBadConn。ErrBadConn应该只在连接已经处于无效(例如关闭)状态时才从Validator, SessionResetter或查询方法返回。

所有Conn实现都应该实现以下接口:Pinger, SessionResetter 和Validator。 如果支持命名参数或上下文，驱动程序的Conn应该实现:ExecerContext, QueryerContext, ConnPrepareContext和ConnBeginTx。

关于database/sql和database/sql/driver包详细的介绍，可以参考：https://pkg.go.dev/database/sql , https://pkg.go.dev/database/sql/driver

go-ora工具
go-ora是原生Go实现的一款不需要依赖CGo的纯GO语言Oracle客户端开源工具。Github社区也比较活跃，一直有更新。主要功能如下

函数清单

类型清单





本文代码中用到的函数BuildUrl通过服务器，端口，服务，用户，密码，urlOptions创建databaseURL，这个函数帮助建立一个将形成的databaseURL字符串用于连接数据库。
func BuildUrl(server string, port int, service, user, password string, options map[string]string) string
其中操作数据库的CRUD例子代码：https://github.com/sijms/go-ora/blob/master/examples/crud/main.go

yaml.v3
yaml包使Go程序能够轻松地对yaml值进行编码和解码。它是作为juju项目的一部分在Canonical内部开发的，基于著名的libyaml C库的纯Go端口，可以快速可靠地解析和生成YAML数据。

可以通过引用程序包gopkg.in/yaml.v3，或者通过运行go get gopkg.in/yaml.v3来安装。


Marshal将提供的值序列化到YAML文档中。生成文档的结构将反映值本身的结构。映射和指针(struct、string、int等)被接受为in值。
Unmarshal解码在in字节片中找到的第一个文档，并将解码后的值赋给out值。
go-pretty
go-pretty是一个用于构建漂亮和可定制的终端输出的开源工具包。它提供了丰富的功能和样式选项，可以帮助我们更好地展示和呈现数据。

go-pretty 可以帮助我们解决在终端输出中遇到的一些问题。它可以处理表格、列表、进度条等各种数据的输出，并且支持自定义颜色、对齐方式、边框样式等。通过使用 go-pretty，我们可以轻松地将数据以优雅的方式展示给用户，提高用户体验。

go-pretty 的主要特性包括：

表格输出：go-pretty 提供了灵活且易于使用的表格输出功能，可以按照需要添加、删除和编辑表头和行数据。同时，我们还可以对表格进行排序、筛选和分页等操作。
列表输出：除了表格，go-pretty 还支持列表输出。我们可以使用不同的样式和符号来呈现列表数据，使其更具可读性。
进度条：go-pretty 可以帮助我们在终端中显示进度条，以便于用户了解任务的进展情况。
颜色和样式：go-pretty 支持自定义颜色、背景色和样式，可以根据需要对输出进行修饰，使其更加美观和易读。
自定义边框：我们可以使用 go-pretty 在表格和列表周围添加自定义边框，以增加输出的可视化效果。
总之，go-pretty 是一个非常实用的工具包，可以帮助我们在终端中以漂亮和定制化的方式展示数据，提升用户体验和可读性。

这个包的当前主要版本是v6，运行go get github.com/jedib0t/go-pretty/v6将其作为一个依赖项添加到你的项目中，并在你的代码中使用以下一个或多个导入包:

github.com/jedib0t/go-pretty/v6/list
github.com/jedib0t/go-pretty/v6/progress
github.com/jedib0t/go-pretty/v6/table
github.com/jedib0t/go-pretty/v6/text
2 程序实现
功能概述
提供的Go程序主要目的是将通过解析配置文件获取Oracle数据库连接信息和要执行巡检的SQL语句，并连接到Oracle数据库执行对应的SQL语句，并格式化输出结果。

以下是程序功能的概述：

配置文件(yaml格式)解析
数据库连接管理
SQL语句执行和结果处理
结果格式化和输出
安装准备
创建用户并安装Go语言
-- https://golang.google.cn/dl/ 下载golang - 文件go1.21.3.linux-amd64.tar.gz
[root@VM-0-10-centos ~]# useradd goapp
[root@VM-0-10-centos ~]# su - goapp
[goapp@VM-0-10-centos ~]$ mkdir app
[goapp@VM-0-10-centos ~]$ tar -C app -zxf ../go1.21.3.linux-amd64.tar.gz
[goapp@VM-0-10-centos ~]$ export PATH=$PATH:$HOME/app/go/bin
[goapp@VM-0-10-centos ~]$ go version
go version go1.21.3 linux/amd64
[goapp@VM-0-10-centos ~]$ go env -w GOPROXY=https://goproxy.cn,direct
获取相关依赖包(在线方式）
[goapp@VM-0-10-centos ~]$ go install github.com/sijms/go-ora/v2@latest
go: downloading github.com/sijms/go-ora/v2 v2.7.19
go: downloading github.com/sijms/go-ora v1.3.2
package github.com/sijms/go-ora/v2 is not a main package
[goapp@VM-0-10-centos ~]$ go install gopkg.in/yaml.v3@latest
go: downloading gopkg.in/yaml.v3 v3.0.1
package gopkg.in/yaml.v3 is not a main package
[goapp@VM-0-10-centos ~]$ go install github.com/jedib0t/go-pretty/v6@latest
go: downloading github.com/jedib0t/go-pretty/v6 v6.4.9
go: downloading github.com/jedib0t/go-pretty v4.3.0+incompatible
package github.com/jedib0t/go-pretty/v6 is not a main package
[goapp@VM-0-10-centos ~]$ go install github.com/jedib0t/go-pretty/v6/text@latest
go: downloading github.com/mattn/go-runewidth v0.0.13
go: downloading golang.org/x/sys v0.1.0
go: downloading github.com/rivo/uniseg v0.2.0
package github.com/jedib0t/go-pretty/v6/text is not a main package
[goapp@VM-0-10-centos ~]$ ls -l go/pkg/mod/cache/download/
...
[goapp@VM-0-10-centos ~]$ tar -cf go_cache_download.tar go/pkg/mod/cache/download/
局域网离线环境中，拷贝上述下载好的包到离线环境通过指定变量GOPROXY来获取
# su - goapp
$ mkdir app
$ tar -C app -zxf ../go1.21.3.linux-amd64.tar.gz
$ tar -xf go_download.tar
$ go env -w GOSUMDB=off
$ go env -w GOPROXY=file://home/goapp/go/pkg/mod/cache/download
主要代码
参数文件解释：通过go语言自带的flag实现
//数据库连接信息和SQL语句在同一YAML文件中
    configFilePath := flag.String("config", "config.yaml", "配置文件路径")
    flag.Parse()
yaml文件及解释：通过yaml.v3提供的方法来实现
//yaml文件格式如下：config.yaml
database:
  server: "localhost"
  service_name: "orcl"
  username: "system"
  password: "***"

check_items:
  - name: 数据库信息
    sql: select b.name,b.DB_UNIQUE_NAME,b.dbid,to_char(b.created,'yyyy-mm-dd hh24:mi:ss') db_created,b.database_role,b.open_mode,c.*
         from v$database b,(select v1||'_'||v2||'.'||v3 nls_lang from
          (select value v1 from nls_database_parameters where parameter='NLS_LANGUAGE') m,
          (select value v2 from nls_database_parameters where parameter='NLS_TERRITORY') n,
          (select value v3 from nls_database_parameters where parameter='NLS_CHARACTERSET') p) c

  - name: 数据库常用参数
    sql: select inst_Id,name,value from gv$parameter
         where name in ('processes','memory_max_target','memory_target','pga_aggregate_limit',
         'pga_aggregate_target','sga_max_size','sga_target','db_cache_size','shared_pool_size','java_pool_size','large_pool_size') order by 1,2

  - name: 检查表空间使用率
    sql: ...

GO语言中实现配置文件解释
//Config结构体定义
type DatabaseConfig struct {
    Server      string `yaml:"server"`
    ServiceName string `yaml:"service_name"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

type CheckItem struct {
    Name string `yaml:"name"`
    SQL  string `yaml:"sql"`
}

type Config struct {
    Database   DatabaseConfig `yaml:"database"`
    CheckItems []CheckItem    `yaml:"check_items"`
}
//读取配置yaml文件并解析给Config结构体
func readConfig(filePath string) (*Config, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }

    config := &Config{}
    err = yaml.Unmarshal(bytes, config)
    if err != nil {
        return nil, err
    }

    return config, nil
}
创建数据库连接
// 创建数据库连接
    connStr := go_ora.BuildUrl(config.Database.Server, 1521,config.Database.ServiceName, config.Database.Username, config.Database.Password, nil)
    db, err := sql.Open("oracle", connStr)
    if err != nil {
        log.Fatalf("无法连接到Oracle数据库：%v", err)
    }
    defer db.Close()

    // 使用提供的用户名和密码来验证数据库连接
    err = db.Ping()
    if err != nil {
        log.Fatalf("数据库连接验证失败：%v", err)
    }
执行对应的SQL，由于不同SQL返回的列数量不一致，此处定义个Result结构来保存所有的结果。具体实现如下：
//结构Result用来存放检查项、结果集列名字和具体的数据
type Result struct {
    ItemName string
    Columns  []string
    Data     [][]string
}

    // 并行执行巡检
    var wg sync.WaitGroup
    resultCh := make(chan Result, len(config.CheckItems))

    for _, item := range config.CheckItems {
        wg.Add(1)
        go execSQL(&wg, item, db, resultCh)
    }

    // 等待所有巡检任务完成
    go func() {
        wg.Wait()
        close(resultCh)
    }()

//通过execSQL来执行具体的SQL
func execSQL(wg *sync.WaitGroup, item CheckItem, db *sql.DB, ch chan Result) {
    defer wg.Done()

    rows, err := db.Query(item.SQL)
    if err != nil {
        log.Printf("%s 查询失败：%v\n", item.Name, err)
        return
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        log.Printf("%s 获取列信息失败：%v\n", item.Name, err)
        return
    }
    // 创建一个与列数量相等的切片用于存储扫描结果
    scanArgs := make([]interface{}, len(columns))
    for i := range columns {
        var result string
        scanArgs[i] = &result
     }

     // 将列名添加到结果中
    result := Result{
        ItemName: item.Name,
        Columns:  columns,
        Data:  [][]string{},
    }

    for rows.Next() {
        err := rows.Scan(scanArgs...)
        if err != nil {
            log.Printf("%s 结果解析失败：%v\n", item.Name, err)
            continue
        }

        data := make([]string, len(columns))
        for i, col := range scanArgs {
            if col == nil {
                data[i] = "NULL"
            } else {
                data[i] = *(col.(*string))
            }
        }
        result.Data = append(result.Data, data)
      }
      ch <- result
}
格式化输出巡检结果
func tableOutput(res chan Result) {
    for result := range res {
        tw := table.NewWriter()
        tw.SetTitle("%s\n", result.ItemName)
        header := convertToRow(result.Columns)
        tw.AppendHeader(header)
        for _, data := range result.Data {
            tw.AppendRow(convertToRow(data))
        }
        fmt.Println(tw.Render())
        fmt.Println()
    }
}
//通过convertToRow函数将[]string转化为go-pretty中的row类型
func convertToRow(str []string) table.Row {
    row := make(table.Row, len(str))
    for i, s := range str {
        row[i] = s
    }
    return row
}
3 程序使用
创建项目
准备目录及项目mygo.space/go2ora
goapp@VM65195:~/go> mkdir -p mygo.space/go2ora
goapp@VM65195:~/go> cd mygo.space/go2ora/
goapp@VM65195:~/go/mygo.space/go2ora> go mod init mygo.space/go2ora
go: creating new go.mod: module mygo.space/go2ora
goapp@VM65195:~/go/src/mygo.space/go2ora> go env -w GOSUMDB=off
goapp@VM65195:~/go/src/mygo.space/go2ora> go env -w GOPROXY=file://home/goapp/go/pkg/mod/cache/download

--添加依赖包到项目中
goapp@VM65195:~/go/src/mygo.space/go2ora> go get github.com/sijms/go-ora/v2
go: added github.com/sijms/go-ora/v2 v2.7.19
goapp@VM65195:~/go/src/mygo.space/go2ora> cat go.mod
module mygo.space/go2ora

go 1.21.3

require github.com/sijms/go-ora/v2 v2.7.19 // indirect

goapp@VM65195:~/go/src/mygo.space/go2ora> go get gopkg.in/yaml.v3
go: added gopkg.in/yaml.v3 v3.0.1
goapp@VM65195:~/go/src/mygo.space/go2ora> go get github.com/jedib0t/go-pretty/v6
go: added github.com/jedib0t/go-pretty/v6 v6.4.9
goapp@VM65195:~/go/src/mygo.space/go2ora> go get github.com/jedib0t/go-pretty/v6/text
go: added github.com/jedib0t/go-pretty/v6 v6.4.9

[goapp@VM-0-10-centos go2ora]$ cat go.mod
module mygo.space/go2ora

go 1.21.3

require (
        github.com/jedib0t/go-pretty/v6 v6.4.9 // indirect
        github.com/mattn/go-runewidth v0.0.13 // indirect
        github.com/rivo/uniseg v0.2.0 // indirect
        github.com/sijms/go-ora/v2 v2.7.19 // indirect
        golang.org/x/sys v0.1.0 // indirect
        gopkg.in/yaml.v3 v3.0.1 // indirect
)
编译运行
准备配置文件
goapp@VM65195:~/go/src/mygo.space/go2ora> vi config.yaml
database:
  server: "localhost"
  service_name: "orcl"
  username: "system"
  password: "***"

check_items:
  - name: 数据库信息
    sql: select b.name,b.DB_UNIQUE_NAME,b.dbid,to_char(b.created,'yyyy-mm-dd hh24:mi:ss') db_created,b.database_role,b.open_mode,c.*
         from v$database b,(select v1||'_'||v2||'.'||v3 nls_lang from
          (select value v1 from nls_database_parameters where parameter='NLS_LANGUAGE') m,
          (select value v2 from nls_database_parameters where parameter='NLS_TERRITORY') n,
          (select value v3 from nls_database_parameters where parameter='NLS_CHARACTERSET') p) c

  - name: 数据库常用参数
    sql: select inst_Id,name,value from gv$parameter
         where name in ('processes','memory_max_target','memory_target','pga_aggregate_limit',
         'pga_aggregate_target','sga_max_size','sga_target','db_cache_size','shared_pool_size','java_pool_size','large_pool_size') order by 1,2

  - name: 检查表空间使用率
    sql:  SELECT
            d.status
                , d.tablespace_name
                , d.contents
                , d.extent_management
                , d.segment_space_management
                , NVL(b.allocatesize - NVL(f.freesize, 0), 0)   used_MB
                , b.allocatesize current_size_MB
                , to_char(NVL((b.allocatesize - NVL(f.freesize, 0)) / b.allocatesize * 100, 0),'990.99')||'%' pct_used
                , a.maxsize canextend_size_MB
                , to_char(NVL((b.allocatesize - NVL(f.freesize, 0)) / a.maxsize * 100, 0),'990.99')||'%' tot_pct_used
        FROM dba_tablespaces d
                , (     SELECT tablespace_name,sum(maxsize) maxsize
                        FROM (  SELECT tablespace_name, decode(autoextensible,'YES',round(sum(maxbytes)/1024/1024),round(sum(bytes)/1024/1024)) maxsize
                                        FROM dba_data_files
                                        GROUP BY tablespace_name,autoextensible
                                ) GROUP BY tablespace_name
                  ) a
                , ( SELECT tablespace_name, sum(bytes)/1024/1024 allocatesize
              from dba_data_files
              group by tablespace_name
                  ) b
                , (     SELECT tablespace_name, sum(bytes)/1024/1024 freesize
                        FROM dba_free_space
                        GROUP BY tablespace_name
                  ) f
        WHERE d.tablespace_name = a.tablespace_name(+)
        AND d.tablespace_name = b.tablespace_name(+)
        AND d.tablespace_name = f.tablespace_name(+)
        AND d.contents='PERMANENT'
        UNION ALL
        SELECT
            d.status
                , d.tablespace_name
                , d.contents
                , d.extent_management
                , d.segment_space_management
                , NVL(b.allocatesize - NVL(f.usedsize, 0), 0)   used_MB
                , b.allocatesize current_size_MB
                , to_char(NVL(NVL(f.usedsize, 0) / b.allocatesize * 100, 0),'990.99')||'%' pct_used
                , a.maxsize canextend_size_MB
                , to_char(NVL(f.usedsize,0) / a.maxsize * 100,'990.99')||'%' tot_pct_used
        FROM
            sys.dba_tablespaces d
                , (     SELECT tablespace_name,sum(maxsize) maxsize
                        FROM (  SELECT tablespace_name, decode(autoextensible,'YES',round(sum(maxbytes)/1024/1024),round(sum(bytes)/1024/1024)) maxsize
                                        FROM dba_temp_files
                                        GROUP BY tablespace_name,autoextensible
                                ) GROUP BY tablespace_name
                  ) a
          , ( select tablespace_name, sum(bytes)/1024/1024  allocatesize
              from dba_temp_files
              group by tablespace_name
            ) b
          , ( select tablespace_name, sum(bytes_cached)/1024/1024 usedsize
              from v$temp_extent_pool
              group by tablespace_name
            ) f
        WHERE d.tablespace_name = a.tablespace_name(+)
          AND d.tablespace_name = b.tablespace_name(+)
          AND d.tablespace_name = f.tablespace_name(+)
          AND d.extent_management like 'LOCAL'
          AND d.contents like 'TEMPORARY'
        ORDER By pct_used

  - name: ASM磁盘空间使用率
    sql:  SELECT
            group_number                             group_number
          , name                                     group_name
          , sector_size                              sector_size
          , block_size                               block_size
          , allocation_unit_size                     allocation_unit_size
          , state                                    state
          , type                                     type
          , database_compatibility
          , total_mb                                 total_mb
          , (total_mb - free_mb)                     used_mb
          , to_char((1- (free_mb / total_mb))*100, '990.99')||'%'   pct_used
          , free_mb/(decode(type,'HIGH',3,'NORMAL',2,1))                                avail_mb
        FROM
            v$asm_diskgroup
        ORDER BY
            name

  - name: 最近2天alert log的ORA-报错
    sql:  select
                to_char(ORIGINATING_TIMESTAMP,'yyyy-mm-dd hh24:mi:ss') originating_timestamp
                , MESSAGE_TEXT
        from V$DIAG_ALERT_EXT
        WHERE
                (MESSAGE_TEXT like '%ORA-%' or upper(MESSAGE_TEXT) like '%ERROR%')
                and ORIGINATING_TIMESTAMP > sysdate - 2
        ORDER BY originating_timestamp
运行程序
goapp@VM65195:~/go/src/mygo.space/go2ora> go run odcheck.go
+----------------------------------------------------------------------------------------------------------------------+
| 数据库信息                                                                                                           |
|                                                                                                                      |
+---------+----------------+------------+---------------------+---------------+------------+---------------------------+
| NAME    | DB_UNIQUE_NAME | DBID       | DB_CREATED          | DATABASE_ROLE | OPEN_MODE  | NLS_LANG                  |
+---------+----------------+------------+---------------------+---------------+------------+---------------------------+
| YC19CDB | yc19cdb        | 3806950467 | 2020-04-22 18:58:11 | PRIMARY       | READ WRITE | AMERICAN_AMERICA.AL32UTF8 |
+---------+----------------+------------+---------------------+---------------+------------+---------------------------+

+----------------------------------------------+
| 数据库常用参数                               |
|                                              |
+---------+----------------------+-------------+
| INST_ID | NAME                 | VALUE       |
+---------+----------------------+-------------+
| 1       | db_cache_size        | 0           |
| 1       | java_pool_size       | 0           |
| 1       | large_pool_size      | 0           |
| 1       | memory_max_target    | 0           |
| 1       | memory_target        | 0           |
| 1       | pga_aggregate_limit  | 51539607552 |
| 1       | pga_aggregate_target | 25769803776 |
| 1       | processes            | 3000        |
| 1       | sga_max_size         | 68719476736 |
| 1       | sga_target           | 0           |
| 1       | shared_pool_size     | 0           |
| 2       | db_cache_size        | 0           |
| 2       | java_pool_size       | 0           |
| 2       | large_pool_size      | 0           |
| 2       | memory_max_target    | 0           |
| 2       | memory_target        | 0           |
| 2       | pga_aggregate_limit  | 51539607552 |
| 2       | pga_aggregate_target | 25769803776 |
| 2       | processes            | 3000        |
| 2       | sga_max_size         | 68719476736 |
| 2       | sga_target           | 0           |
| 2       | shared_pool_size     | 0           |
+---------+----------------------+-------------+

+--------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| 检查表空间使用率                                                                                                                                                   |
|                                                                                                                                                                    |
+--------+-------------------+-----------+-------------------+--------------------------+------------+-----------------+----------+-------------------+--------------+
| STATUS | TABLESPACE_NAME   | CONTENTS  | EXTENT_MANAGEMENT | SEGMENT_SPACE_MANAGEMENT | USED_MB    | CURRENT_SIZE_MB | PCT_USED | CANEXTEND_SIZE_MB | TOT_PCT_USED |
+--------+-------------------+-----------+-------------------+--------------------------+------------+-----------------+----------+-------------------+--------------+
| ONLINE | TEMP              | TEMPORARY | LOCAL             | MANUAL                   | 325        | 395             |   17.72% | 32768             |    0.21%     |
| ONLINE | SYSAUX            | PERMANENT | LOCAL             | AUTO                     | 684.4375   | 820             |   83.47% | 32768             |    2.09%     |
| ONLINE | SYSTEM            | PERMANENT | LOCAL             | MANUAL                   | 480.125    | 490             |   97.98% | 32768             |    1.47%     |
...
+--------+-------------------+-----------+-------------------+--------------------------+------------+-----------------+----------+-------------------+--------------+

+------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| ASM磁盘空间使用率                                                                                                                                                      |
|                                                                                                                                                                        |
+--------------+------------+-------------+------------+----------------------+-----------+--------+------------------------+-----------+----------+----------+----------+
| GROUP_NUMBER | GROUP_NAME | SECTOR_SIZE | BLOCK_SIZE | ALLOCATION_UNIT_SIZE | STATE     | TYPE   | DATABASE_COMPATIBILITY | TOTAL_MB  | USED_MB  | PCT_USED | AVAIL_MB |
+--------------+------------+-------------+------------+----------------------+-----------+--------+------------------------+-----------+----------+----------+----------+
| 1            | DATAC1     | 512         | 4096       | 4194304              | CONNECTED | HIGH   | 11.2.0.4.0             | 230326272 | 69064908 |   29.99% | 53753788 |
| 2            | DBFS_DG    | 512         | 4096       | 4194304              | MOUNTED   | HIGH   | 11.2.0.4.0             | 1030560   | 1764     |    0.17% | 342932   |
| 3            | RECOC1     | 512         | 4096       | 4194304              | CONNECTED | NORMAL | 11.2.0.4.0             | 38423808  | 5689496  |   14.81% | 16367156 |
+--------------+------------+-------------+------------+----------------------+-----------+--------+------------------------+-----------+----------+----------+----------+

+--------------------------------------+
| 最近2天alert log的ORA-报错           |
|                                      |
+-----------------------+--------------+
| ORIGINATING_TIMESTAMP | MESSAGE_TEXT |
+-----------------------+--------------+
+-----------------------+--------------+
更换配置文件再次执行
/home/goapp/go/src/mygo.space/go2ora
goapp@VM65195:~/go/src/mygo.space/go2ora> cat 208.yaml |more
database:
  server: "localhost"
  service_name: "orcl"
  username: "system"
  password: "***"

check_items:
  - name: 数据库信息
    sql: select b.name,b.DB_UNIQUE_NAME,b.dbid,to_char(b.created,'yyyy-mm-dd hh24:mi:ss') db_created,b.database_role,b.open_mode,c.*
         from v$database b,(select v1||'_'||v2||'.'||v3 nls_lang from
          (select value v1 from nls_database_parameters where parameter='NLS_LANGUAGE') m,
          (select value v2 from nls_database_parameters where parameter='NLS_TERRITORY') n,
          (select value v3 from nls_database_parameters where parameter='NLS_CHARACTERSET') p) c

goapp@VM65195:~/go/src/mygo.space/go2ora> ./odcheck -h
Usage of ./odcheck:
  -config string
        配置文件路径 (default "config.yaml")
goapp@VM65195:~/go/src/mygo.space/go2ora> ./odcheck -config 208.yaml
+--------------------------------------------------------------------------------------------------------------------+
| 数据库信息                                                                                                         |
|                                                                                                                    |
+-------+----------------+------------+---------------------+---------------+------------+---------------------------+
| NAME  | DB_UNIQUE_NAME | DBID       | DB_CREATED          | DATABASE_ROLE | OPEN_MODE  | NLS_LANG                  |
+-------+----------------+------------+---------------------+---------------+------------+---------------------------+
| JSFMS | jsfms          | 1344789788 | 2023-10-10 12:15:24 | PRIMARY       | READ WRITE | AMERICAN_AMERICA.ZHS16GBK |
+-------+----------------+------------+---------------------+---------------+------------+---------------------------+
...
改换下代码
可能某些DBA场景需要连接不同的数据库执行不同的SQL脚本，将上述代码可以稍微进行修改下。程序接受数据库连接字符串作为参数还有指定SQL脚本文件。假设程序代码为odacheck1.go则执行过程如下：

goapp@VM65195:~/go/src/mygo.space/go2ora> go run odcheck1.go -dsn system/***@localhost/orcl -sql tbs.yaml
...
上述执行的tbs.yml的定义了关于表空间想的巡检项及其对应的SQL代码。

goapp@VM65195:~/go/src/mygo.space/go2ora> cat tbs.yaml
check_items:
  - name: 数据库信息
    sql: ...

  - name: 检查表空间使用率
    sql:  ...
为了实现上述代码修改，可在原有代码进行如下的修改：

Config 结构体移除掉数据库信息
type Config struct {
    //Database   DatabaseConfig `yaml:"database"`
    CheckItems []CheckItem    `yaml:"check_items"`
}
参数化解释增加-dsn用于接受数据库连接字符串参数
dsnStr := flag.String("dsn", "system/oracle@localhost:1521/orcl", "数据库连接配置")
    sqlFilePath := flag.String("sql", "sql.yaml", "配置文件路径")
    flag.Parse()
新增函数parseDSN来解释DNS并将结果赋值 DatabaseConfig
func parseDSN(dsn *string) (DatabaseConfig, int, error) {
    var dbConfig DatabaseConfig
    port := 1521
    parts := strings.Split(*dsn, "@")
    if len(parts) != 2 {
        return DatabaseConfig{}, port, fmt.Errorf("连接字符串格式错误")
    }
    // 解析用户名和密码部分   
    userPassword := strings.Split(parts[0], "/")
    ipPortService := strings.Split(parts[1], "/")
    if len(userPassword) != 2 || len(ipPortService) != 2 {
        return DatabaseConfig{}, port, fmt.Errorf("连接字符串格式错误")
    }
    // 解析 IP 地址和端口号部分
    ipPort := strings.Split(ipPortService[0], ":")
    if len(ipPort) == 2 {
        port, err := strconv.Atoi(ipPort[1])
        if err != nil {
            return DatabaseConfig{}, port, fmt.Errorf("连接字符串格式错误")
        }
    }
    // 设置数据库配置
    dbConfig.Username = userPassword[0]
    dbConfig.Password = userPassword[1]
    dbConfig.Server = ipPort[0]
    dbConfig.ServiceName = ipPortService[1]

    return dbConfig, port ,nil
}
最后执行的记过信息如下：

+-------------------------------------------------------------------------------------------------------------------+
| 数据库信息                                                                                                        |
|                                                                                                                   |
+------+----------------+------------+---------------------+---------------+------------+---------------------------+
| NAME | DB_UNIQUE_NAME | DBID       | DB_CREATED          | DATABASE_ROLE | OPEN_MODE  | NLS_LANG                  |
+------+----------------+------------+---------------------+---------------+------------+---------------------------+
| ORCL | orcl           | 1673481055 | 2023-08-24 10:05:19 | PRIMARY       | READ WRITE | AMERICAN_AMERICA.AL32UTF8 |
+------+----------------+------------+---------------------+---------------+------------+---------------------------+

+----------------------------------------------------------------------------------------------------------------------------------------------------------------+
| 检查表空间使用率                                                                                                                                               |
|                                                                                                                                                                |
+--------+-----------------+-----------+-------------------+--------------------------+----------+-----------------+----------+-------------------+--------------+
| STATUS | TABLESPACE_NAME | CONTENTS  | EXTENT_MANAGEMENT | SEGMENT_SPACE_MANAGEMENT | USED_MB  | CURRENT_SIZE_MB | PCT_USED | CANEXTEND_SIZE_MB | TOT_PCT_USED |
+--------+-----------------+-----------+-------------------+--------------------------+----------+-----------------+----------+-------------------+--------------+
| ONLINE | TEST_IDX        | PERMANENT | LOCAL             | AUTO                     | 1        | 50              |    2.00% | 32768             |    0.00%     |
| ONLINE | USERO_TBS       | PERMANENT | LOCAL             | AUTO                     | 1        | 50              |    2.00% | 50                |    2.00%     |
| ONLINE | HR_TBS          | PERMANENT | LOCAL             | AUTO                     | 2.5625   | 100             |    2.56% | 100               |    2.56%     |
| ONLINE | HR_TEMP         | TEMPORARY | LOCAL             | MANUAL                   | 97       | 100             |    3.00% | 100               |    3.00%     |
| ONLINE | TEST_TAB        | PERMANENT | LOCAL             | AUTO                     | 3        | 50              |    6.00% | 32768             |    0.01%     |
| ONLINE | USERS           | PERMANENT | LOCAL             | AUTO                     | 2.6875   | 5               |   53.75% | 32768             |    0.01%     |
| ONLINE | TEST_TBS        | PERMANENT | LOCAL             | AUTO                     | 432.5625 | 500             |   86.51% | 32768             |    1.32%     |
| ONLINE | SYSAUX          | PERMANENT | LOCAL             | AUTO                     | 1497.625 | 1580            |   94.79% | 32768             |    4.57%     |
| ONLINE | TEMP            | TEMPORARY | LOCAL             | MANUAL                   | 1        | 132             |   99.24% | 32768             |    0.40%     |
| ONLINE | SYSTEM          | PERMANENT | LOCAL             | MANUAL                   | 955.75   | 960             |   99.56% | 32768             |    2.92%     |
+--------+-----------------+-----------+-------------------+--------------------------+----------+-----------------+----------+-------------------+--------------+
总结
通过使用Go语言和一些开源包，我们实现了可配置化的数据库巡检，为DBA提供了更高效和灵活的工作方式。传统的巡检方法往往需要手动编写复杂的SQL查询和脚本来检查各种指标和配置，这不仅耗时繁琐，还容易出错。而通过使用go-ora开源包，我们可以轻松地连接、查询和分析数据库，大大简化了巡检过程。总而言之，借助Go语言和一些开源包，实现可配置化的数据库巡检为DBA提供了更高效、灵活和可靠的工作方式。这种方法不仅简化了巡检流程，还提升了数据库管理的效率和质量，为企业的数据库运维工作带来了更大的便利和价值。

文中涉及到代码如下: 链接：https://pan.baidu.com/s/1ul3m0RefPerUpLfU-2lWJg?pwd=6wsj 提取码：6wsj --来自百度网盘超级会员V5的分享

参考文献
Go标准库包：https://pkg.go.dev/std
database/sql包：https://pkg.go.dev/database/sql
go-ora Github: https://github.com/sijms/go-ora
go-ora pkg.go.dev: https://pkg.go.dev/github.com/sijms/go-ora/v2
yaml Github链接： https://github.com/go-yaml/yaml
yaml.v3 API文档：https://gopkg.in/yaml.v3
yaml.v3 pkg.go.dev: https://pkg.go.dev/gopkg.in/yaml.v3
go-pretty github链接：https://github.com/jedib0t/go-pretty
go-pretty v6 pkg.go.dev: https://pkg.go.dev/github.com/j
