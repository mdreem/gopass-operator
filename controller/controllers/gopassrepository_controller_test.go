package controllers

import (
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	gopass_repository_grpc "github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"net"
	"time"

	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("GopassRepository", func() {

	const (
		GopassRepositoryName      = "test-gopass-repository"
		GopassRepositoryNamespace = "default"

		UserName      = "Henry.Dorsett.Case"
		RepositoryUrl = "someUrl"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating GopassRepository", func() {
		It("Should initialize a new gopass repository", func() {
			ctx := context.Background()

			originalCreateRepositoryServiceClientFunc := createRepositoryServiceClient
			originalGetRelevantDeploymentFunc := getRelevantDeployment

			createRepositoryServiceClientFunc = createRepositoryServiceClientForTesting
			getRelevantDeploymentFunc = getRelevantDeploymentMock

			defer func() {
				createRepositoryServiceClientFunc = originalCreateRepositoryServiceClientFunc
				getRelevantDeploymentFunc = originalGetRelevantDeploymentFunc
			}()

			testRepositoryServiceServer := InitializeTestRepositoryServer()
			go func() {
				var err error
				testRepositoryServiceServer, err = initializeTestServer(testRepositoryServiceServer)
				Expect(err).ToNot(HaveOccurred())
			}()

			gopassRepository := &gopassv1alpha1.GopassRepository{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      GopassRepositoryName,
					Namespace: GopassRepositoryNamespace,
				},
				Spec: gopassv1alpha1.GopassRepositorySpec{
					RepositoryURL:   RepositoryUrl,
					RefreshInterval: "",
					UserName:        UserName,
					SecretKeyRef:    gopassv1alpha1.SecretKeyRefSpec{},
				},
				Status: gopassv1alpha1.GopassRepositoryStatus{},
			}

			Expect(k8sClient.Create(ctx, gopassRepository)).Should(Succeed())

			gopassRepositoryLookupKey := types.NamespacedName{Name: GopassRepositoryName, Namespace: GopassRepositoryNamespace}
			createdGopassRepository := &gopassv1alpha1.GopassRepository{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, gopassRepositoryLookupKey, createdGopassRepository)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue(), "wait for GopassRepository to be created")

			Expect(createdGopassRepository.Spec.UserName).Should(Equal(UserName))

			Eventually(func() bool {
				return len((*testRepositoryServiceServer).Calls["InitializeRepository"]) > 0
			}, timeout, interval).Should(BeTrue(), "wait for repository server to be called")

			initializeRepositoryCalls := (*testRepositoryServiceServer).Calls["InitializeRepository"]
			Expect(len(initializeRepositoryCalls)).Should(Equal(1), "correct number of calls to InitializeRepository")
			Expect(initializeRepositoryCalls[0]).Should(Equal(RepositoryUrl), "calls to InitializeRepository should contain correct URL")

			updateRepositoryCalls := (*testRepositoryServiceServer).Calls["UpdateRepository"]
			Expect(len(updateRepositoryCalls)).Should(Equal(1), "correct number of calls to UpdateRepository")
			Expect(updateRepositoryCalls[0]).Should(Equal(RepositoryUrl), "calls to UpdateRepository should contain correct URL")

			updateAllPasswordsCalls := (*testRepositoryServiceServer).Calls["UpdateAllPasswords"]
			Expect(len(updateAllPasswordsCalls)).Should(Equal(1), "correct number of calls to UpdateAllPasswords")
			Expect(updateAllPasswordsCalls[0]).Should(Equal(RepositoryUrl), "calls to UpdateAllPasswords should contain correct URL")
		})
	})
})

func initializeTestServer(server *TestRepositoryServer) (*TestRepositoryServer, error) {
	grpcServer := grpc.NewServer()
	gopass_repository_grpc.RegisterRepositoryServiceServer(grpcServer, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 12345))
	if err != nil {
		return nil, err
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func createRepositoryServiceClientForTesting(targetUrl string) (gopass_repository.RepositoryServiceClient, *grpc.ClientConn, error) {
	fmt.Printf("createRepositoryServiceClientForTesting called with targetUrl='%s'", targetUrl)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:12345", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return gopass_repository.NewRepositoryServiceClient(conn), conn, nil
}

func getRelevantDeploymentMock(deployments *[]appsv1.Deployment) (*appsv1.Deployment, error) {
	return &appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			AvailableReplicas: 1,
		},
	}, nil
}
