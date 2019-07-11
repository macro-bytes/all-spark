package template

type DockerTemplate struct {
	NanoCpus    int64
	MemBytes    int64
	Network     string
	ClusterID   string
	WorkerNodes int
}
