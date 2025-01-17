package test_engine

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis/container-engine-lib/lib/backend_interface/objects/container_status"
	"github.com/kurtosis-tech/kurtosis/container-engine-lib/lib/backend_interface/objects/port_spec"
	"github.com/kurtosis-tech/kurtosis/container-engine-lib/lib/backend_interface/objects/service"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/service_network"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_instruction/add_service"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_starlark_framework/kurtosis_plan_instruction"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/kurtosis_types/service_config"
	"github.com/kurtosis-tech/kurtosis/core/server/api_container/server/startosis_engine/runtime_value_store"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.starlark.net/starlark"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type addServicesTestCase struct {
	*testing.T
	serviceNetwork    *service_network.MockServiceNetwork
	runtimeValueStore *runtime_value_store.RuntimeValueStore
}

func (suite *KurtosisPlanInstructionTestSuite) TestAddServices() {
	suite.serviceNetwork.EXPECT().ExistServiceRegistration(TestServiceName).Times(1).Return(false, nil)
	suite.serviceNetwork.EXPECT().ExistServiceRegistration(TestServiceName2).Times(1).Return(false, nil)
	suite.serviceNetwork.EXPECT().UpdateServices(
		mock.Anything,
		map[service.ServiceName]*service.ServiceConfig{},
		mock.Anything,
	).Times(1).Return(
		map[service.ServiceName]*service.Service{},
		map[service.ServiceName]error{},
		nil,
	)
	suite.serviceNetwork.EXPECT().AddServices(
		mock.Anything,
		mock.MatchedBy(func(configs map[service.ServiceName]*service.ServiceConfig) bool {
			suite.Require().Len(configs, 2)
			suite.Require().Contains(configs, TestServiceName)
			suite.Require().Contains(configs, TestServiceName2)

			expectedServiceConfig1 := service.NewServiceConfig(
				TestContainerImageName,
				map[string]*port_spec.PortSpec{},
				map[string]*port_spec.PortSpec{},
				nil,
				nil,
				map[string]string{},
				nil,
				nil,
				0,
				0,
				service_config.DefaultPrivateIPAddrPlaceholder,
				0,
				0,
			)
			actualServiceConfig1 := configs[TestServiceName]
			suite.Assert().Equal(expectedServiceConfig1, actualServiceConfig1)

			expectedServiceConfig2 := service.NewServiceConfig(
				TestContainerImageName,
				map[string]*port_spec.PortSpec{},
				map[string]*port_spec.PortSpec{},
				nil,
				nil,
				map[string]string{},
				nil,
				nil,
				TestCpuAllocation,
				TestMemoryAllocation,
				service_config.DefaultPrivateIPAddrPlaceholder,
				0,
				0,
			)
			actualServiceConfig2 := configs[TestServiceName2]
			suite.Assert().Equal(expectedServiceConfig2, actualServiceConfig2)
			return true
		}),
		mock.Anything,
	).Times(1).Return(
		map[service.ServiceName]*service.Service{
			TestServiceName:  service.NewService(service.NewServiceRegistration(TestServiceName, TestServiceUuid, TestEnclaveUuid, nil, string(TestServiceName)), container_status.ContainerStatus_Running, nil, nil, nil),
			TestServiceName2: service.NewService(service.NewServiceRegistration(TestServiceName2, TestServiceUuid2, TestEnclaveUuid, nil, string(TestServiceName2)), container_status.ContainerStatus_Running, nil, nil, nil),
		},
		map[service.ServiceName]error{},
		nil,
	)

	suite.serviceNetwork.EXPECT().HttpRequestService(
		mock.Anything,
		string(TestServiceName),
		TestReadyConditionsRecipePortId,
		TestGetRequestMethod,
		"",
		TestReadyConditionsRecipeEndpoint,
		"",
	).Times(1).Return(&http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Request: &http.Request{ //nolint:exhaustruct
			Method: TestGetRequestMethod,
			URL:    &url.URL{}, //nolint:exhaustruct
		},
		Close:            true,
		ContentLength:    -1,
		Body:             io.NopCloser(strings.NewReader("{}")),
		Trailer:          nil,
		TransferEncoding: nil,
		Uncompressed:     true,
		TLS:              nil,
	}, nil)

	suite.serviceNetwork.EXPECT().HttpRequestService(
		mock.Anything,
		string(TestServiceName2),
		TestReadyConditions2RecipePortId,
		TestGetRequestMethod,
		"",
		TestReadyConditions2RecipeEndpoint,
		"",
	).Times(1).Return(&http.Response{
		Status:     "201 OK",
		StatusCode: 201,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{},
		Request: &http.Request{
			Method: TestGetRequestMethod,
			URL: &url.URL{ //nolint:exhaustruct
				Path:        "",
				Scheme:      "",
				Opaque:      "",
				User:        nil,
				Host:        "",
				RawPath:     "",
				ForceQuery:  false,
				RawQuery:    "",
				Fragment:    "",
				RawFragment: "",
			},
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           http.Header{},
			Body:             nil,
			GetBody:          nil,
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            false,
			Host:             "",
			Form:             nil,
			PostForm:         nil,
			MultipartForm:    nil,
			Trailer:          nil,
			RemoteAddr:       "",
			RequestURI:       "",
			TLS:              nil,
			Cancel:           nil,
			Response:         nil,
		},
		Close:            true,
		ContentLength:    -1,
		Body:             io.NopCloser(strings.NewReader("{}")),
		Trailer:          nil,
		TransferEncoding: nil,
		Uncompressed:     true,
		TLS:              nil,
	}, nil)

	suite.run(&addServicesTestCase{
		T:                 suite.T(),
		serviceNetwork:    suite.serviceNetwork,
		runtimeValueStore: suite.runtimeValueStore,
	})
}

func (t *addServicesTestCase) GetInstruction() *kurtosis_plan_instruction.KurtosisPlanInstruction {
	return add_service.NewAddServices(t.serviceNetwork, t.runtimeValueStore)
}

func (t *addServicesTestCase) GetStarlarkCode() string {
	service1ReadyConditionsScriptPart := getDefaultReadyConditionsScriptPart()
	service2ReadyConditionsScriptPart := getCustomReadyConditionsScripPart(
		TestReadyConditions2RecipePortId,
		TestReadyConditions2RecipeEndpoint,
		TestReadyConditions2RecipeExtract,
		TestReadyConditions2Field,
		TestReadyConditions2Assertion,
		TestReadyConditions2Target,
		TestReadyConditions2Interval,
		TestReadyConditions2Timeout,
	)
	serviceConfig1 := fmt.Sprintf("ServiceConfig(image=%q, ready_conditions=%s)", TestContainerImageName, service1ReadyConditionsScriptPart)
	serviceConfig2 := fmt.Sprintf("ServiceConfig(image=%q, cpu_allocation=%d, memory_allocation=%d, ready_conditions=%s)", TestContainerImageName, TestCpuAllocation, TestMemoryAllocation, service2ReadyConditionsScriptPart)
	return fmt.Sprintf(`%s(%s={%q: %s, %q: %s})`, add_service.AddServicesBuiltinName, add_service.ConfigsArgName, TestServiceName, serviceConfig1, TestServiceName2, serviceConfig2)
}

func (t *addServicesTestCase) GetStarlarkCodeForAssertion() string {
	return ""
}

func (t *addServicesTestCase) Assert(interpretationResult starlark.Value, executionResult *string) {
	resultDict, ok := interpretationResult.(*starlark.Dict)
	require.True(t, ok, "interpretation result should be a dictionary")
	require.Equal(t, resultDict.Len(), 2)
	require.Contains(t, resultDict.Keys(), starlark.String(TestServiceName))
	require.Contains(t, resultDict.Keys(), starlark.String(TestServiceName2))

	require.Contains(t, *executionResult, "Successfully added the following '2' services:")
	require.Contains(t, *executionResult, fmt.Sprintf("Service '%s' added with UUID '%s'", TestServiceName, TestServiceUuid))
	require.Contains(t, *executionResult, fmt.Sprintf("Service '%s' added with UUID '%s'", TestServiceName2, TestServiceUuid2))
}
