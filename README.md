# gopass-operator

:warning: this is only an experiment for now :warning:

This operator handles deploying GoPass-repositories as a `Secret` into the cluster.

The GoPass repository to use is given in `repositoryUrl`, the credentials to access it are given in the
section `secretKeyRef`. Currently only `username`/`password` are supported.

The GPG key needed to access the secrets in the repository is references in the section `gpgKeyRef`, where `name`
references the `Secret` and `key` references the key inside this `Secret` that contains the GPG key. The `Secret` needs
to reside in the same namespace as the operator itself.

```yaml
apiVersion: gopass.gopass.operator/v1alpha1
kind: GopassRepository
metadata:
  name: gopassrepository-sample
spec:
  refreshInterval: 30s
  repositoryUrl: "someRepositoryUrl"
  userName: "git"
  secretKeyRef:
    name: "gopass-test-secret"
    key: "gopass-test-secret-key"
  gpgKeyRef:
    name: "gpg-key"
    key: "gpg-key"
```

The created `Secret` will consist of all accessible entries in the GoPass repository. All characters that are not
alphanumeric will be replaced with `-` to become compatible with names in kubernetes resources. If two entries result in
the same key, one will be overridden.

## How does it work

When a new `GopassRepository` is created, it spins up a new repository server in the same namespace as the controller
itself. It is mainly used to separate different GPG keys from each other. Every repository-server will only hold one GPG
key.
