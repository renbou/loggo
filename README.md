# Loggo

Loggo is an incredibly lightweight log aggregation system which is stupidly simple to setup and use. If you, like me, are tired from the complexity of setup and configuration of systems such as Graylog, Grafana Loki, and others, then Loggo should suit you pretty well.

![Preview](/docs/preview.png)

## Getting started

To understand how easy it is to setup and use loggo, I suggest checking out the [example](/example), which sets up a server instance and runs the loggo client alongside the app like this:

```yaml
# compose.yml
services:
  ...
  loggo:
    image: ghcr.io/renbou/loggo:latest
    ports:
      - "20080:20080" # web panel for log-viewing
```

```Dockerfile
# Dockerfile
FROM ghcr.io/renbou/loggo:latest as loggo
FROM python:3.11.7-slim-bookworm
...
COPY --from=loggo /bin/loggo /bin/loggo
CMD poetry run app | loggo pigeon --hq.addr loggo:20081
```

Loggo containers and binaries can be built for any platform, so you aren't even forced to use it inside containers. As long as you can launch the client process alongside the process from which you want to collect logs, you should be good to go!

## Design

Loggo was designed and built with a few key points in mind:

- Painless to run, which is crucial, in my opinion, for simple home-grown and other small/medium projects;
- Performant enough to collect logs from multiple sources, and then be able to filter through them;
- Fault-tolerant, so client logs wouldn't be lost if the server goes down.

To suit the storage-related needs, Loggo uses the awesome embeddable [BadgerDB](https://github.com/dgraph-io/badger) for storing logs in a flattened format with timestamps as keys, all of which allows filtering old logs and streaming newly incoming filtered logs. For transportation from the client to the server, gRPC streams with acknowledgements and reestablishment are used, which means that logs shouldn't be lost in any realistic case. Finally, to make this whole system easily usable, the gRPC-web API allows viewing and streaming logs filtered using an expression language with a Graylog-like syntax.

## Roadmap

With the current state of Loggo it's safe to say that usability-related tasks are more high-priority than others, such as performance, but since Loggo was designed and is already built with performance in mind, there shouldn't be any overly difficult to overcome bottlenecks. That being said, here are some of the features I'd like to support in Loggo in the near future:

- Configuring the whole system through the web interface, so that there's no need to touch the YAML file;
- More ways to receive logs on the client, e.g. files, syslog, etc;
- Indexing for custom log fields in the DB, so that most searches can be filtered faster.

Overall Loggo progress can be tracked through the [issues](https://github.com/renbou/loggo/issues) and linked [project](https://github.com/users/renbou/projects/2).
