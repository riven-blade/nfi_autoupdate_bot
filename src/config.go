package src

// FileInfo 定义与 YAML 结构相匹配的结构体
type FileInfo struct {
	FilePath       string `yaml:"file_path"`
	GithubFilePath string `yaml:"github_file_path"`
	HasVersion     bool   `yaml:"has_version"`
	Update         bool   `yaml:"update"`
}

type UpdateInfo struct {
	RestartAPI string     `yaml:"restart_api"`
	Name       string     `yaml:"name"`
	Version    string     `yaml:"version"`
	Files      []FileInfo `yaml:"files"`
}

type Config struct {
	UpdateInfos []UpdateInfo `yaml:"update_infos"`
	Github      string       `yaml:"github"`
	TgBot       string       `yaml:"tg_bot"`
	TgUserID    int64        `yaml:"tg_user_id"`
	Status      bool         `yaml:"status"`
	Username    string       `yaml:"username"`
	Password    string       `yaml:"password"`
}

type ConfigOutput struct {
	Name       string   `yaml:"name"`
	Version    string   `yaml:"version"`
	UpdateFile []string `yaml:"update_file"`
}
