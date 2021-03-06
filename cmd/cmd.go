package cmd

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/kataras/iris-cli/utils"

	"github.com/spf13/cobra"
)

var proxyAddr string

// New returns the root command.
func New(buildRevision, buildTime string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "iris-cli",
		Short: "Command Line Interface for Iris",
		Long: `Iris CLI is a tool for Iris Web Framework.
It can be used to install starter kits and project structures 
Complete documentation is available at https://github.com/kataras/iris-cli`,
		SilenceErrors:              true,
		SilenceUsage:               true,
		TraverseChildren:           true,
		SuggestionsMinimumDistance: 1,
		Run:                        func(cmd *cobra.Command, args []string) {},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			transport := &http.Transport{
				DisableCompression: true,
				DisableKeepAlives:  true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: utils.IsInsideDocker(),
				},
			}

			if proxyAddr != "" {
				if proxyAddr == "env" {
					transport.Proxy = http.ProxyFromEnvironment
				} else {
					u := &url.URL{Scheme: "http", Host: proxyAddr}
					transport.Proxy = func(req *http.Request) (*url.URL, error) {
						return u, nil
					}
				}
			}

			utils.DefaultClient.Transport = transport
		},
	}

	helpTemplate := HelpTemplate{
		BuildRevision:        buildRevision,
		BuildTime:            buildTime,
		ShowGoRuntimeVersion: true,
	}
	rootCmd.SetHelpTemplate(helpTemplate.String())
	rootCmd.PersistentFlags().StringVar(&proxyAddr, "proxy", proxyAddr, "--proxy=env to load from system or ip:port form, e.g. 51.158.178.4:3128")

	// Commands.
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(newCommand())
	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(cleanCommand())
	rootCmd.AddCommand(unistallCommand())
	rootCmd.AddCommand(addCommand())
	rootCmd.AddCommand(checkCommand())

	return rootCmd
}

var shared = make(map[string]map[string]interface{}) // key = root command/app and value a map of key-value pair.

// SetValue sets a value to the shared store for specific app based on the root "cmd".
func SetValue(cmd *cobra.Command, key string, value interface{}) {
	name := cmd.Root().Name()

	m := shared[name]
	if m == nil {
		m = make(map[string]interface{})
		shared[name] = m
	}

	m[key] = value
}

// GetValue retrieves a value from the shared store from a specific app based on the root "cmd".
func GetValue(cmd *cobra.Command, key string) (interface{}, bool) {
	if m, ok := shared[cmd.Root().Name()]; ok {
		v := m[key]
		if v != nil {
			return m, true
		}
	}

	return nil, false
}

// RunCommand runs a command.
func RunCommand(from *cobra.Command, commandToRun string, args ...string) error {
	cmd, _, err := from.Root().Find(append([]string{commandToRun}, args...))
	if err != nil {
		return err
	}

	if err = cmd.ParseFlags(args); err != nil {
		return err
	}

	if fn := cmd.PreRunE; fn != nil {
		if err = fn(cmd, args); err != nil {
			return err
		}
	}

	if err = cmd.RunE(cmd, args); err != nil {
		return err
	}

	if fn := cmd.PostRunE; fn != nil {
		return fn(cmd, args)
	}

	return nil
}

func readDataFile(cmd *cobra.Command, path string, ptr interface{}) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
		return
	}

	err = json.Unmarshal(b, &ptr)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
