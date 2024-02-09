package main

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

// 创建数据库连接
// osqlInfo := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", *username, *password, *oraclehost, *oracleport, *dbname)
// db, err := sql.Open("oracle", osqlInfo)

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