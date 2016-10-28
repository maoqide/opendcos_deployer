# opendcos_deployer
deploy opendcos with rest-api

## requirements
- a machine installed centos 7
- docker installed


## api
### create cluster
    method: post
    url: ip:port/v1/deploy/cluster
    contentType: json
*json example*
```json
{
    "userName": "mao",
    "skipInstall": true,
    "privateKey": "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",
    "privateNicName": "enp0s8",
    "config": {
    "agent_list": [
        "192.168.56.111",
        "192.168.56.112"
    ],
    "bootstrap_url": "file:///opt/dcos_install_tmp",
    "cluster_name" : "dcos1",
    "exhibitor_storage_backend": "static",
    "ip_detect_filename": "/genconf/ip-detect",
    "master_discovery": "static",
    "master_list": [
        "192.168.56.110"
    ],
    "process_timeout": 600,
    "resolvers": [
        "8.8.8.8",
        "8.8.4.4"
    ],
    "ssh_port": 22,
    "ssh_user": "root"
    }
}
```
- skipInstall: false if you want to install prereqs on cluster nodes  
- privateKey: Base64 encoded private key  
- privateNicName: network name used on each cluster node
- config: config parameters of DC/OS config.yaml


### add nodes
    method: post  
    url: ip:port/v1/deploy/nodes  
    contentType: json  
*json example*
```json
{
  "userName": "mao",
  "clusterName": "dcos1",
  "privateNicName": "enp0s8",
  "sshUser": "root",
  "nodes": [
    {
        "ip": "192.168.56.106",
        "skipInstall": true,
        "slaveType": "slave"
    },
    {
        "ip": "192.168.56.110",
        "skipInstall": true,
        "slaveType": "slave"
    }
  ]
}
```
- slaveType: 'slave' or 'slave_public'
- skipInstall: false if you want to install prereqs on cluster nodes
- privateNicName: network name used on each cluster node


### delete cluster
    method: delete  
    url: ip:port/v1/deploy/cluster/$(username)/$(clustername)


### delete node
    method: delete  
    url: ip:port/v1/deploy/nodes/$(username)/$(clustername)/$(nodeip)
