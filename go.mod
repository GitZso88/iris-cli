module github.com/kataras/iris-cli

go 1.13

replace github.com/AlecAivazis/survey/v2 => github.com/kataras/survey/v2 v2.0.6

require (
	github.com/AlecAivazis/survey/v2 v2.0.8
	github.com/cheggaaa/pb/v3 v3.0.4
	github.com/creack/pty v1.1.11
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/denormal/go-gitignore v0.0.0-20180930084346-ae8ad1d07817
	github.com/fsnotify/fsnotify v1.4.9
	github.com/kataras/golog v0.0.18
	github.com/kataras/neffos v0.0.16
	github.com/spf13/cobra v1.0.0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
