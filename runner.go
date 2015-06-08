package main

import "github.com/hashicorp/consul/api"
import "encoding/json"
import "time"
import "errors"

type CheckState string

const (
	StateOK       CheckState = "ok"
	StateCritical CheckState = "critical"
)

type CheckResult struct {
	State    CheckState
	Metadata map[string]string
}

type CheckMetadata map[string]string
type Checks map[string]CheckMetadata

func runner() {
	checkInterval := time.Tick(time.Minute)

	for {
		select {
		case <-checkInterval:
			checks, err := getChecks()
			if err != nil {
				log.Printf("failed to retrieve checks to run: %s", err)
				continue
			}

			for checkName, checkMeta := range checks {
				err = runCheck(checkName, checkMeta)
				if err != nil {
					log.Printf("failed to run check '%s': %s", checkName, err)
				}
			}
		}
	}
}

func getChecks() (Checks, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	kv := client.KV()

	// Get all of the keys, which correspond to checks to run.
	checkNames, _, err := kv.Keys("/hcf/checks", "", nil)
	if err != nil {
		return nil, err
	}

	checks := make(Checks)
	for _, checkName := range checkNames {
		check, _, err := kv.Get("/hcf/checks/"+checkName, nil)
		if err != nil {
			return nil, err
		}

		checkMeta := make(CheckMetadata)
		err = json.Unmarshal(check.Value, checkMeta)
		if err != nil {
			return nil, err
		}

		checks[checkName] = checkMeta
	}

	return checks, nil
}

func runCheck(checkName string, checkMeta CheckMetadata) error {
	checkType, ok := checkMeta["type"]
	if !ok {
		return errors.New("check must specify a type!")
	}

	var checkResult CheckResult
	switch checkType {
	case "mysql_alive":
		checkResult, err = runMysqlAliveCheck(checkMeta)
		if err != nil {
			return fmt.Errorf("error during mysql_alive check: %s", err)
		}
	}

	err = submitResult(checkName, checkResult)
	if err != nil {
		return fmt.Errorf("error submitting result: %s", err)
	}

	return nil
}

func submitResult(checkName string, checkResult CheckResult) error {
	return nil
}
