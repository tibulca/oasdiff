package diff_test

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func getReqPropFile(file string) string {
	return fmt.Sprintf("../data/required-properties/%s", file)
}

func TestBreaking_NewRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// new required property in request header breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_NewNonRequiredProperty(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// new optional property in request header doesn't break client
	require.Empty(t, d)
}

func TestBreaking_PropertyRequiredEnabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in request header to required breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_PropertyRequiredDisabled(t *testing.T) {
	s1 := l(t, 1)
	s2 := l(t, 1)

	sr := openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Type:        "string",
			Description: "Unique ID of the course",
		},
	}

	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s1.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{"courseId"}

	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Properties["courseId"] = &sr
	s2.Paths[installCommandPath].Get.Parameters.GetByInAndName(openapi3.ParameterInHeader, "network-policies").Schema.Value.Required = []string{}

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in request header to optional doesn't break client
	require.Empty(t, d)
}

func TestBreaking_RespBodyRequiredPropertyDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-required-properties-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-required-properties-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in response body to optional breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_RespBodyRequiredPropertyEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-required-properties-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-required-properties-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in response body to required doesn't break client
	require.Empty(t, d)
}

func TestBreaking_ReqBodyRequiredPropertyDisabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-required-properties-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-required-properties-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in request body to optional doesn't break client
	require.Empty(t, d)
}

func TestBreaking_ReqBodyRequiredPropertyEnabled(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-required-properties-revision.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-required-properties-base.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// changing an existing property in request body to required breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_ReqBodyNewRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-new-required-properties-base.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-new-required-properties-revision.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// adding a new required property in request body breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_ReqBodyDeleteRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("request-new-required-properties-revision.yaml"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("request-new-required-properties-base.yaml"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// deleting a required property in request doesn't break client
	require.Empty(t, d)
}

func TestBreaking_RespBodyNewRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-new-required-properties-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-new-required-properties-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// adding a new required property in response body doesn't break client
	require.Empty(t, d)
}

func TestBreaking_RespBodyDeleteRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-new-required-properties-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-new-required-properties-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// deleting a required property in response body breaks client
	require.NotEmpty(t, d)
}

func TestBreaking_RespBodyNewAllOfRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-allof-required-properties-base.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-allof-required-properties-revision.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// adding a new required property in response body doesn't break client
	require.Empty(t, d)
}

func TestBreaking_RespBodyDeleteAllOfRequiredProperty(t *testing.T) {
	loader := openapi3.NewLoader()

	s1, err := loader.LoadFromFile(getReqPropFile("response-allof-required-properties-revision.json"))
	require.NoError(t, err)

	s2, err := loader.LoadFromFile(getReqPropFile("response-allof-required-properties-base.json"))
	require.NoError(t, err)

	d, err := diff.Get(&diff.Config{
		BreakingOnly: true,
	}, s1, s2)
	require.NoError(t, err)

	// deleting a required property in response body breaks client
	require.NotEmpty(t, d)
}
