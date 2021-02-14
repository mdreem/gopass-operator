package cluster

import (
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"os/exec"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                 string
		wantKubernetesClient bool
		wantErr              bool
	}{
		{
			name:                 "successful creation of KubernetesClient",
			wantKubernetesClient: true,
			wantErr:              false,
		},
		{
			name:                 "unsuccessful creation of KubernetesClient",
			wantKubernetesClient: false,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalGetRestInClusterConfigFunc := getRestInClusterConfigFunc
			if tt.wantErr {
				getRestInClusterConfigFunc = mockedFailingRestInClusterConfig
			} else {
				getRestInClusterConfigFunc = mockedGetRestInClusterConfig
			}
			defer func() {
				getRestInClusterConfigFunc = originalGetRestInClusterConfigFunc
			}()

			kubernetesClient, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			kubernetesClientInitialized := kubernetesClient != KubernetesClient{}
			if kubernetesClientInitialized != tt.wantKubernetesClient {
				t.Errorf("New() kubernetesClient = %v, kubernetesClient: %v", kubernetesClient, tt.wantKubernetesClient)
			}
		})
	}
}

func mockedGetRestInClusterConfig() (*rest.Config, error) {
	return &rest.Config{
		Host: "testHost",
	}, nil
}

func mockedFailingRestInClusterConfig() (*rest.Config, error) {
	return nil, fmt.Errorf("I failed")
}

func TestKubernetesClient_GetRepositoryCredentials(t *testing.T) {
	type fields struct {
		clientset kubernetes.Interface
	}
	type args struct {
		ctx            context.Context
		authentication *gopass_repository.Authentication
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            Secret
		wantErr         bool
		wantedErrorText string
	}{
		{
			name: "successfully fetched secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "someRef",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"someKey": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: nil,
				authentication: &gopass_repository.Authentication{
					Namespace: "testNameSpace",
					Username:  "molly.millions",
					SecretRef: "someRef",
					SecretKey: "someKey",
				},
			},
			want: Secret{
				Name:     "molly.millions",
				Password: "my secret",
			},
			wantErr:         false,
			wantedErrorText: "",
		},
		{
			name: "unable to find secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "someRef",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"someKey": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: nil,
				authentication: &gopass_repository.Authentication{
					Namespace: "testNameSpace",
					Username:  "molly.millions",
					SecretRef: "wrongRef",
					SecretKey: "someKey",
				},
			},
			want:            Secret{},
			wantErr:         true,
			wantedErrorText: "secrets \"wrongRef\" not found",
		},
		{
			name: "unable to find secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "someRef",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"someKey": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: nil,
				authentication: &gopass_repository.Authentication{
					Namespace: "testNameSpace",
					Username:  "molly.millions",
					SecretRef: "someRef",
					SecretKey: "wrongKey",
				},
			},
			want:            Secret{},
			wantErr:         true,
			wantedErrorText: "unable to find key 'wrongKey' in secret 'someRef' in namespace 'testNameSpace'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KubernetesClient{
				clientset: tt.fields.clientset,
			}
			got, err := k.GetRepositoryCredentials(tt.args.ctx, tt.args.authentication)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRepositoryCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && (tt.wantedErrorText != err.Error()) {
				t.Errorf("GetRepositoryCredentials() error = %v, wantedErrorText %v", err, tt.wantedErrorText)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRepositoryCredentials() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKubernetesClient_GetGpgKey(t *testing.T) {
	type fields struct {
		clientset kubernetes.Interface
	}
	type args struct {
		ctx             context.Context
		gpgKeyReference *gopass_repository.GpgKeyReference
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantErr         bool
		wantedErrorText string
		letCommandFail  bool
	}{
		{
			name: "successfully add key",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "gpg-key",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"gpg-key-ref": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				gpgKeyReference: &gopass_repository.GpgKeyReference{
					GpgKeyRef:    "gpg-key",
					GpgKeyRefKey: "gpg-key-ref",
				},
			},
			letCommandFail:  false,
			wantErr:         false,
			wantedErrorText: "",
		},
		{
			name: "unable to find Secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "unkownSecretMap",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"gpg-key-ref": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				gpgKeyReference: &gopass_repository.GpgKeyReference{
					GpgKeyRef:    "gpg-key",
					GpgKeyRefKey: "gpg-key-ref",
				},
			},
			letCommandFail:  false,
			wantErr:         true,
			wantedErrorText: "secrets \"gpg-key\" not found",
		},
		{
			name: "unable to find key in Secret",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "gpg-key",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"wrong-key": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				gpgKeyReference: &gopass_repository.GpgKeyReference{
					GpgKeyRef:    "gpg-key",
					GpgKeyRefKey: "wrong-gpg-key-ref",
				},
			},
			letCommandFail:  false,
			wantErr:         true,
			wantedErrorText: "unable to find key 'wrong-gpg-key-ref' in secret 'gpg-key' in namespace 'testNameSpace'",
		},
		{
			name: "fail to add key due to issue executing gpg",
			fields: fields{
				clientset: fake.NewSimpleClientset(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "gpg-key",
							Namespace: "testNameSpace",
						},
						Data: map[string][]byte{
							"gpg-key-ref": []byte("my secret"),
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				gpgKeyReference: &gopass_repository.GpgKeyReference{
					GpgKeyRef:    "gpg-key",
					GpgKeyRefKey: "gpg-key-ref",
				},
			},
			letCommandFail:  true,
			wantErr:         true,
			wantedErrorText: "exec: \"thisCommandDoesNotExist\": executable file not found in $PATH",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KubernetesClient{
				clientset: tt.fields.clientset,
			}

			originalExecCommandContext := execCommandContext
			if tt.letCommandFail {
				execCommandContext = mockFailedCommandContext
			} else {
				execCommandContext = mockExecCommandContext
			}
			defer func() {
				execCommandContext = originalExecCommandContext
			}()

			err := k.GetGpgKey(tt.args.ctx, "testNameSpace", tt.args.gpgKeyReference)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetGpgKey() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && (tt.wantedErrorText != err.Error()) {
				t.Errorf("GetGpgKey() error = '%v', wantedErrorText '%v'", err, tt.wantedErrorText)
				return
			}
		})
	}
}

func mockExecCommandContext(ctx context.Context, _ string, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "echo", args...)
}

func mockFailedCommandContext(ctx context.Context, _ string, _ ...string) *exec.Cmd {
	return exec.CommandContext(ctx, "thisCommandDoesNotExist")
}
