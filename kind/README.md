# kind

This directory contains a Makefile that is intended to be invoked via the top-level make target `kind-init-cluster`.
It creates a [kind](https://kind.sigs.k8s.io/) cluster and deploys into it:

- one MySQL instance in a Deployment, in the `mysql` namespace
- two ganache instances into a StatefulSet, into the `ganache` namespace
- and then into `SISU_COUNT` (defaults to 2) namespaces starting at `sisud-0`, `sisud-1`, ...:
    - one deyes instance in a Deployment
    - one dheart instance in a Deployment
    - one sisud instance in a Deployment

If you want a different number of sisu instances, you can run as e.g. `make kind-init-cluster SISU_COUNT=3`.

`ganache-0` and `ganache-1` can respectively be accessed from the host machine via `http://localhost:7545` and `http://localhost:7546`.
Likewise, `sisud-0` is available at `http://localhost:25456`, `sisud-1` is available at `http://localhost:25457`, and so on.

You can remove the cluster with `make kind-delete-cluster`.

## Why kind?

kind is widely available and gives a realistic Kubernetes cluster running inside Docker.
It serves well as a proof of concept that your Kubernetes resource definitions work.
The resource files defined in kind could be trivially ported to a production Kubernetes cluster;
and it is simple to replicate production behavior such as a pod disappearing without a clean shutdown, for local testing.

## Interacting with the cluster

### What is running?

To see what is running in the cluster, you can run:

```console
% kubectl --context kind-sisu -n sisu-1 get pods
NAME                     READY   STATUS             RESTARTS   AGE
deyes-857c8b89b-j9qkp    1/1     Running            3          18m
dheart-f75c87445-s72gc   1/1     Running            0          18m
sisud-57b7d4dbb4-qx98j   1/2     CrashLoopBackOff   8          18m
```

### Inspecting a pod

If you need to inspect something in a running pod, you can run this (use a real pod name):

```console
% kubectl --context kind-sisu -n sisu-1 exec -it dheart-f75c87445-lp7kj -- bash
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

At time of writing, we are still stabilizing sisud.
Kubernetes does not allow you to `exec` into a pod that is not running.
To work around that, we have a `sleep` container in the sisud pod which you can `exec` into in order to inspect files, etc.

```console
% kubectl --context kind-sisu -n sisu-1 exec -it sisud-57b7d4dbb4-qx98j -c sleep -- bash
root@sisud-57b7d4dbb4-qx98j:~# ls -al /root/.sisu/main
total 24
drwxr-xr-x 6 root root 4096 Feb  4 03:49 .
drwxrwxrwx 3 root root 4096 Feb  4 03:48 ..
drwxr-xr-x 2 root root 4096 Feb  4 03:48 config
drwx------ 8 root root 4096 Feb  4 03:49 data
drwx------ 2 root root 4096 Feb  4 03:48 keyring-test
drwxr-xr-x 3 root root 4096 Feb  4 03:49 private
```

### Logs for a failing pod

Suppose you have a pod that is failing or is in a crash loop, like this:

```console
% kubectl --context kind-sisu -n sisu-1 get po
NAME                     READY   STATUS             RESTARTS   AGE
deyes-857c8b89b-j9qkp    1/1     Running            3          27m
dheart-f75c87445-s72gc   1/1     Running            0          27m
sisud-57b7d4dbb4-qx98j   1/2     CrashLoopBackOff   10         27m
```

You want to see the logs for the failing pod.
You could just take a guess at which container is failing.
If you request the logs for the pod without specifying a container, you will see output like this:

```console
% kubectl --context kind-sisu -n sisu-1 logs sisud-57b7d4dbb4-qx98j
error: a container name must be specified for pod sisud-57b7d4dbb4-qx98j, choose one of: [sisud sleep] or one of the init containers: [populate-sisu-dir init-sisu]
```

In this case, since it is in a CrashLoopBackoff, as opposed to a failed to initialize state, I can confidently request the logs of the sisud container:

```console
% kubectl --context kind-sisu -n sisu-1 logs sisud-57b7d4dbb4-qx98j -c sisud
...
Error: error during handshake: error on replay: validator set is nil in genesis and still empty after InitChain
...
```

But if you wanted more detail about what was wrong, without guessing about which container's logs to inspect, you can describe the pod to see what is failing:

```console
% kubectl --context kind-sisu -n sisu-1 describe pod sisud-57b7d4dbb4-qx98j
...
Init Containers:
  populate-sisu-dir:
    Container ID:  containerd://a3a1779fffdc955d6868a4f76617c0ed69e21e71f4ba213acc3e8749df171a92
    Image:         sisu.test/sisu-all:latest
    Image ID:      sha256:67cd05fe85befb1166539a6316f0c50320420d80fcbe6e6372fee380e8cfd5ec
    Port:          <none>
    Host Port:     <none>
    Command:
      bash
      -c
      cp -R /root/.sisu/. /sisu/
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
...
Containers:
  sisud:
    Container ID:  containerd://e4a91f8da4d18c68418877b381e3a9a46fb9aacf7af88e5591baee9d62a9e502
    Image:         sisu.test/sisu-all:latest
    Image ID:      sha256:67cd05fe85befb1166539a6316f0c50320420d80fcbe6e6372fee380e8cfd5ec
    Port:          25456/TCP
    Host Port:     0/TCP
    Command:
      sisud
      start
    State:          Waiting
      Reason:       CrashLoopBackOff
    Last State:     Terminated
      Reason:       Error
      Exit Code:    1
```

### Reach a pod from your host machine

If you want to interact with a pod from the host machine, you can run e.g. `kubectl --context kind-sisu -n sisu-1 port-forward service/dheart 5678`.
This creates a listener on port 5678 of your host machine, that proxies requests to port 5678 of the Kubernetes Service named `dheart`.
That command will run in the foreground until you ^C it.

Likewise you can also run `kubectl --context kind-sisu -n sisu port-forward service/deyes 31001`.
or `kubectl --context kind-sisu -n sisu port-forward service/sisud 25456`.

In the future we can configure the Kind cluster to automatically create a listener for those services so that you can simply access
a corresponding port on localhost.

### Reaching a pod from another pod

Within any pod in a `sisu-N` namespace, you should be able to resolve `dheart`, `deyes`, or `sisud` simply through that name.
Because ganache runs in a StatefulSet, in a separate namespace, you need to qualify `ganache-N.ganache.ganache` (that is `POD_NAME.SERVICE_NAME.NAMESPACE`).
Likewise, you can resolve `mysql.mysql` (format `SERVICE_NAME.NAMESPACE`) from any pod in the cluster.

Here is some console output demonstrating how to resolve the other pods:

```console
% kubectl --context kind-sisu -n sisu-1 exec -it dheart-f75c87445-s72gc -- bash
root@dheart-f75c87445-s72gc:~# (apt update && apt install -y curl) 2>&1 > /dev/null
root@dheart-f75c87445-s72gc:~# curl -v deyes:31001
*   Trying 10.96.20.201:31001...
* Connected to deyes (10.96.20.201) port 31001 (#0)
> GET / HTTP/1.1
> Host: deyes:31001
> User-Agent: curl/7.74.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Date: Fri, 04 Feb 2022 04:10:50 GMT
< Content-Length: 0
<
* Connection #0 to host deyes left intact
root@dheart-f75c87445-s72gc:~# curl -v mysql.mysql:3306 # expected to connect, not expected to do anything meaningful*   Trying 10.244.0.13:3306...
* Connected to mysql.mysql (10.244.0.13) port 3306 (#0)
> GET / HTTP/1.1
> Host: mysql.mysql:3306
> User-Agent: curl/7.74.0
> Accept: */*
>
* Received HTTP/0.9 when not allowed

* Closing connection 0
curl: (1) Received HTTP/0.9 when not allowed

root@dheart-f75c87445-s72gc:~# curl -v ganache-0.ganache.ganache:7545
*   Trying 10.244.0.11:7545...
* Connected to ganache-0.ganache.ganache (10.244.0.11) port 7545 (#0)
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
root@dheart-f75c87445-s72gc:~# curl -v ganache-1.ganache.ganache:7545
*   Trying 10.244.0.12:7545...
* Connected to ganache-1.ganache.ganache (10.244.0.12) port 7545 (#0)
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
404 Not Found
```

(If this approach is not desirable, it is trivial to change Deployments to StatefulSets or vice versa,
so that exact instances can be reached (like ganache) as opposed to being load balanced (like dheart, deyes, and sisud right now,
even though they are only a single instance each at the moment).)
