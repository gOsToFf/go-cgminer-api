package cgminer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type CGMiner struct {
	server string
}

type status struct {
	Code        int
	Description string
	Status      string `json:"STATUS"`
	When        int64
}

type Summary struct {
	Accepted               int64
	BestShare              int64   `json:"Best Share"`
	DeviceHardwarePercent  float64 `json:"Device Hardware%"`
	DeviceRejectedPercent  float64 `json:"Device Rejected%"`
	DifficultyAccepted     float64 `json:"Difficulty Accepted"`
	DifficultyRejected     float64 `json:"Difficulty Rejected"`
	DifficultyStale        float64 `json:"Difficulty Stale"`
	Discarded              int64
	Elapsed                int64
	FoundBlocks            int64 `json:"Found Blocks"`
	GetFailures            int64 `json:"Get Failures"`
	Getworks               int64
	HardwareErrors         int64   `json:"Hardware Errors"`
	LocalWork              int64   `json:"Local Work"`
	MHS5s                  float64 `json:"MHS 5s"`
	MHSav                  float64 `json:"MHS av"`
	NetworkBlocks          int64   `json:"Network Blocks"`
	PoolRejectedPercentage float64 `json:"Pool Rejected%"`
	PoolStalePercentage    float64 `json:"Pool Stale%"`
	Rejected               int64
	RemoteFailures         int64 `json:"Remote Failures"`
	Stale                  int64
	TotalMH                float64 `json:"Total MH"`
	Utilty                 float64
	WorkUtility            float64 `json:"Work Utility"`
}

type statsTempResponse struct {
	StatsSum []byte  `json:"SUMMARY"`
	Stats    []Stats `json:"SUMMARY"`
}

type statsResponse struct {
	Status []status `json:"STATUS"`
	Stats  []Stats  `json:"STATS"`
	Id     int64    `json:"id"`
}

type Stats struct {
	Elapsed        int64
	Ghs5s          string  `json:"GHS 5s"`
	Ghsav          float64 `json:"GHS av"`
	Frequency      string  `json:"frequency"`
	Temp1          float64
	Temp2          float64
	Temp3          float64
	Temp4          float64
	Temp5          float64
	Temp6          float64
	Temp7          float64
	Temp8          float64
	Temp21         float64 `json:"temp2_1"`
	Temp22         float64 `json:"temp2_2"`
	Temp23         float64 `json:"temp2_3"`
	Temp24         float64 `json:"temp2_4"`
	Temp25         float64 `json:"temp2_5"`
	Temp26         float64 `json:"temp2_6"`
	Temp27         float64 `json:"temp2_7"`
	Temp28         float64 `json:"temp2_8"`
	ChanRate1      string  `json:"chain_rate1"`
	ChanRate2      string  `json:"chain_rate2"`
	ChanRate3      string  `json:"chain_rate3"`
	ChanRate4      string  `json:"chain_rate4"`
	ChanRate5      string  `json:"chain_rate5"`
	ChanRate6      string  `json:"chain_rate6"`
	ChanRate7      string  `json:"chain_rate7"`
	ChanRate8      string  `json:"chain_rate8"`
	ChanIdealRate6 float64 `json:"chain_rateideal6"`
	ChanIdealRate7 float64 `json:"chain_rateideal7"`
	ChanIdealRate8 float64 `json:"chain_rateideal8"`
}

type Devs struct {
	GPU                 int64
	Enabled             string
	Status              string
	Temperature         float64
	FanSpeed            int     `json:"Fan Speed"`
	FanPercent          int64   `json:"Fan Percent"`
	GPUClock            int64   `json:"GPU Clock"`
	MemoryClock         int64   `json:"Memory Clock"`
	GPUVoltage          float64 `json:"GPU Voltage"`
	Powertune           int64
	MHSav               float64 `json:"MHS av"`
	MHS5s               float64 `json:"MHS 5s"`
	Accepted            int64
	Rejected            int64
	HardwareErrors      int64 `json:"Hardware Errors"`
	Utility             float64
	Intensity           string
	LastSharePool       int64   `json:"Last Share Pool"`
	LashShareTime       int64   `json:"Lash Share Time"`
	TotalMH             float64 `json:"TotalMH"`
	Diff1Work           int64   `json:"Diff1 Work"`
	DifficultyAccepted  float64 `json:"Difficulty Accepted"`
	DifficultyRejected  float64 `json:"Difficulty Rejected"`
	LastShareDifficulty float64 `json:"Last Share Difficulty"`
	LastValidWork       int64   `json:"Last Valid Work"`
	DeviceHardware      float64 `json:"Device Hardware%"`
	DeviceRejected      float64 `json:"Device Rejected%"`
	DeviceElapsed       int64   `json:"Device Elapsed"`
}

type Pool struct {
	Accepted               int64
	BestShare              int64   `json:"Best Share"`
	Diff1Shares            int64   `json:"Diff1 Shares"`
	DifficultyAccepted     float64 `json:"Difficulty Accepted"`
	DifficultyRejected     float64 `json:"Difficulty Rejected"`
	DifficultyStale        float64 `json:"Difficulty Stale"`
	Discarded              int64
	GetFailures            int64 `json:"Get Failures"`
	Getworks               int64
	HasGBT                 bool    `json:"Has GBT"`
	HasStratum             bool    `json:"Has Stratum"`
	LastShareDifficulty    float64 `json:"Last Share Difficulty"`
	LastShareTime          int64   `json:"Last Share Time"`
	LongPoll               string  `json:"Long Poll"`
	Pool                   int64   `json:"POOL"`
	PoolRejectedPercentage float64 `json:"Pool Rejected%"`
	PoolStalePercentage    float64 `json:"Pool Stale%"`
	Priority               int64
	ProxyType              string `json:"Proxy Type"`
	Proxy                  string
	Quota                  int64
	Rejected               int64
	RemoteFailures         int64 `json:"Remote Failures"`
	Stale                  int64
	Status                 string
	StratumActive          bool   `json:"Stratum Active"`
	StratumURL             string `json:"Stratum URL"`
	URL                    string
	User                   string
	Works                  int64
}

type summaryResponse struct {
	Status  []status  `json:"STATUS"`
	Summary []Summary `json:"SUMMARY"`
	Id      int64     `json:"id"`
}

type devsResponse struct {
	Status []status `json:"STATUS"`
	Devs   []Devs   `json:"DEVS"`
	Id     int64    `json:"id"`
}

type poolsResponse struct {
	Status []status `json:"STATUS"`
	Pools  []Pool   `json:"POOLS"`
	Id     int64    `json:"id"`
}

type addPoolResponse struct {
	Status []status `json:"STATUS"`
	Id     int64    `json:"id"`
}

// New returns a CGMiner pointer, which is used to communicate with a running
// CGMiner instance. Note that New does not attempt to connect to the miner.
func New(hostname string, port int64) *CGMiner {
	miner := new(CGMiner)
	server := fmt.Sprintf("%s:%d", hostname, port)
	miner.server = server

	return miner
}

func (miner *CGMiner) runCommand(command, argument string) (string, error) {
	conn, err := net.DialTimeout("tcp", miner.server, 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	type commandRequest struct {
		Command   string `json:"command"`
		Parameter string `json:"parameter,omitempty"`
	}

	request := &commandRequest{
		Command: command,
	}

	if argument != "" {
		request.Parameter = argument
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		_ = conn.Close()
		return "", err
	}

	fmt.Fprintf(conn, "%s", requestBody)
	result, err := bufio.NewReader(conn).ReadString('\x00')
	_ = conn.Close()
	if err != nil {
		return "", err
	}

	result = strings.ReplaceAll(result, string('\n'), "")
	result = strings.ReplaceAll(result, "}{", "},{")

	return strings.TrimRight(result, "\x00"), nil
}

// Devs returns basic information on the miner. See the Devs struct.
func (miner *CGMiner) Devs() (*[]Devs, error) {
	result, err := miner.runCommand("devs", "")
	if err != nil {
		return nil, err
	}

	var devsResponse devsResponse
	err = json.Unmarshal([]byte(result), &devsResponse)
	if err != nil {
		return nil, err
	}

	var devs = devsResponse.Devs
	return &devs, err
}

// Summary returns basic information on the miner. See the Summary struct.
func (miner *CGMiner) Stats() (*Stats, error) {
	result, err := miner.runCommand("stats", "")
	if err != nil {
		return nil, err
	}

	var statsResponse statsResponse
	err = json.Unmarshal([]byte(result), &statsResponse)
	if err != nil {
		return nil, err
	}

	if len(statsResponse.Stats) != 1 {
		var stats = statsResponse.Stats[1]
		return &stats, err
	}

	var stats = statsResponse.Stats[0]
	return &stats, err
}

// Summary returns basic information on the miner. See the Summary struct.
func (miner *CGMiner) Summary() (*Summary, error) {
	result, err := miner.runCommand("summary", "")
	if err != nil {
		return nil, err
	}

	var summaryResponse summaryResponse
	err = json.Unmarshal([]byte(result), &summaryResponse)
	if err != nil {
		return nil, err
	}

	if len(summaryResponse.Summary) != 1 {
		return nil, errors.New("Received multiple Summary objects")
	}

	var summary = summaryResponse.Summary[0]
	return &summary, err
}

// Pools returns a slice of Pool structs, one per pool.
func (miner *CGMiner) Pools() ([]Pool, error) {
	result, err := miner.runCommand("pools", "")
	if err != nil {
		return nil, err
	}

	var poolsResponse poolsResponse
	err = json.Unmarshal([]byte(result), &poolsResponse)
	if err != nil {
		return nil, err
	}

	var pools = poolsResponse.Pools
	return pools, nil
}

// AddPool adds the given URL/username/password combination to the miner's
// pool list.
func (miner *CGMiner) AddPool(url, username, password string) error {
	// TODO: Don't allow adding a pool that's already in the pool list
	// TODO: Escape commas in the URL, username, and password
	parameter := fmt.Sprintf("%s,%s,%s", url, username, password)
	result, err := miner.runCommand("addpool", parameter)
	if err != nil {
		return err
	}

	var addPoolResponse addPoolResponse
	err = json.Unmarshal([]byte(result), &addPoolResponse)
	if err != nil {
		// If there an error here, it's possible that the pool was actually added
		return err
	}

	status := addPoolResponse.Status[0]

	if status.Status != "S" {
		return errors.New(fmt.Sprintf("%d: %s", status.Code, status.Description))
	}

	return nil
}

func (miner *CGMiner) Enable(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("enablepool", parameter)
	return err
}

func (miner *CGMiner) Disable(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("disablepool", parameter)
	return err
}

func (miner *CGMiner) Delete(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("removepool", parameter)
	return err
}

func (miner *CGMiner) SwitchPool(pool *Pool) error {
	parameter := fmt.Sprintf("%d", pool.Pool)
	_, err := miner.runCommand("switchpool", parameter)
	return err
}

func (miner *CGMiner) Restart() error {
	_, err := miner.runCommand("restart", "")
	return err
}

func (miner *CGMiner) Quit() error {
	_, err := miner.runCommand("quit", "")
	return err
}
