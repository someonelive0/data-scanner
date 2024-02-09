package main


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
