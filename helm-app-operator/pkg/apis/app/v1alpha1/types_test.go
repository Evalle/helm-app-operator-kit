package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/helm/pkg/proto/hapi/release"
)

const (
	testNamespaceName = "helm-test"
	testReleaseName   = "helm-test-dory"
)

var now = metav1.Now()

func TestSetPhase(t *testing.T) {
	newStatus, err := newTestStatus().SetPhase(PhaseApplying, ReasonCustomResourceUpdated, "working on it").ToMap()
	assert.NoError(t, err)

	assert.Equal(t, string(PhaseApplying), newStatus["phase"])
	assert.Equal(t, string(ReasonCustomResourceUpdated), newStatus["reason"])
	assert.Equal(t, "working on it", newStatus["message"])
	assert.NotEqual(t, metav1.Now(), newStatus["lastUpdateTime"])
	assert.NotEqual(t, metav1.Now(), newStatus["lastTransitionTime"])
}

func TestStatusForEmpty(t *testing.T) {
	status := StatusFor(newTestResource())

	assert.Equal(t, &HelmAppStatus{}, status)
}

func TestStatusForFilled(t *testing.T) {
	expectedResource := newTestResource()
	expectedResource.Object["status"] = newTestStatusRaw()
	status := StatusFor(expectedResource)

	assert.EqualValues(t, newTestStatus().Phase, status.Phase)
	assert.EqualValues(t, newTestStatus().Reason, status.Reason)
	assert.EqualValues(t, newTestStatus().Message, status.Message)
}

func TestSetRelease(t *testing.T) {
	releaseName := "TestRelease"
	release := release.Release{Name: releaseName}
	newStatus := newTestStatus().SetRelease(&release)
	assert.EqualValues(t, newStatus.Release.Name, releaseName)
}

func newTestResource() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Character",
			"apiVersion": "stable.nicolerenee.io",
			"metadata": map[string]interface{}{
				"name":      "dory",
				"namespace": testNamespaceName,
			},
			"spec": map[string]interface{}{
				"Name": "Dory",
				"From": "Finding Nemo",
				"By":   "Disney",
			},
		},
	}
}

func newTestStatus() *HelmAppStatus {
	return &HelmAppStatus{
		Phase:              PhaseApplied,
		Reason:             ReasonApplySuccessful,
		Message:            "some message",
		LastUpdateTime:     now,
		LastTransitionTime: now,
	}
}

func newTestStatusRaw() map[string]interface{} {
	return map[string]interface{}{
		"phase":              PhaseApplied,
		"reason":             ReasonApplySuccessful,
		"message":            "some message",
		"lastUpdateTime":     now.UTC(),
		"lastTransitionTime": now.UTC(),
	}
}
