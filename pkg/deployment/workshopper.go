package deployment

import (
	"strings"

	openshiftv1alpha1 "github.com/redhat/openshift-workshop-operator/pkg/apis/openshift/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewWorkshopperDeployment(cr *openshiftv1alpha1.Workshop, name string, namespace string,
	projectName string, infraProjectName string, username string, appsHostnameSuffix string,
	openshiftConsoleURL string, openshiftAPIURL string) *appsv1.Deployment {
	workshopperImage := "quay.io/osevg/workshopper:latest"
	labels := GetLabels(cr, name)

	guidePath := strings.TrimPrefix(cr.Spec.Source.GitURL, "https://github.com/")
	guideBranch := cr.Spec.Source.GitBranch

	workshopperEnv := []corev1.EnvVar{
		{
			Name:  "CONTENT_URL_PREFIX",
			Value: "https://raw.githubusercontent.com/" + guidePath + "/" + guideBranch + "/guide",
		},
		{
			Name:  "WORKSHOPS_URLS",
			Value: "https://raw.githubusercontent.com/" + guidePath + "/" + guideBranch + "/guide/_workshop.yml",
		},
		{
			Name:  "PROJECT",
			Value: projectName,
		},
		{
			Name:  "INFRA_PROJECT",
			Value: infraProjectName,
		},
		{
			Name:  "OPENSHIFT_CONSOLE_URL",
			Value: openshiftConsoleURL,
		},
		{
			Name:  "OPENSHIFT_API_URL",
			Value: openshiftAPIURL,
		},
		{
			Name:  "OPENSHIFT_USER",
			Value: username,
		},
		{
			Name:  "OPENSHIFT_PASSWORD",
			Value: cr.Spec.User.Password,
		},
		{
			Name:  "APPS_HOSTNAME_SUFFIX",
			Value: appsHostnameSuffix,
		},
		{
			Name:  "LABS_GIT_REPO",
			Value: "https://github.com/" + guidePath + ".git#" + guideBranch,
		},
		{
			Name:  "LABS_DOWNLOAD_URL",
			Value: "https://github.com/" + guidePath + "/archive/" + guideBranch + ".zip",
		},
		{
			Name:  "WEB_NODEJS_GIT_REPO",
			Value: "https://github.com/" + guidePath + "/tree/" + guideBranch + "/labs/web-nodejs",
		},
		{
			Name:  "CATALOG_GO_GIT_REPO",
			Value: "https://github.com/" + guidePath + "/tree/" + guideBranch + "/labs/catalog-go",
		},
		{
			Name:  "CHE_URL",
			Value: "http://che-eclipse-che." + appsHostnameSuffix,
		},
		{
			Name:  "NEXUS_URL",
			Value: "http://nexus-opentlc-shared" + appsHostnameSuffix,
		},
		{
			Name:  "KIALI_URL",
			Value: "https://kiali-istio-system." + appsHostnameSuffix,
		},
		{
			Name:  "JAEGER_URL",
			Value: "https://tracing-istio-system." + appsHostnameSuffix,
		},
		{
			Name:  "KIBANA_URL",
			Value: "https://kibana-openshift-logging." + appsHostnameSuffix,
		},
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            name,
							Image:           workshopperImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          name,
									ContainerPort: 8080,
									Protocol:      "TCP",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: int32(8080),
										},
										Scheme: corev1.URISchemeHTTP,
									},
								},
								InitialDelaySeconds: 25,
								FailureThreshold:    10,
								TimeoutSeconds:      5,
							},
							Env: workshopperEnv,
						},
					},
				},
			},
		},
	}
}
