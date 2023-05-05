/*
解析gen.yaml文件并读取配置
解析gen.pro.yaml文件读取并覆盖配置
*/
package gen

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type GenConf struct {
	// model监听目录
	Dir string `yaml:"dir"`
	// 生成类型
	Codes []*GenCodeConf `yaml:"codes"`
}

type GenCodeConf struct {
	Typ    string `yaml:"typ"`
	OutDir string `yaml:"out_dir"`
}

// GenConf 生成配置
var genConf *GenConf

// GetGenConf 获取生成配置
func GetGenConf() *GenConf {
	return genConf
}

// ParseGenConf 解析生成配置
func ParseGenConf() {
	// 解析gen.yaml文件
	genConf = &GenConf{}
	err := parseYaml("gen.yaml", genConf)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("WARN gen.yaml is not found. use default config")
			initDefaultConf()
		} else {
			panic(err)
		}
		return
	}
	// 解析gen.pro.yaml文件
	err = parseYaml("gen.pro.yaml", genConf)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	// 检查配置
	checkYaml()
}

// parseYaml 解析yaml文件
func parseYaml(fileName string, conf *GenConf) error {
	// 解析yaml文件
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return err
	}
	return nil
}

// checkYaml 检查yaml文件
func checkYaml() {
	if genConf == nil {
		panic("genConf is nil")
	}
	if genConf.Dir == "" {
		panic("genConf.Dir is empty")
	}
	if genConf.Codes == nil {
		panic("genConf.Codes is nil")
	}
	if len(genConf.Codes) == 0 {
		panic("genConf.Codes is empty")
	}
	for _, code := range genConf.Codes {
		if code.Typ == "" {
			panic("genConf.Codes.Typ is empty")
		}
		if code.Typ != "chi_route" && code.Typ != "ts" && code.Typ != "dart" && code.Typ != "swift" {
			panic("genConf.Codes.Typ is invalid")
		}
		if code.OutDir == "" {
			panic("genConf.Codes.OutDir is empty")
		}
	}
}

// initDefaultConf init default config
func initDefaultConf() {
	genConf = &GenConf{
		Dir: "./internal/model",
		Codes: []*GenCodeConf{
			{
				Typ:    "chi_route",
				OutDir: "./internal/route",
			},
			{
				Typ:    "ts",
				OutDir: "./static/ts",
			},
		},
	}
}
