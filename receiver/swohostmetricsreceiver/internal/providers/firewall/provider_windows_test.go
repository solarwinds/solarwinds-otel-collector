package firewall

import (
	"fmt"
	"testing"

	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/registry"

	"github.com/stretchr/testify/assert"
)

func Test_Provide_ProvidesFirewallContainerOnSucceedingCliCommand(t *testing.T) {
	expectedModel := Container{
		FirewallProfiles: []Profile{
			{Name: "Domain", Enabled: 1},
			{Name: "Private", Enabled: 0},
			{Name: "Public", Enabled: 1},
		},
		Error: nil,
	}

	sut := provider{
		getRegistryValue: getMockRegistryValueFunc(map[string]uint64{
			"StandardProfile": 0,
			"DomainProfile":   1,
			"PublicProfile":   1,
		}),
	}

	ch := sut.Provide()
	actualModel := <-ch
	_, open := <-ch

	assert.Equal(t, expectedModel, actualModel, "exact model must be provided")
	assert.False(t, open, "channel must be closed at the end")
}

func Test_Provide_ProvidesErrorfullFirewallContainerOnFailingCommandWithPartialResult(t *testing.T) {
	expectedProfiles := []Profile{
		{Name: "Domain", Enabled: 1},
		{Name: "Public", Enabled: 0},
	}

	sut := provider{
		getRegistryValue: getMockRegistryValueFunc(map[string]uint64{
			"DomainProfile": 1,
			"PublicProfile": 0,
		}),
	}

	ch := sut.Provide()
	actualModel := <-ch
	_, open := <-ch

	assert.ElementsMatch(t, expectedProfiles, actualModel.FirewallProfiles, "partial profiles must be provided")
	assert.ErrorContains(t, actualModel.Error, "failed to get profile 'StandardProfile'", "error must be provided in result")
	assert.False(t, open, "channel must be closed at the end")
}

func getMockRegistryValueFunc(results map[string]uint64) registry.GetKeyUIntValueTypeFunc {
	return func(_ registry.RootKey, _, keyName string, _ string) (uint64, error) {
		result, exists := results[keyName]
		if !exists {
			return 0, fmt.Errorf("test failure getting %s", keyName)
		}
		return result, nil
	}
}
