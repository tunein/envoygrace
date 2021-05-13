# envoygrace

Graceful shutdown for Envoy running in a Sidecar model.

The problem today with Envoy running as a sidecar next to your application is that when Kubernetes goes and sends a `SIGTERM` to Envoy, it doesn't gracefully drain the active connections automatically, leading to errors and not a pleasant user experience.
There are many Github issues and articles that say to be able to gracefully shutdown envoy you need to send a `POST /healthcheck/fail` request to it's admin interface - that way it'll stop accepting new connections and drain existing connections.

Some links to follow for more information:

- https://github.com/envoyproxy/envoy/issues/2920#issuecomment-540709506
- https://www.envoyproxy.io/docs/envoy/latest/operations/admin#operations-admin-interface-healthcheck-fail
- https://github.com/envoyproxy/envoy/issues/1990#issuecomment-341509716

This package is an alternative to doing the follow in your K8s prestop hook:

- `curl -XPOST localhost:8001/healthcheck/fail`
- `sleep 5`

## Why use this package, when I can just use cURL?

- You don't want to install cURL in your final Docker image (but are willing to put this binary in it).
- Future versions of `envoygrace` can add additional functionality that you would otherwise need to write yourself.

## Installation

Download a version stamped binary from the [Releases](https://github.com/tunein/envoygrace/releases) page into your application's Dockerfile. For example:

```dockerfile
ARG ENVOYGRACE_VERSION=v1.0.0

RUN wget -qO/bin/envoygrace https://github.com/tunein/envoygrace/releases/download/${ENVOYGRACE_VERSION}/envoygrace-linux-amd64 && \
    chmod +x /bin/envoygrace
```

## Usage

Configure the binary to run on a preStop hook in your Kubernetes deployment.

```yaml
lifecycle:
  preStop:
    exec:
      command:
        - "/bin/envoygrace --envoy-url localhost:8001 --eviction-period 5000 --envoy-timeout 1000 --pre-sleep 5000"
```

**MAKE SURE YOU SET A PRESTOP HOOK ON YOUR ENVOY CONTAINER WITH SLEEP THAT'S LONGER THAN THE ABOVE CONFIG**

### Command Line Options

```
  --help
      shows this help output
  --envoy-timeout int
      envoy request timeout in milliseconds (default 1000)
  --envoy-url string
    	the base url where envoy lives (default "localhost:8001")
  --eviction-period int
    	amount of milliseconds to sleep after calling envoy graceful commands (default 5000)
  --pre-sleep int
    	how long to sleep before sending requests to envoy in milliseconds (default 5000)
```
