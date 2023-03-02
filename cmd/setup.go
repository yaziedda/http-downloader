package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/linuxsuren/http-downloader/pkg/installer"
	"github.com/linuxsuren/http-downloader/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSetupCommand(v *viper.Viper, stdio terminal.Stdio) (cmd *cobra.Command) {
	opt := &setupOption{
		stdio: stdio,
		v:     v,
	}
	cmd = &cobra.Command{
		Use:     "setup",
		Short:   "Init the configuration of hd",
		RunE:    opt.runE,
		GroupID: configGroup.ID,
	}
	return
}

type setupOption struct {
	stdio terminal.Stdio
	v     *viper.Viper
}

func (o *setupOption) runE(cmd *cobra.Command, args []string) (err error) {
	var (
		proxyGitHub string
		provider    string
	)

	logger := log.GetLoggerFromContextOrDefault(cmd)
	if proxyGitHub, err = selectFromList([]string{"", "ghproxy.com", "gh.api.99988866.xyz", "mirror.ghproxy.com"},
		o.v.GetString("proxy-github"),
		"Select proxy-github", o.stdio); err == nil {
		o.v.Set("proxy-github", proxyGitHub)
	} else {
		return
	}

	if provider, err = selectFromList([]string{"github", "gitee"}, o.v.GetString("provider"),
		"Select provider", o.stdio); err == nil {
		o.v.Set("provider", provider)
	} else {
		return
	}

	var configDir string
	fetcher := &installer.DefaultFetcher{}
	if configDir, err = fetcher.GetHomeDir(); err != nil {
		return
	}
	if err = os.MkdirAll(configDir, 0750); err != nil {
		err = fmt.Errorf("failed to create directory: %s, error: %v", configDir, err)
		return
	}

	configPath := filepath.Join(configDir, ".config", "hd.yaml")
	logger.Info("write config into:", configPath)
	err = o.v.WriteConfigAs(configPath)
	return
}

func selectFromList(items []string, defaultItem, title string, stdio terminal.Stdio) (val string, err error) {
	selector := &survey.Select{
		Message: title,
		Options: items,
		Default: defaultItem,
	}
	err = survey.AskOne(selector, &val, survey.WithStdio(stdio.In, stdio.Out, stdio.Err))
	return
}
