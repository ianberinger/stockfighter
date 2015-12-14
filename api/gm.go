package api

import "fmt"

const (
	gmURL string = "https://www.stockfighter.io/gm/"
)

//LevelState contains all data returned by API call to GameMaster-API.
type LevelState struct {
	ErrorResult
	Account              string       `json:"account"`
	InstanceID           int          `json:"instanceId"`
	Instructions         Instructions `json:"instructions"`
	SecondsPerTradingDay int          `json:"secondsPerTradingDay"`
	Venues               []string     `json:"venues"`
	Symbols              []string     `json:"tickers"`
}

//Instructions returned from the stockfighter GameMaster-API.
type Instructions struct {
	Instructions string `json:"Instructions"`
	OrderTypes   string `json:"Order Types"`
}

//StartLevel starts a level and sets the instance state. It returns the levelstate.
func (i *Instance) StartLevel(level string) (v LevelState) {
	i.doHTTP("POST", fmt.Sprintf("%slevels/%s", gmURL, level), nil, &v)

	if v.InstanceID != 0 {
		i.setState(v.InstanceID, v.Account, v.Venues[0], v.Symbols[0])
	}
	return
}

//RestartLevel restarts the current level. An instanceID needs to already set.
func (i *Instance) RestartLevel() (v LevelState) {
	i.doHTTP("POST", fmt.Sprintf("%sinstances/%d/restart", gmURL, i.instanceID), nil, &v)

	if v.InstanceID != 0 {
		i.setState(v.InstanceID, v.Account, v.Venues[0], v.Symbols[0])
	}
	return
}

//StopLevel stops the current level and clears the instance state. An instanceID needs to already set.
func (i *Instance) StopLevel() (v ErrorResult) {
	i.doHTTP("POST", fmt.Sprintf("%sinstances/%d/stop", gmURL, i.instanceID), nil, &v)

	i.setState(i.instanceID, "", "", "")
	return
}

//ResumeLevel resumes the current level (and sets the instance state) and returns the levelstate. An instanceID needs to already set.
func (i *Instance) ResumeLevel() (v LevelState) {
	i.doHTTP("POST", fmt.Sprintf("%sinstances/%d/resume", gmURL, i.instanceID), nil, &v)

	if v.InstanceID != 0 {
		i.setState(v.InstanceID, v.Account, v.Venues[0], v.Symbols[0])
	}
	return
}

//JudgeLevel tells the API to judge the current level and clears the instance state. An instanceID needs to already set.
func (i *Instance) JudgeLevel() (v ErrorResult) {
	i.doHTTP("POST", fmt.Sprintf("%sinstances/%d/judge", gmURL, i.instanceID), nil, &v)

	i.setState(i.instanceID, "", "", "")
	return
}
