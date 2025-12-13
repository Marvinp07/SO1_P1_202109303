module proyecto/go-daemon

go 1.21

require (
    github.com/docker/docker v24.0.7+incompatible
    gopkg.in/yaml.v3 v3.0.1
)

replace github.com/docker/distribution/reference => github.com/distribution/reference v0.6.0
