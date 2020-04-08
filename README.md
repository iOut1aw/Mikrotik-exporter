# Mikrotik Prometheus Exporter

[![GitHub license](https://img.shields.io/github/license/hatamiarash7/Mikrotik-Exporter)](https://github.com/hatamiarash7/Mikrotik-Exporter/blob/master/LICENSE) ![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/hatamiarash7/mikrotik-exporter) ![Docker Image Version (latest by date)](https://img.shields.io/docker/v/hatamiarash7/mikrotik-exporter?sort=date)

![banner](banner.jpg)

It's a Prometheus exporter for mikrotik devices

## Prerequisites

You should create a user on the device that has API and read-only access :

```mikrotik
/user group add name=prometheus policy=api,read,winbox
/user add name=prometheus group=prometheus password=changeme
```

## Install

There is a simple Docker image that you can use it on your host / stack :

```bash
docker pull hatamiarash7/mikrotik-exporter:1.2
```

## Run

You can use this image in two ways :

#### Single device

```bash
docker-compose -f docker-compose-single.yml up -d
```

Or using `docker run` :

```bash
docker run -e DEVICE=home -e ADDRESS=192.168.1.4 -e USER=prometheus -e PASSWORD=changeme -p 9436:9436 hatamiarash7/mikrotik-exporter:1.2
```

#### Multiple device

To monitor multiple devices in your network you should create a `config.yml` file :

```yaml
devices:
  - name: home
    address: 192.168.1.1
    user: foo
    password: bar
  - name: company
    address: 192.168.2.1
    user: test
    password: 1234

features:
  bgp: true
  dhcp: true
  dhcpl: true
  routes: true
```

Then :

```bash
docker-compose up -d
```

## Access

Open [http://localhost:9436/metrics](http://localhost:9436/metrics)

## Dashboard

You can [Mikrotik Dashboard - Prometheus](https://grafana.com/grafana/dashboards/12055)

## Configure

There is some options and features that you can use them to configure your devices.  
**IMPORTANT : You can only use this options in `config.yml`**

| Option      | Description                                             | Default  |
| ----------- | ------------------------------------------------------- | -------- |
| address     | address of the device to monitor                        |
| config-file | config file to load                                     |
| device      | single device to monitor ( name )                       |
| insecure    | skips verification of server certificate when using TLS | false    |
| log-format  | logformat text or json                                  | json     |
| log-level   | log level                                               | info     |
| path        | path to answer requests on                              | /metrics |
| password    | password for authentication for single device           |          |
| deviceport  | port for single device                                  | 8728     |
| port        | port number to listen on                                | :9436    |
| timeout     | timeout when connecting to devices                      |          |
| tls         | use tls to connect to routers                           | false    |
| user        | user for authentication with single device              |          |

| Feature      | Description                               | Default |
| ------------ | ----------------------------------------- | ------- |
| with-bgp     | retrieves BGP routing infrormation        | false   |
| with-routes  | retrieves routing table information       | false   |
| with-dhcp    | retrieves DHCP server metrics             | false   |
| with-dhcpl   | retrieves DHCP server lease metrics       | false   |
| with-health  | retrieves board Health metrics            | false   |
| with-poe     | retrieves PoE metrics                     | false   |
| with-optics  | retrieves optical diagnostic metrics      | false   |
| with-w60g    | retrieves w60g interface metrics          | false   |
| with-wlansta | retrieves connected wlan station metrics  | false   |
| with-wlanif  | retrieves wlan interface metrics          | false   |
| with-monitor | retrieves ethernet interface monitor info | false   |
| with-ipsec   | retrieves ipsec metrics                   | false   |

You can see a complex example of `config.yml` :

```yaml
devices:
  - name: router
    address: 10.1.1.25
    user: prometheus
    password: changeme
  - name: second_router
    address: 192.168.1.8
    port: 8999
    user: prometheus
    password: changeme
  - name: domain
    srv:
      record: mikrotik.example.com
    user: prometheus
    password: changeme
  - name: custom_dns
    srv:
      record: mikrotik2.example.com
      dns:
        address: 8.8.8.8
        port: 53
    user: prometheus
    password: changeme

features:
  bgp: true
  dhcp: true
  dhcpl: true
  routes: true
```

**Note :** You can use `srv` instead of `address` to access your mikrotik by domain

## To-Do

- [ ] Test IPv6
- [ ] Other architectures ( Arm, Arm64, 386 )
- [ ] Extra OS support ( FreeBSD, NetBSD )
- [ ] Create Docker Swarm example
- [ ] Create Kubernetes example
- [x] Create specific Grafana dashboard

## Support

[![ko-fi](https://www.ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/D1D1WGU9)

## Contributing

1. Fork it!
2. Create your feature branch : `git checkout -b my-new-feature`
3. Commit your changes : `git commit -am 'Add some feature'`
4. Push to the branch : `git push origin my-new-feature`
5. Submit a pull request :D

## Issues

Each project may have many problems. Contributing to the better development of this project by reporting them
