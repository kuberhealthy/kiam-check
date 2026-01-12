# kiam-check

The `kiam-check` validates that [KIAM](https://github.com/uswitch/kiam) can intercept AWS metadata requests by listing Lambda functions with an assumed IAM role. The check reports success when it can list Lambda functions, or when the list matches a configured count.

## Configuration

Set these environment variables in the `HealthCheck` spec:

- `AWS_REGION` (optional): AWS region to query. Defaults to `us-west-2`.
- `LAMBDA_COUNT` (optional): expected number of Lambda functions. When set, the check requires an exact match.
- `DEBUG` (optional): set to `true` to enable debug logging.

The KIAM role is passed via `spec.extraAnnotations` using the `iam.amazonaws.com/role` annotation.

## Build

- `just build` builds the container image locally.
- `just test` runs unit tests.
- `just binary` builds the binary in `bin/`.

## Example HealthCheck

Apply the example below or the provided `healthcheck.yaml`:

```yaml
apiVersion: kuberhealthy.github.io/v2
kind: HealthCheck
metadata:
  name: kiam
  namespace: kuberhealthy
spec:
  extraAnnotations:
    iam.amazonaws.com/role: <role-arn>
  runInterval: 5m
  timeout: 15m
  podSpec:
    spec:
      containers:
        - name: kiam
          image: kuberhealthy/kiam-check:sha-<short-sha>
          imagePullPolicy: IfNotPresent
          env:
            - name: AWS_REGION
              value: us-west-2
          resources:
            requests:
              cpu: 15m
              memory: 10Mi
            limits:
              cpu: 30m
      restartPolicy: Always
```

## Install

Configure a valid IAM role ARN with permissions to read Lambda functions (for example, `AWSLambdaReadOnlyAccess`). Replace `<role-arn>` in `healthcheck.yaml` and apply it with `kubectl apply -f healthcheck.yaml`.
