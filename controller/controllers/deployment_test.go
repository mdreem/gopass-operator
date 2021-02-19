package controllers

import (
	"context"
	"github.com/go-logr/logr"
	logr_testing "github.com/go-logr/logr/testing"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestGopassRepositoryReconciler_createRepositoryServer(t *testing.T) {
	type fields struct {
		Client    client.Client
		Log       logr.Logger
		Scheme    *runtime.Scheme
		Namespace string
	}
	type args struct {
		ctx            context.Context
		namespacedName types.NamespacedName
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		want              bool
		wantErr           bool
		wantedServices    int
		wantedDeployments int
	}{
		{
			name: "test-repository",
			fields: fields{
				Client:    fake.NewClientBuilder().WithRuntimeObjects().Build(),
				Log:       logr_testing.NullLogger{},
				Namespace: "test-namespace",
			},
			args: args{
				ctx: context.Background(),
				namespacedName: types.NamespacedName{
					Namespace: "repoNamespace",
					Name:      "repoName",
				},
			},
			want:              false,
			wantErr:           false,
			wantedServices:    1,
			wantedDeployments: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GopassRepositoryReconciler{
				Client:    tt.fields.Client,
				Log:       tt.fields.Log,
				Scheme:    tt.fields.Scheme,
				Namespace: tt.fields.Namespace,
			}
			got, err := r.createRepositoryServer(tt.args.ctx, tt.args.namespacedName)
			if (err != nil) != tt.wantErr {
				t.Errorf("createRepositoryServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createRepositoryServer() got = %v, want %v", got, tt.want)
			}

			var services = &corev1.ServiceList{}
			err = tt.fields.Client.List(context.Background(), services)
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

			for _, deployment := range deployments.Items {
				repoName, ok := deployment.Labels["gopassRepoName"]
				if ok != true {
					t.Errorf("gopassRepoName not found in labels of deployment")
				}
				if repoName != tt.args.namespacedName.Name {
					t.Errorf("repoName of deployment was '%s', wanted '%s'", repoName, tt.args.namespacedName.Name)
				}

				repoNamespace, ok := deployment.Labels["gopassRepoNamespace"]
				if ok != true {
					t.Errorf("gopassRepoName not found in labels of deployment")
				}
				if repoNamespace != tt.args.namespacedName.Namespace {
					t.Errorf("repoNamespace of deployment was '%s', wanted '%s'", repoNamespace, tt.args.namespacedName.Namespace)
				}
			}

			for _, service := range services.Items {
				repoName, ok := service.Labels["gopassRepoName"]
				if ok != true {
					t.Errorf("gopassRepoName not found in labels of service")
				}
				if repoName != tt.args.namespacedName.Name {
					t.Errorf("repoName of service was '%s', wanted '%s'", repoName, tt.args.namespacedName.Name)
				}

				repoNamespace, ok := service.Labels["gopassRepoNamespace"]
				if ok != true {
					t.Errorf("gopassRepoName not found in labels of service")
				}
				if repoNamespace != tt.args.namespacedName.Namespace {
					t.Errorf("repoNamespace of service was '%s', wanted '%s'", repoNamespace, tt.args.namespacedName.Namespace)
				}
			}
		})
	}
}
