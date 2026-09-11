package docker

const VOLUMES_PATH string = "/var/lib/docker/volumes"

type Container struct {
	Name        string
	Image       string
	Environment map[string]string
	Volumes     []string
	Ports       map[string]string
	DependsOn   []string
}
