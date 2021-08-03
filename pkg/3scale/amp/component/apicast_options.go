package component

import (
	"crypto/md5"
	"fmt"

	"github.com/go-playground/validator/v10"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/3scale/3scale-operator/pkg/helper"
)

type CustomPolicy struct {
	Name      string
	Version   string
	SecretRef v1.LocalObjectReference
}

func (c CustomPolicy) VolumeName() string {
	return fmt.Sprintf("policy-%s-%s", helper.DNS1123Name(c.Version), helper.DNS1123Name(c.Name))
}

func (c CustomPolicy) AnnotationKey() string {
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// prefix/name: value
	// The name segment is required and must be 63 characters or less
	// Currently: len(CustomPoliciesAnnotationKeyPrefix) + 32 (from the hash) = 54
	return fmt.Sprintf("%s-%x", CustomPoliciesAnnotationKeyPrefix, md5.Sum([]byte(c.VolumeName())))
}

func (c CustomPolicy) AnnotationValue() string {
	return c.VolumeName()
}

type ApicastOptions struct {
	ManagementAPI                  string `validate:"required"`
	OpenSSLVerify                  string `validate:"required"`
	ResponseCodes                  string `validate:"required"`
	ImageTag                       string `validate:"required"`
	ExtendedMetrics                bool
	ProductionResourceRequirements v1.ResourceRequirements `validate:"-"`
	StagingResourceRequirements    v1.ResourceRequirements `validate:"-"`
	ProductionReplicas             int32
	StagingReplicas                int32
	CommonLabels                   map[string]string `validate:"required"`
	CommonStagingLabels            map[string]string `validate:"required"`
	CommonProductionLabels         map[string]string `validate:"required"`
	StagingPodTemplateLabels       map[string]string `validate:"required"`
	ProductionPodTemplateLabels    map[string]string `validate:"required"`
	ProductionAffinity             *v1.Affinity      `validate:"-"`
	ProductionTolerations          []v1.Toleration   `validate:"-"`
	StagingAffinity                *v1.Affinity      `validate:"-"`
	StagingTolerations             []v1.Toleration   `validate:"-"`
	ProductionWorkers              *int32            `validate:"-"`

	// Used for monitoring objects
	// Those objects are namespaced. However, objects includes labels, rules and expressions
	// that need namespace filtering because they are "global" once imported
	// to the prometheus or grafana services.
	Namespace string `validate:"required"`

	ProductionLogLevel *string `validate:"-"`
	StagingLogLevel    *string `validate:"-"`

	ProductionCustomPolicies []CustomPolicy `validate:"-"`
	StagingCustomPolicies    []CustomPolicy `validate:"-"`
}

func NewApicastOptions() *ApicastOptions {
	return &ApicastOptions{}
}

func (a *ApicastOptions) Validate() error {
	validate := validator.New()
	return validate.Struct(a)
}

func DefaultProductionResourceRequirements() v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1000m"),
			v1.ResourceMemory: resource.MustParse("128Mi"),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("500m"),
			v1.ResourceMemory: resource.MustParse("64Mi"),
		},
	}
}

func DefaultStagingResourceRequirements() v1.ResourceRequirements {
	return v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("100m"),
			v1.ResourceMemory: resource.MustParse("128Mi"),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("50m"),
			v1.ResourceMemory: resource.MustParse("64Mi"),
		},
	}
}
