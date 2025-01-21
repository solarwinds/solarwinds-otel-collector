package installedsoftware

import (
	"fmt"
	"testing"

	"github.com/solarwinds/solarwinds-otel-collector/receiver/swohostmetricsreceiver/internal/cli"
	"github.com/stretchr/testify/assert"
)

const rpmPayload = `
libgcc;11.4.1;1705517422
crypto-policies;20230731;1705517422
tzdata;2023d;1705517422
subscription-manager-rhsm-certificates;20220623;1705517422
libreport-filesystem;2.15.2;1705517422
`

func Test_Rpm_GetSoftware_CommandFailsErrorIsReturned(t *testing.T) {
	cli := cli.CreateNewCliExecutorMock("", "", fmt.Errorf("failed to process command"))
	sut := createRpmProvider(cli)

	is, err := sut.GetSoftware()

	assert.Error(t, err, "call must produce error")
	assert.Zero(t, len(is), "no installed software must be provided")
}

func Test_Rpm_GetSoftware_CommandSucceedsSoftwareIsReturned(t *testing.T) {
	expectedSoftwareCount := 5

	cli := cli.CreateNewCliExecutorMock(rpmPayload, "", nil)
	sut := createRpmProvider(cli)

	is, err := sut.GetSoftware()

	assert.Equal(t, expectedSoftwareCount, len(is), "installed software must be in exact count")
	assert.NoError(t, err, "call must succeed")
	assert.Equal(t, is[0].Name, "libgcc", "the first installed software must be equal")
	assert.Equal(t, is[4].Name, "libreport-filesystem", "the last installed software must be equal")
}

func Test_Rpm_GetSoftware_UnsupportedFormatEmptyCollectionIsReturnedWithNoError(t *testing.T) {
	cli := cli.CreateNewCliExecutorMock("kokoha_unsupported", "", nil)
	sut := createRpmProvider(cli)

	is, err := sut.GetSoftware()

	assert.Zero(t, len(is), "installed software collection must be empty")
	assert.NoError(t, err, "no error is returned")
}
