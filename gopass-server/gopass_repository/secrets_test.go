package gopass_repository

import (
	"context"
	"github.com/gopasspw/gopass/pkg/gopass/apimock"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/fake"
	clientgotesting "k8s.io/client-go/testing"
	"reflect"
	"testing"
)

func TestRepositoryServer_updateAllPasswords(t *testing.T) {
	type fields struct {
		Repositories     map[string]*gopassRepo
		Client           cluster.Client
		KubernetesClient *fake.Clientset
	}
	type args struct {
		ctx        context.Context
		repository *gopass_repository.Repository
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		wantedActions []clientgotesting.Action
	}{
		{
			name: "Secret map does not exist yet. Creating it.",
			fields: fields{
				Repositories: map[string]*gopassRepo{
					"testUrl": {
						store:      apimock.New(),
						directory:  "",
						repository: nil,
					},
				},
				Client:           &cluster.KubernetesTestClient{},
				KubernetesClient: fake.NewSimpleClientset(),
			},
			args: args{
				ctx: nil,
				repository: &gopass_repository.Repository{
					RepositoryURL:  "testUrl",
					Authentication: &gopass_repository.Authentication{},
					SecretName: &gopass_repository.NamespacedName{
						Namespace: "testNamespace",
						Name:      "someSecret",
					},
				},
			},
			wantErr: false,
			wantedActions: []clientgotesting.Action{
				clientgotesting.NewCreateAction(schema.GroupVersionResource{Version: "v1", Resource: "secrets"}, "testNamespace", &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "someSecret",
						Namespace: "testNamespace",
					},
					StringData: make(map[string]string),
				}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := &RepositoryServer{
				Repositories:     tt.fields.Repositories,
				Client:           tt.fields.Client,
				KubernetesClient: tt.fields.KubernetesClient,
			}
			if err := r.updateAllPasswords(tt.args.ctx, tt.args.repository); (err != nil) != tt.wantErr {
				t.Errorf("updateAllPasswords() error = %v, wantErr %v", err, tt.wantErr)
			}

			actions := tt.fields.KubernetesClient.Actions()
			for _, wantedAction := range tt.wantedActions {
				if !containsAction(actions, wantedAction) {
					t.Errorf("List of actions does not contain: %v\n contained: %v", wantedAction, actions)
				}
			}
		})
	}
}

func containsAction(actions []clientgotesting.Action, wantedAction clientgotesting.Action) bool {
	for _, foundAction := range actions {
		if reflect.DeepEqual(foundAction, wantedAction) {
			return true
		}
	}
	return false
}
