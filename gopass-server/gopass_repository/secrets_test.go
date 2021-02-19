package gopass_repository

import (
	"context"
	"github.com/gopasspw/gopass/pkg/gopass/apimock"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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
		passwords     map[string]string
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
		{
			name: "Update existing Secret map.",
			fields: fields{
				Repositories: map[string]*gopassRepo{
					"testUrl": {
						store:      apimock.New(),
						directory:  "",
						repository: nil,
					},
				},
				Client: &cluster.KubernetesTestClient{},
				KubernetesClient: fake.NewSimpleClientset(&corev1.Secret{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "testNamespace",
						Name:      "someSecret",
					},
				}),
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
				clientgotesting.NewUpdateAction(schema.GroupVersionResource{Version: "v1", Resource: "secrets"}, "testNamespace", &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "someSecret",
						Namespace: "testNamespace",
					},
					StringData: map[string]string{"secretKey": "secretSecret"},
				}),
			},
			passwords: map[string]string{"secretKey": "secretSecret"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for key, password := range tt.passwords {
				sec := apimock.Secret{
					Buf: []byte(password),
				}

				err := tt.fields.Repositories["testUrl"].store.Set(context.Background(), key, &sec)
				if err != nil {
					t.Errorf("unable to set key in store: %v", err)
				}
			}

			r := &RepositoryServer{
				Repositories:     tt.fields.Repositories,
				Client:           tt.fields.Client,
				KubernetesClient: tt.fields.KubernetesClient,
			}

			response, err := r.UpdateAllPasswords(tt.args.ctx, tt.args.repository)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateAllPasswords() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && !response.Successful {
				t.Error("response not successful")
			}

			if err != nil && response.ErrorMessage != "" {
				t.Errorf("received error message '%s', expected empty error message", response.ErrorMessage)
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

func TestRepositoryServer_deleteSecretMap(t *testing.T) {
	type fields struct {
		Repositories     map[string]*gopassRepo
		Client           cluster.Client
		KubernetesClient *fake.Clientset
	}
	type args struct {
		ctx            context.Context
		namespacedName types.NamespacedName
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          bool
		wantErr       bool
		wantedActions []clientgotesting.Action
	}{
		{
			name: "Delete existing Secret",
			fields: fields{
				KubernetesClient: fake.NewSimpleClientset(&corev1.Secret{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "someNamespace",
						Name:      "someName",
					},
					StringData: map[string]string{"aSecretSecret": "secretPassword"},
				}),
			},
			args: args{
				ctx: nil,
				namespacedName: types.NamespacedName{
					Namespace: "someNamespace",
					Name:      "someName",
				},
			},
			want:    true,
			wantErr: false,
			wantedActions: []clientgotesting.Action{
				clientgotesting.NewDeleteAction(schema.GroupVersionResource{Version: "v1", Resource: "secrets"}, "someNamespace", "someName"),
			},
		},
		{
			name: "Delete Secret that does not exist",
			fields: fields{
				KubernetesClient: fake.NewSimpleClientset(&corev1.Secret{
					TypeMeta: metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "someNamespace",
						Name:      "unknownSecret",
					},
					StringData: map[string]string{"aSecretSecret": "secretPassword"},
				}),
			},
			args: args{
				ctx: nil,
				namespacedName: types.NamespacedName{
					Namespace: "someNamespace",
					Name:      "someName",
				},
			},
			want:          true,
			wantErr:       false,
			wantedActions: []clientgotesting.Action{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepositoryServer{
				Repositories:     tt.fields.Repositories,
				Client:           tt.fields.Client,
				KubernetesClient: tt.fields.KubernetesClient,
			}

			successful, err := r.deleteSecretMap(tt.args.ctx, tt.args.namespacedName)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteSecretMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if successful != tt.want {
				t.Errorf("deleteSecretMap() got = %v, want %v", successful, tt.want)
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
