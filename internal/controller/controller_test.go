package controller

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Deployment Controller", func() {
	const (
		DeploymentName      = "test-deployment"
		DeploymentNamespace = "default"
		DeploymentImage     = "nginx"
		timeout             = time.Second * 20
		interval            = time.Millisecond * 20
	)
	Context("When reconciling a Deployment", func() {
		It("Should first create a test deployment", func() {
			By("Creating the deployment resource")
			selector := make(map[string]string)
			selector["app"] = DeploymentImage
			mdata := metav1.ObjectMeta{
				Name:      DeploymentName,
				Namespace: DeploymentNamespace,
				Labels:    selector,
			}
			depSelector := metav1.LabelSelector{
				MatchLabels: selector,
			}
			var containers []corev1.Container
			container := corev1.Container{
				Name:  DeploymentImage,
				Image: DeploymentImage,
			}
			podTemp := corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: append(containers, container),
				},
				ObjectMeta: mdata,
			}
			replicas := int32(1)
			testDep := &appsv1.Deployment{
				ObjectMeta: mdata,
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &depSelector,
					Template: podTemp,
				},
			}
			Expect(k8sClient.Create(ctx, testDep)).To(Succeed())
		})
		It("Should label the deployment with the found ReplicaSet", func() {
			By("Checking for the corresponding label on the test deployment")
			dep := &appsv1.Deployment{}
			Eventually(func(g Gomega) {
				k8sClient.Get(ctx, client.ObjectKey{Name: DeploymentName, Namespace: DeploymentNamespace}, dep)
				g.Expect(dep.ObjectMeta.Labels["owned-replica-set"]).ToNot(BeEmpty())
			}, timeout, interval).Should(Succeed())

			By("Checking if the replica set exists")
			rs := &appsv1.ReplicaSet{}
			rsKey := client.ObjectKey{
				Name:      dep.Labels["owned-replica-set"],
				Namespace: DeploymentNamespace,
			}
			Expect(k8sClient.Get(ctx, rsKey, rs)).To(Succeed())

			By("Checking if the right replica set is labeled on the deployment")
			found := false
			for _, owner := range rs.OwnerReferences {
				if owner.Name == DeploymentName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})
		It("Should clean up the test environment", func() {
			By("By deleting the test deployment")
			dep := &appsv1.Deployment{}
			k8sClient.Get(ctx, client.ObjectKey{Name: DeploymentName, Namespace: DeploymentNamespace}, dep)
			Expect(k8sClient.Delete(ctx, dep)).To(Succeed())
		})
	})
})
