import {createEnclave} from "../../test_helpers/enclave_setup";
import {
    DEFAULT_DRY_RUN,
    JEST_TIMEOUT_MS,
} from "./shared_constants";
import log from "loglevel";
import {KurtosisFeatureFlag} from "kurtosis-sdk/build/core/kurtosis_core_rpc_api_bindings/api_container_service_pb";

const DEFAULT_STARLARK_RUN_FUNC_NAME = "run"
const NO_EXPERIMENTAL_FEATURE = new Array<KurtosisFeatureFlag>()

const VALID_SCRIPT_INPUT_TEST_NAME = "valid-package-with-input"
const STARLARK_SCRIPT = `
def run(plan, args):
    plan.print(args["greetings"])
`
jest.setTimeout(JEST_TIMEOUT_MS)

test("Test valid Starlark script with input", async () => {
    // ------------------------------------- ENGINE SETUP ----------------------------------------------
    const createEnclaveResult = await createEnclave(VALID_SCRIPT_INPUT_TEST_NAME)

    if (createEnclaveResult.isErr()) {
        throw createEnclaveResult.error
    }

    const {enclaveContext, stopEnclaveFunction} = createEnclaveResult.value

    try {
        // ------------------------------------- TEST SETUP ----------------------------------------------
        const params = `{"greetings": "bonjour!"}`
        const runResult = await enclaveContext.runStarlarkScriptBlocking(
            DEFAULT_STARLARK_RUN_FUNC_NAME,
            STARLARK_SCRIPT,
            params,
            DEFAULT_DRY_RUN,
            NO_EXPERIMENTAL_FEATURE,
        )

        if (runResult.isErr()) {
            log.error(`An error occurred executing Starlark script`);
            throw runResult.error
        }

        expect(runResult.value.interpretationError).toBeUndefined()
        expect(runResult.value.validationErrors).toEqual([])
        expect(runResult.value.executionError).toBeUndefined()

        const expectedScriptOutput = "bonjour!\n"
        expect(runResult.value.runOutput).toEqual(expectedScriptOutput)
        expect(runResult.value.instructions).toHaveLength(1)
    } finally {
        stopEnclaveFunction()
    }
})

test("Test valid Starlark package with input - missing key in params", async () => {
    // ------------------------------------- ENGINE SETUP ----------------------------------------------
    const createEnclaveResult = await createEnclave(VALID_SCRIPT_INPUT_TEST_NAME)

    if (createEnclaveResult.isErr()) {
        throw createEnclaveResult.error
    }

    const {enclaveContext, stopEnclaveFunction} = createEnclaveResult.value

    try {
        // ------------------------------------- TEST SETUP ----------------------------------------------
        const params = `{"hello": "world"}` // expecting key 'greetings' here
        const runResult = await enclaveContext.runStarlarkScriptBlocking(
            DEFAULT_STARLARK_RUN_FUNC_NAME,
            STARLARK_SCRIPT,
            params,
            DEFAULT_DRY_RUN,
            NO_EXPERIMENTAL_FEATURE,
        )

        if (runResult.isErr()) {
            log.error(`An error occurred execute Starlark package`);
            throw runResult.error
        }

        expect(runResult.value.interpretationError).not.toBeUndefined()
        expect(runResult.value.interpretationError?.getErrorMessage()).toContain("Evaluation error: key \"greetings\" not in dict")
        expect(runResult.value.validationErrors).toEqual([])
        expect(runResult.value.executionError).toBeUndefined()

        expect(runResult.value.runOutput).toEqual("")
        expect(runResult.value.instructions).toHaveLength(0)
    } finally {
        stopEnclaveFunction()
    }
})
