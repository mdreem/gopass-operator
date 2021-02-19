package controllers

import (
	"context"
	"github.com/go-logr/logr"
	logr_testing "github.com/go-logr/logr/testing"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"google.golang.org/grpc"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
	"time"
)

type TestRepositoryServiceClient struct {
	Calls map[string][]string
}

func NewTestRepositoryServiceClient() *TestRepositoryServiceClient {
	return &TestRepositoryServiceClient{
		Calls: map[string][]string{
			"InitializeRepository": {},
			"UpdateRepository":     {},
			"UpdateAllPasswords":   {},
			"DeleteSecret":         {},
		},
	}
}

func (r *TestRepositoryServiceClient) InitializeRepository(_ context.Context, _ *gopass_repository.RepositoryInitialization, _ ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) UpdateRepository(_ context.Context, _ *gopass_repository.Repository, _ ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) UpdateAllPasswords(_ context.Context, _ *gopass_repository.Repository, _ ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) DeleteSecret(_ context.Context, _ *gopass_repository.Repository, _ ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func init() {
	err := gopassv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}
}

func TestGopassRepositoryReconciler_handleDeletionOfResource(t *testing.T) {

	type fields struct {
		Client    client.Client
		Log       logr.Logger
		Namespace string
	}
	type args struct {
		ctx           context.Context
		req           controllerruntime.Request
		repository    *gopassv1alpha1.GopassRepository
		serviceClient gopass_repository.RepositoryServiceClient
	}

	const repositoryName = "test-repository"
	const testNameSpace = "test-namespace"
	const finalizerName = "gopass.repository.finalizer"

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantedResult      controllerruntime.Result
		wantedError       error
		wantedSuccess     bool
		wantedServices    int
		wantedDeployments int
	}{
		{
			name: "No deletion scheduled. Add finalizer.",
			fields: fields{
				Client: fake.NewClientBuilder().WithRuntimeObjects(
					&gopassv1alpha1.GopassRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      repositoryName,
							Namespace: testNameSpace,
							DeletionTimestamp: &metav1.Time{
								Time: time.Time{},
							},
						},
					},
				).Build(),
				Log:       logr_testing.NullLogger{},
				Namespace: "",
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: testNameSpace,
						Name:      repositoryName,
					},
				},
				repository: &gopassv1alpha1.GopassRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      repositoryName,
						Namespace: testNameSpace,
						DeletionTimestamp: &metav1.Time{
							Time: time.Time{},
						},
					},
				},
				serviceClient: NewTestRepositoryServiceClient(),
			},
			wantedResult: controllerruntime.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
			wantedError:       nil,
			wantedSuccess:     false,
			wantedServices:    0,
			wantedDeployments: 0,
		},
		{
			name: "Deletion scheduled. Finalizer exists. Deployment and Service do not exist.",
			fields: fields{
				Client: fake.NewClientBuilder().WithRuntimeObjects(
					&gopassv1alpha1.GopassRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      repositoryName,
							Namespace: testNameSpace,
							DeletionTimestamp: &metav1.Time{
								Time: time.Time{},
							},
						},
					},
				).Build(),
				Log:       logr_testing.NullLogger{},
				Namespace: "",
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: testNameSpace,
						Name:      repositoryName,
					},
				},
				repository: &gopassv1alpha1.GopassRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      repositoryName,
						Namespace: testNameSpace,
						DeletionTimestamp: &metav1.Time{
							Time: time.Unix(123, 456),
						},
						Finalizers: []string{finalizerName},
					},
				},
				serviceClient: NewTestRepositoryServiceClient(),
			},
			wantedResult: controllerruntime.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
			wantedError:       nil,
			wantedSuccess:     true,
			wantedServices:    0,
			wantedDeployments: 0,
		},
		{
			name: "Deletion scheduled. Finalizer exists. Deployment and Service do exist.",
			fields: fields{
				Client: fake.NewClientBuilder().WithRuntimeObjects(
					&gopassv1alpha1.GopassRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      repositoryName,
							Namespace: testNameSpace,
							DeletionTimestamp: &metav1.Time{
								Time: time.Time{},
							},
						},
					},
					&appsv1.Deployment{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-deployment",
							Namespace: testNameSpace,
							Labels: map[string]string{
								"gopassRepoName":      repositoryName,
								"gopassRepoNamespace": testNameSpace,
							},
						},
					},
					&corev1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-deployment",
							Namespace: testNameSpace,
							Labels: map[string]string{
								"gopassRepoName":      repositoryName,
								"gopassRepoNamespace": testNameSpace,
							},
						},
					},
				).Build(),
				Log:       logr_testing.NullLogger{},
				Namespace: "",
			},
			args: args{
				ctx: context.Background(),
				req: controllerruntime.Request{
					NamespacedName: types.NamespacedName{
						Namespace: testNameSpace,
						Name:      repositoryName,
					},
				},
				repository: &gopassv1alpha1.GopassRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      repositoryName,
						Namespace: testNameSpace,
						DeletionTimestamp: &metav1.Time{
							Time: time.Unix(123, 456),
						},
						Finalizers: []string{finalizerName},
					},
				},
				serviceClient: NewTestRepositoryServiceClient(),
			},
			wantedResult: controllerruntime.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
			wantedError:       nil,
			wantedSuccess:     true,
			wantedServices:    0,
			wantedDeployments: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GopassRepositoryReconciler{
				Client:    tt.fields.Client,
				Log:       tt.fields.Log,
				Scheme:    scheme.Scheme,
				Namespace: tt.fields.Namespace,
			}

			gotResult, gotError, deletionSuccessful := r.handleDeletionOfResource(tt.args.ctx, tt.args.req, tt.args.repository, tt.args.serviceClient)
			if !reflect.DeepEqual(gotResult, tt.wantedResult) {
				t.Errorf("handleDeletionOfResource() gotResult = %v, wantedResult %v", gotResult, tt.wantedResult)
			}
			if !reflect.DeepEqual(gotError, tt.wantedError) {
				t.Errorf("handleDeletionOfResource() gotError = %v, wantedResult %v", gotError, tt.wantedError)
			}
			if deletionSuccessful != tt.wantedSuccess {
				t.Errorf("handleDeletionOfResource() deletionSuccessful = %v, wantedResult %v", deletionSuccessful, tt.wantedSuccess)
			}

			var services = &corev1.ServiceList{}
			err := tt.fields.Client.List(context.Background(), services)
			if err != nil {
				t.Errorf("unable to fetch list of services: %v", err)
			}
			if len(services.Items) != tt.wantedServices {
				t.Errorf("number of remaining services was '%d', wanted '%v'", len(services.Items), tt.wantedServices)
			}

			var deployments = &appsv1.DeploymentList{}
			err = tt.fields.Client.List(context.Background(), deployments)
			if err != nil {
				t.Errorf("unable to fetch list of deployments: %v", err)
			}
			if len(deployments.Items) != tt.wantedDeployments {
				t.Errorf("number of remaining deployments was '%d', wanted '%v'", len(deployments.Items), tt.wantedDeployments)
			}
		})
	}
}
