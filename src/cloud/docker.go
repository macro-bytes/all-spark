package cloud

import (
	"context"
	"daemon"
	"errors"
	"logger"
	"strconv"
	"time"
	"util/netutil"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

// DockerEnvironment interface
type DockerEnvironment struct {
	NanoCpus    int64
	MemBytes    int64
	ClusterID   string
	WorkerNodes int
	Image       string
	Mounts      []mount.Mount
	EnvParams   []string
}

const (
	allsparkBridgedNetwork = "allspark_bridged_newtork"
)

func (e *DockerEnvironment) getDockerClient() *client.Client {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.GetError().Println(err)
	}
	cli.NegotiateAPIVersion(ctx)
	return cli
}

func (e *DockerEnvironment) computeExecutorMemory() string {
	return strconv.FormatInt(e.MemBytes/1024/1024/1024-1, 10)
}

// CreateCluster - creates a spark cluster in docker
func (e *DockerEnvironment) CreateCluster() (string, error) {
	expectedWorkers := "EXPECTED_WORKERS=" + strconv.Itoa(e.WorkerNodes)

	var envVariables []string
	envVariables = []string{expectedWorkers,
		"SPARK_WORKER_PORT=7078",
		"CLUSTER_ID=" + e.ClusterID,
		"EXECUTOR_MEMORY=" + e.computeExecutorMemory(),
		"ALLSPARK_CALLBACK=" + daemon.GetAllSparkConfig().CallbackURL}

	envVariables = append(envVariables, e.EnvParams...)

	containerID, err := e.createSparkNode(e.ClusterID+masterIdentifier, envVariables)
	if err != nil {
		logger.GetError().Println(err)
	}

	masterIP, err := e.getIPAddress(containerID)
	if err != nil {
		logger.GetError().Println(err)
	}

	if netutil.IsListeningOnPort(masterIP, sparkPort, 30*time.Second, 120) {
		envVariables = append([]string{"MASTER_IP=" + masterIP,
			"SPARK_WORKER_PORT=7078"},
			envVariables...)

		for i := 1; i <= e.WorkerNodes; i++ {
			identifier := e.ClusterID + workerIdentifier + strconv.Itoa(i)
			e.createSparkNode(identifier, envVariables)
		}
	} else {
		return "", errors.New("master node has failed to come online")
	}

	webURL := "http://" + masterIP + ":8080"
	return webURL, nil
}

// DestroyCluster - destroys the spark cluster in docker
func (e *DockerEnvironment) DestroyCluster() error {
	cli := e.getDockerClient()
	defer cli.Close()

	clusterNodes, err := e.getClusterNodes()
	if err != nil {
		return err
	}

	for _, el := range clusterNodes {
		err = cli.ContainerRemove(context.Background(), el,
			types.ContainerRemoveOptions{Force: true})
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *DockerEnvironment) getClusterNodes() ([]string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	filters := filters.NewArgs()
	filters.Add("name", e.ClusterID)

	resp, err := cli.ContainerList(context.Background(),
		types.ContainerListOptions{Filters: filters})

	if err != nil {
		return nil, err
	}

	var result []string
	for _, el := range resp {
		result = append(result, el.Names[0])
	}
	return result, nil
}

func (e *DockerEnvironment) getIPAddress(id string) (string, error) {
	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return "", nil
	}

	return resp.NetworkSettings.Networks[allsparkBridgedNetwork].IPAddress, nil
}

func (e *DockerEnvironment) createSparkNode(identifier string,
	envParams []string) (string, error) {

	cli := e.getDockerClient()
	defer cli.Close()

	resp, err := cli.ContainerCreate(context.Background(),
		&container.Config{
			Image: e.Image,
			Env:   envParams,
		},
		&container.HostConfig{
			Resources: container.Resources{
				NanoCPUs: e.NanoCpus,
				Memory:   e.MemBytes,
			},
			Mounts:      e.Mounts,
			NetworkMode: allsparkBridgedNetwork,
		},
		&network.NetworkingConfig{},
		identifier)
	if err != nil {
		return "", err
	}

	if err = cli.ContainerStart(context.Background(),
		resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}
