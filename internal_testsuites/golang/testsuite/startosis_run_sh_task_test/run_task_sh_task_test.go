package startosis_run_sh_task_test

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-cli/golang_internal_testsuite/test_helpers"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	runshTest           = "run-sh-test"
	runshStarlarkSimple = `
def run(plan):
  result1 = plan.run_sh(run="echo kurtosis")
  result2 = plan.run_sh(run="mkdir -p /src/{0} && cd /src/{0} && echo $(pwd)".format(result1.output))
  plan.assert(result2.output, "==", "/src/kurtosis\n")
`
	runshStarlarkFileArtifact = `
def run(plan):
  result = plan.run_sh(run="mkdir -p /src && echo kurtosis > /src/tech.txt", store=["/src/tech.txt", "/src"], image="ethpandaops/ethereum-genesis-generator:1.0.14")
  file_artifacts = result.files_artifacts
  result2 = plan.run_sh(run="cat /temp/tech.txt", files={"/temp": file_artifacts[0]})
  plan.assert(result2.output, "==", "kurtosis\n")
  result3 = plan.run_sh(run="cat /task/src/tech.txt", files={"/task": file_artifacts[1]})
  plan.assert(result3.output, "==", "kurtosis\n")
`

	runshStarlarkFileArtifactFailure = `
def run(plan):
  result = plan.run_sh(run="cat /tmp/kurtosis.txt")
  plan.assert(value=result.code, assertion="==", target_value="0")
`

	runshStarlarkWithTimeout = `
def run(plan):
  result = plan.run_sh(run="sleep 45s", wait="30s")
  plan.assert(value=result.code, assertion="==", target_value="0")
`

	runshStarlarkIgnoreParentDir = `
def run(plan):
  result1 = plan.run_sh(run="mkdir -p /src/data && echo kurtosis > /src/data/kurtosis.txt",store=["/src/*"])
  files_artifacts = result1.files_artifacts
  result2 = plan.run_sh(files={"/temp": files_artifacts[0]}, run="cat /temp/data/kurtosis.txt")
  plan.assert(result2.output, "==", "kurtosis\n")
`
)

func TestStarlark_RunshTaskSimple(t *testing.T) {
	ctx := context.Background()
	runResult, _ := test_helpers.SetupSimpleEnclaveAndRunScript(t, ctx, runshTest, runshStarlarkSimple)
	expectedOutput := "Command returned with exit code '0' and the following output:\n--------------------\nkurtosis\n\n--------------------\nCommand returned with exit code '0' and the following output:\n--------------------\n/src/kurtosis\n\n--------------------\nAssertion succeeded. Value is '\"/src/kurtosis\\n\"'.\n"
	require.Equal(t, expectedOutput, string(runResult.RunOutput))
}

func TestStarlark_RunshTaskFileArtifact(t *testing.T) {
	ctx := context.Background()
	runResult, _ := test_helpers.SetupSimpleEnclaveAndRunScript(t, ctx, runshTest, runshStarlarkFileArtifact)
	expectedOutput := "Command returned with exit code '0' with no output\nCommand returned with exit code '0' and the following output:\n--------------------\nkurtosis\n\n--------------------\nAssertion succeeded. Value is '\"kurtosis\\n\"'.\nCommand returned with exit code '0' and the following output:\n--------------------\nkurtosis\n\n--------------------\nAssertion succeeded. Value is '\"kurtosis\\n\"'.\n"
	require.Equal(t, expectedOutput, string(runResult.RunOutput))
}

func TestStarlark_RunshTaskFileArtifactFailure(t *testing.T) {
	ctx := context.Background()
	runResult, _ := test_helpers.SetupSimpleEnclaveAndRunScript(t, ctx, runshTest, runshStarlarkFileArtifactFailure)
	expectedErrorMessage := "cat: can't open '/tmp/kurtosis.txt': No such file or directory"
	require.NotNil(t, runResult.ExecutionError)
	require.Contains(t, runResult.ExecutionError.GetErrorMessage(), expectedErrorMessage)
}

func TestStarlark_RunshTimesoutSuccess(t *testing.T) {
	ctx := context.Background()
	runResult, _ := test_helpers.SetupSimpleEnclaveAndRunScript(t, ctx, runshTest, runshStarlarkWithTimeout)
	expectedErrorMessage := "The exec request timed out after 30 seconds"
	require.NotNil(t, runResult.ExecutionError)
	require.Contains(t, runResult.ExecutionError.GetErrorMessage(), expectedErrorMessage)
}

func TestStarlark_RunshFileArtifactWithoutParentDir(t *testing.T) {
	ctx := context.Background()
	runResult, _ := test_helpers.SetupSimpleEnclaveAndRunScript(t, ctx, runshTest, runshStarlarkIgnoreParentDir)
	expectedOutput := "Command returned with exit code '0' with no output\nCommand returned with exit code '0' and the following output:\n--------------------\nkurtosis\n\n--------------------\nAssertion succeeded. Value is '\"kurtosis\\n\"'.\n"
	require.Equal(t, expectedOutput, string(runResult.RunOutput))
}
