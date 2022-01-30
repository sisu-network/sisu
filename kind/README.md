# kind

This directory contains a Makefile that is intended to be invoked via the top-level make target `create-kind-cluster`.
It creates a [kind](https://kind.sigs.k8s.io/) cluster and deploys into it:

- one MySQL instance in a Deployment
- one deyes instance in a Deployment
- one dheart instance in a Deployment
- one sisud instance in a Deployment
- and two ganache instances into a StatefulSet, into a separate `ganache` namespace

**NOTE**: currently sisud crashes at startup due to misconfiguration.
Getting sisud to run properly is the next milestone for a local cluster.
That will likely involve modifying `kind/include/sisu.toml` and updating `kind/Dockerfile.sisu-all`.

## Why kind?

kind is widely available and gives a realistic Kubernetes cluster running inside Docker.
It serves well as a proof of concept that your Kubernetes resource definitions work.
The resource files defined in kind could be trivially ported to a production Kubernetes cluster;
and it is simple to replicate production behavior such as a pod disappearing without a clean shutdown, for local testing.

## Interacting with the cluster

### What is running?

To see what is running in the cluster, you can run:

```console
% kubectl --context kind-sisu -n sisu get pods
NAME                     READY   STATUS             RESTARTS   AGE
deyes-857c8b89b-6hc7f    1/1     Running            3          7m45s
dheart-f75c87445-lp7kj   1/1     Running            3          7m45s
mysql-545ccdb445-lb9hd   1/1     Running            0          7m34s
sisud-6bdf6bb89c-86g79   0/1     CrashLoopBackOff   6          7m45s
```

### Inspecting a pod

If you need to inspect something in a running pod, you can run this (use a real pod name):

```console
% kubectl --context kind-sisu -n sisu exec -it dheart-f75c87445-lp7kj -- bash
root@dheart-f75c87445-lp7kj:~# cat /root/dheart.toml
port = 5678

sisu-server-url = "http://sisud:25456"

[db]
host = "mysql"
port = 3306
username = "root"
password = "password"
schema = "deyes"
```

deyes, dheart, and sisud all use the same locally built Docker image, so those commands are all available in any of the images.

### Reach a pod from your host machine

If you want to interact with a pod from the host machine, you can run e.g. `kubectl --context kind-sisu -n sisu port-forward service/dheart 5678`.
This creates a listener on port 5678 of your host machine, that proxies requests to port 5678 of the Kubernetes Service named `dheart`.
That command will run in the foreground until you ^C it.

Likewise you can also run `kubectl --context kind-sisu -n sisu port-forward service/deyes 31001`.
or `kubectl --context kind-sisu -n sisu port-forward service/sisud 25456`.

In the future we can configure the Kind cluster to automatically create a listener for those services so that you can simply access
a corresponding port on localhost.

### Reaching a pod from another pod

Within any pod in the `sisu` namespace, you should be able to resolve `dheart`, `deyes`, `sisud`, or `mysql` simply through that name.
Because ganache runs in a StatefulSet, in a separate namespace, you need to qualify `ganache-N.ganache.ganache` (that is `POD_NAME.SERVICE_NAME.NAMESPACE`).

```console
% kubectl --context kind-sisu -n sisu exec -it dheart-f75c87445-lp7kj -- bash
root@dheart-f75c87445-lp7kj:~#  (apt update && apt install -y curl) 2>&1 > /dev/null
debconf: delaying package configuration, since apt-utils is not installed
root@dheart-f75c87445-lp7kj:~# curl -v deyes:31001
*   Trying 10.96.66.244:31001...
* Connected to deyes (10.96.66.244) port 31001 (#0)
> GET / HTTP/1.1
> Host: deyes:31001
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Sat, 29 Jan 2022 03:57:40 GMT
< Content-Length: 0
<
* Connection #0 to host deyes left intact
root@dheart-f75c87445-lp7kj:~# curl -v mysql:3306 # expected to connect, not expected to do anything meaningful
*   Trying 10.244.0.10:3306...
* Connected to mysql (10.244.0.10) port 3306 (#0)
> GET / HTTP/1.1
> Host: mysql:3306
> User-Agent: curl/7.74.0
> Accept: */*
>
* Received HTTP/0.9 when not allowed

* Closing connection 0
curl: (1) Received HTTP/0.9 when not allowed

root@dheart-f75c87445-lp7kj:~# curl -v ganache-0.ganache.ganache:7545
*   Trying 10.244.0.8:7545...
* Connected to ganache-0.ganache.ganache (10.244.0.8) port 7545 (#0)
> GET / HTTP/1.1
> Host: ganache-0.ganache.ganache:7545
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404
< Content-Type: text/plain
< Content-Length: 13
<
* Connection #0 to host ganache-0.ganache.ganache left intact
404 Not Found
root@dheart-f75c87445-lp7kj:~# curl -v ganache-1.ganache.ganache:7545
*   Trying 10.244.0.9:7545...
* Connected to ganache-1.ganache.ganache (10.244.0.9) port 7545 (#0)
> GET / HTTP/1.1
> Host: ganache-1.ganache.ganache:7545
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404
< Content-Type: text/plain
< Content-Length: 13
<
* Connection #0 to host ganache-1.ganache.ganache left intact
```

(If this approach is not desirable, it is trivial to change Deployments to StatefulSets or vice versa,
so that exact instances can be reached (like ganache) as opposed to being load balanced (like dheart, deyes, and sisud right now,
even though they are only a single instance each at the moment).)
