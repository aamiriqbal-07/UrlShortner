package config

type Config struct {
    Server struct {
        Port string `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
    
    Database struct {
        Host     string `yaml:"host"`
        Port     string `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    
    ShortURL struct {
        Length    int    `yaml:"length"`
        BaseURL   string `yaml:"base_url"`
    } `yaml:"short_url"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    // Default values
    cfg.Server.Port = "8080"
    cfg.Server.Host = "localhost"
    
    cfg.Database.Host = "localhost"
    cfg.Database.Port = "3306"
    cfg.Database.User = "root"
    cfg.Database.Password = "12345"
    cfg.Database.Name = "urlshortner"
    
    cfg.ShortURL.Length = 6
    cfg.ShortURL.BaseURL = "http://localhost:8080"
    
    return cfg, nil
}