package cloudhypervisor

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"github.com/jonaprince/envy/virtualmachine"
)

type Client struct {
	httpClient *http.Client
}

func NewClient(socketPath string) *Client {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) CreateVM(virtualmachine *virtualmachine.Virtualmachine) error {
	// Construct the VM configuration
	vmConfig := VMConfig{
		Cpus: &CpusConfig{
			BootVcpus: virtualmachine.CPU,
			MaxVcpus:  virtualmachine.CPU,
		},
		Disks: []DiskConfig{
			{Path: "/var/lib/envy/vms/" + virtualmachine.ID + "/disk.img"},
		},
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(vmConfig)
	req, err := http.NewRequest("POST", "http://unix/vm.create", buf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type CpusConfig struct {
	BootVcpus int `json:"boot_vcpus"`
	MaxVcpus  int `json:"max_vcpus"`
}

type DiskConfig struct {
	Path string `json:"path"`
}

// type PayloadConfig struct {
// 	Cmdline   *string `json:"cmdline,omitempty"`
// 	Firmware  *string `json:"firmware,omitempty"`
// 	Initramfs *string `json:"initramfs,omitempty"`
// 	Kernel    *string `json:"kernel,omitempty"`
// }

type VMConfig struct {
	Cpus *CpusConfig `json:"cpus,omitempty"`
	// Payload *PayloadConfig `json:"payload,omitempty"`
	Disks []DiskConfig `json:"disks,omitempty"`
}
