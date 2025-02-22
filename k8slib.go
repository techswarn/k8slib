package k8slib

import (
	"log"
    "time"
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

type Replica struct {
	Name string
	Status bool
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
			Name: i.Namespace,
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
			Replicas: to.Int32Ptr(1),
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

func (i *Instance) WaitForReadyReplicas(ctx context.Context, clientSet *kubernetes.Clientset, deployment *appv1.Deployment) *Replica {
	log.Printf("Waiting for ready replicas in deployment %q\n", deployment.Name)
	for {
		expectedReplicas := *deployment.Spec.Replicas
		readyReplicas := getReadyReplicasForDeployment(ctx, clientSet, deployment)
		if readyReplicas == expectedReplicas {
			log.Printf("replicas are ready!\n\n")
			return &Replica{
				Name: deployment.Name,
				Status: true,
			}
			break
		}

		log.Printf("replicas are not ready yet. %d/%d\n", readyReplicas, expectedReplicas)
		time.Sleep(1 * time.Second)
	}

	return &Replica{
		Name: "",
		Status: false,
	}
}
//*corev1.Namespace
func DeleteNamespace(ctx context.Context, clientSet *kubernetes.Clientset, ns string) {
	log.Printf("\n\nDeleting namespace %q.\n", ns)
	result := clientSet.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})
	log.Printf("\n\nDeleting namespace %#v.\n", result)
}

func getReadyReplicasForDeployment(ctx context.Context, clientSet *kubernetes.Clientset, deployment *appv1.Deployment) int32 {
	dep, err := clientSet.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
	utils.PanicIfError(err)
	return dep.Status.ReadyReplicas
}

func (i *Instance) GetDeployment(ctx context.Context, clientSet *kubernetes.Clientset) int32 {
	var dep int32
	for {
		deployment, err := clientSet.AppsV1().Deployments(i.Namespace).Get(ctx, i.Name, metav1.GetOptions{})
		log.Println(err)
		dep := deployment.Status.ReadyReplicas
		log.Printf("number %d", dep)
		if dep == 1 {
			return dep
			break
		}
		time.Sleep(1 * time.Second)
		
	}
	// log.Printf("number %d", dep)
	// return dep
	return dep
}
