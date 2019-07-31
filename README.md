# quickdemo

uses [demoinfocs](https://github.com/markus-wa/demoinfocs-golang) to retrieve stats required for [de_stats](https://github.com/grandprixgp/de_stats)

accepts demos as space seperated list of filenames with argument `-d`

chunks demos to parse as many as available memory allows

highly parallel (thanks to Go and [demoinfocs](https://github.com/markus-wa/demoinfocs-golang))
