package cloudhypervisor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
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

func (c *Client) CreateVM(vmConfig VMConfig) error {
	// Construct the VM configuration

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(vmConfig)
	req, err := http.NewRequest("PUT", "http://unix/api/v1/vm.create", buf)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		// Log the response body for debugging
		var respBody bytes.Buffer
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			slog.Error("failed to read response body", slog.Any("error", err))
		} else {
			slog.Error("failed to create VM", "response", respBody.String())
		}
		return fmt.Errorf("failed to create VM, status code: %d", resp.StatusCode)
	}
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
