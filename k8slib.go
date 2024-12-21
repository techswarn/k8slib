package k8slib

import (
	"log"
	"time"
	"github.com/techswarn/k8slib/utils"
	"k8s.io/client-go/kubernetes"
	// appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "github.com/Azure/go-autorest/autorest/to"
	_ "errors"
	"context"
)

type Deploy struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
	Namespace string `json:namespace`
	Status bool `json:status`
	CreatedAt time.Time `json:"createdat"`
}

var cs *kubernetes.Clientset
func Connect() *kubernetes.Clientset {
    cs, _ = utils.GetKubehandle()
	return cs
}

func (d *Deploy) CreateNamespace(ctx context.Context, clientSet *kubernetes.Clientset) *corev1.Namespace {
	log.Printf("Creating namespace %q.\n\n", d.Name)
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.Name,
		},
	}
	ns, err := clientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	utils.PanicIfError(err)
	return ns
}