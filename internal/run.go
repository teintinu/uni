package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/evanw/esbuild/pkg/api"
)

func Run(repo *Repository, entrypoint string, args []string) error {
	if err := EnsureTmp(repo); err != nil {
		return err
	}

	dir, err := TempDir(repo, "run")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entrypoint},
		Outfile:     path.Join(dir, "bundle.js"),
		Platform:    api.PlatformNode,
		Format:      api.FormatCommonJS,
		Write:       true,
		LogLevel:    api.LogLevelWarning,
	})

	if len(result.Errors) > 0 {
		return fmt.Errorf("build error")
	}

	script := `
		const { main } = require('./bundle.js');
		const args = process.argv.slice(2);
		void main(...args);
	`
	scriptPath := path.Join(dir, "script.js")
	if err := ioutil.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return err
	}

	nodeArgs := append([]string{scriptPath}, args...)
	node := exec.Command("node", nodeArgs...)
	node.Stdin = os.Stdin
	node.Stdout = os.Stdout
	node.Stderr = os.Stderr
	return node.Run()
}
