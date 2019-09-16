# aws-sm-kubernetes-helper

The `aws-sm-kubernetes-helper` is a small application/container that authenticates with [AWS Secrets Manager][aws-sm], pulls a secret, and places it in a well-known, configurable location. It is most commonly used as an init container to supply credentials to applications or services.

[aws-sm]: https://aws.amazon.com/secrets-manager/


## Configuration

- `AWS_REGION` - the AWS region the secret is stored in (ex. us-west-2)

- `AWS_ACCESS_KEY_ID` - the AWS access key ID of the IAM user that has secretsmanager:GetSecretValue access to the saved secret

- `AWS_SECRET_ACCESS_KEY` - the AWS secret access key of the same IAM user

- `SECRET_NAME` - the name of the AWS secret to be retrieved

- `SECRET_DEST_PATH` - the destination path on disk to store the secret. Usually this is a shared volume. Defaults to `/var/run/secrets/aws-sm/.secret`.

## Example Usage

- AWS credentials retrieved from a generic Kubernetes secret 'aws-secret', containing id, key, and region key/value pairs
- AWS secret name 'myapp-creds'

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: alpine
  name: alpine
spec:
  securityContext:
    runAsUser: 1001
    fsGroup: 1001
  volumes:
  - name: app-secret
    emptyDir:
      medium: Memory
  initContainers:
  - image: stoyansbg/aws-sm-kubernetes-helper:v.1.0-runonce
    name: aws-sm-init
    volumeMounts:
    - name: app-secret
      mountPath: /var/run/secrets/aws-sm
    env:
      - name: SECRET_NAME
        value: myapp-creds
      - name: AWS_ACCESS_KEY_ID
        valueFrom:
          secretKeyRef:
            name: aws-secret
            key: id
      - name: AWS_SECRET_ACCESS_KEY
        valueFrom:
          secretKeyRef:
            name: aws-secret
            key: secret
      - name: AWS_REGION
        valueFrom:
          secretKeyRef:
            name: aws-secret
            key: region
    securityContext:
      allowPrivilegeEscalation: false
  containers:
  - image: alpine
    name: alpine
    command: ["sleep"]
    args: ["3600"]
    volumeMounts:
    - name: app-secret
      mountPath: /var/run/secrets/aws-sm

  # ...
```

## References

   https://github.com/sethvargo/vault-kubernetes-authenticator  
   https://github.com/sethvargo/vault-init  
   AWS Secrets Manager sample code  
