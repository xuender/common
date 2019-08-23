package common

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command 命令行
type Command struct {
	cmd *cobra.Command
}

// NewCommand 新疆命令行参数
func NewCommand(cmd *cobra.Command) *Command {
	return &Command{cmd: cmd}
}

// Int 读取配置Int
func (c *Command) Int(cmd *cobra.Command, name string) int {
	f := cmd.Flag(name)
	// 命令行优先
	if f.Changed {
		if i, err := strconv.Atoi(f.Value.String()); err == nil {
			return i
		}
	}
	ret := viper.GetInt(name)
	if ret == 0 {
		if i, err := strconv.Atoi(f.Value.String()); err == nil {
			return i
		}
	}
	return ret
}

// String 读取配置String
func (c *Command) String(cmd *cobra.Command, name string) string {
	f := cmd.Flag(name)
	// 命令行优先
	if f.Changed {
		return f.Value.String()
	}
	ret := viper.GetString(name)
	if ret == "" {
		return f.Value.String()
	}
	return ret
}

// Bool 读取配置String
func (c *Command) Bool(cmd *cobra.Command, name string) bool {
	f := cmd.Flag(name)
	b, _ := strconv.ParseBool(f.Value.String())
	// 命令行优先
	if f.Changed {
		return b
	}
	ret := viper.GetString(name)
	if ret == "" {
		return b
	}
	return viper.GetBool(name)
}
