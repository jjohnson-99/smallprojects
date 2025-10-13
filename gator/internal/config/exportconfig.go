package config
import (
    "os"
    "log"
    "encoding/json"
)

type Config struct {
    DB_URL string `json:"db_url"`
    Current_user_name string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func Read() Config {
    var cfg Config

    filepath, _ := getConfigFilePath()
    data, err := os.ReadFile(filepath)
    if err != nil {
        return cfg
    }
    if err := json.Unmarshal(data, &cfg); err != nil {
        return cfg
    }
    return cfg
}

func (c Config) SetUser(name string) {
    c.Current_user_name = name
    err := write(c)
    if err != nil {
        log.Fatal(err)
    }
}

func getConfigFilePath() (string, error) {
    dir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    filepath := dir + configFileName
    return filepath, nil
}

func write(cfg Config) error {
    filename, _ := getConfigFilePath()
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    cfgJson, err := json.Marshal(cfg)
    if err != nil {
        return err
    }
    err = os.WriteFile(filename, cfgJson, 0644)
    if err != nil {
        log.Fatal(err)
    }
    return nil
}
