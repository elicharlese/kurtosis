package startosis_remote_package_test

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-cli/golang_internal_testsuite/test_helpers"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/kurtosis_core_rpc_api_bindings"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testName               = "package"
	defaultDryRun          = false
	remotePackage          = "github.com/kurtosis-tech/datastore-army-package"
	executeParams          = `{"num_datastores": 2}`
	dataStoreService0Name  = "datastore-0"
	dataStoreService1Name  = "datastore-1"
	datastorePortId        = "grpc"
	defaultParallelism     = 4
	useDefaultMainFile     = ""
	useDefaultFunctionName = ""
)

var (
	noExperimentalFeature = []kurtosis_core_rpc_api_bindings.KurtosisFeatureFlag{}
)

func TestStartosisRemotePackage(t *testing.T) {
	ctx := context.Background()

	// ------------------------------------- ENGINE SETUP ----------------------------------------------
	enclaveCtx, _, destroyEnclaveFunc, err := test_helpers.CreateEnclave(t, ctx, testName)
	require.NoError(t, err, "An error occurred creating an enclave")
	defer func() {
		err = destroyEnclaveFunc()
		require.NoError(t, err, "An error occurred destroying the enclave after the test finished")
	}()

	// ------------------------------------- TEST RUN ----------------------------------------------
	logrus.Debugf("Executing Starlark Package: '%v'", remotePackage)

	runResult, err := enclaveCtx.RunStarlarkRemotePackageBlocking(ctx, remotePackage, useDefaultMainFile, useDefaultFunctionName, executeParams, defaultDryRun, defaultParallelism, noExperimentalFeature)
	require.NoError(t, err, "Unexpected error executing starlark package")

	require.Nil(t, runResult.InterpretationError, "Unexpected interpretation error. This test requires you to be online for the read_file command to run")
	require.Empty(t, runResult.ValidationErrors, "Unexpected validation error")
	require.Empty(t, runResult.ExecutionError, "Unexpected execution error")
	logrus.Infof("Successfully ran Starlark Package")

	// Check that the service added by the script is functional
	logrus.Infof("Checking that services are all healthy")
	require.NoError(
		t,
		test_helpers.ValidateDatastoreServiceHealthy(context.Background(), enclaveCtx, dataStoreService0Name, datastorePortId),
		"Error validating datastore server '%s' is healthy",
		dataStoreService0Name,
	)
	require.NoError(
		t,
		test_helpers.ValidateDatastoreServiceHealthy(context.Background(), enclaveCtx, dataStoreService1Name, datastorePortId),
		"Error validating datastore server '%s' is healthy",
		dataStoreService1Name,
	)
	logrus.Infof("All services added via the package work as expected")
}
