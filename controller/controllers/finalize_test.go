package controllers

import (
	"context"
	"github.com/go-logr/logr"
	logr_testing "github.com/go-logr/logr/testing"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"google.golang.org/grpc"
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

func (r *TestRepositoryServiceClient) InitializeRepository(ctx context.Context, in *gopass_repository.RepositoryInitialization, opts ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) UpdateRepository(ctx context.Context, in *gopass_repository.Repository, opts ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) UpdateAllPasswords(ctx context.Context, in *gopass_repository.Repository, opts ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServiceClient) DeleteSecret(ctx context.Context, in *gopass_repository.Repository, opts ...grpc.CallOption) (*gopass_repository.RepositoryResponse, error) {
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func init() {
	gopassv1alpha1.AddToScheme(scheme.Scheme)
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
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantedResult  controllerruntime.Result
		wantedError   error
		wantedSuccess bool
	}{
		{
			name: "No deletion scheduled. Add finalizer.",
			fields: fields{
				Client: fake.NewClientBuilder().WithRuntimeObjects(
					&gopassv1alpha1.GopassRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-repository",
							Namespace: "test-namespace",
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
						Namespace: "test-namespace",
						Name:      "test-repository",
					},
				},
				repository: &gopassv1alpha1.GopassRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-repository",
						Namespace: "test-namespace",
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
			wantedError:   nil,
			wantedSuccess: false,
		},
		{
			name: "Deletion scheduled. Finalizer exists. Deployment and Service do not exist.",
			fields: fields{
				Client: fake.NewClientBuilder().WithRuntimeObjects(
					&gopassv1alpha1.GopassRepository{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-repository",
							Namespace: "test-namespace",
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
						Namespace: "test-namespace",
						Name:      "test-repository",
					},
				},
				repository: &gopassv1alpha1.GopassRepository{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-repository",
						Namespace: "test-namespace",
						DeletionTimestamp: &metav1.Time{
							Time: time.Unix(123, 456),
						},
						Finalizers: []string{"gopass.repository.finalizer"},
					},
				},
				serviceClient: NewTestRepositoryServiceClient(),
			},
			wantedResult: controllerruntime.Result{
				Requeue:      false,
				RequeueAfter: 0,
			},
			wantedError:   nil,
			wantedSuccess: true,
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
		})
	}
}
