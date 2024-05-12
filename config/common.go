package config

type Common struct {
	PackSize uint8  `json:"pack_size" yaml:"pack_size"`
	CacheDir string `json:"cache_dir" yaml:"cache_dir"`
	TmpDir   string `json:"tmp_dir" yaml:"tmp_dir"`
}
