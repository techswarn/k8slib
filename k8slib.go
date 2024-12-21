package k8slib

import (
	"log"
	_ "time"
	"github.com/techswarn/k8slib/utils"
	"k8s.io/client-go/kubernetes"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/Azure/go-autorest/autorest/to"
	_ "errors"
	"context"
)

type Instance struct {
	Name string `json:"name"`
	Image string `json:"image"`
	Namespace string `json:namespace`
}

var cs *kubernetes.Clientset
func Connect() *kubernetes.Clientset {
    cs, _ = utils.GetKubehandle()
	return cs
}

func (i *Instance) CreateNamespace(ctx context.Context, clientSet *kubernetes.Clientset) *corev1.Namespace {
	log.Printf("Creating namespace %q.\n\n", i.Name)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: i.Name,
		},
	}
	ns, err := clientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	utils.PanicIfError(err)
	return ns
}

// func Deploy(ctx context.Context, clientSet *kubernetes.Clientset, ns *corev1.Namespace, name string, image string) {
// 	deployment := createNginxDeployment(ctx, clientSet, ns, name)
// 	waitForReadyReplicas(ctx, clientSet, deployment)
// 	createNginxService(ctx, clientSet, ns, name)
// 	createNginxIngress(ctx, clientSet, ns, name)
// }

func (i *Instance) CreateDeployment(ctx context.Context, clientSet *kubernetes.Clientset, ns *corev1.Namespace) *appv1.Deployment {
	var (
		matchLabel = map[string]string{"app": "nginx"}
		objMeta    = metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: ns.Name,
			Labels:    matchLabel,
		}
	)

	deployment := &appv1.Deployment{
		ObjectMeta: objMeta,
		Spec: appv1.DeploymentSpec{
			Replicas: to.Int32Ptr(2),
			Selector: &metav1.LabelSelector{MatchLabels: matchLabel},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchLabel,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  i.Name,
							Image: i.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	deployment, err := clientSet.AppsV1().Deployments(ns.Name).Create(ctx, deployment, metav1.CreateOptions{})
	utils.PanicIfError(err)
	return deployment
}
