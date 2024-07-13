# YoshiK3S

This is a Library made for the [Terraform Provider YoshiK3S](https://github.com/HideyoshiNakazone/terraform-provider-yoshik3s) project,
it is primarily a library but it can be used as a standalone CLI tool.

This tool is used to create, update and delete K3S nodes and it does not manage the state of the nodes,
it is recommended to use the [Terraform Provider YoshiK3S](https://github.com/HideyoshiNakazone/terraform-provider-yoshik3s)
to manage the state of the nodes via [Terraform](https://www.terraform.io/).

## Features
    
- [x] Create a K3s cluster
- [x] Create a K3s master node
- [x] Create a K3s worker node
- [x] Delete a K3s cluster
- [x] Delete a K3s master node
- [x] Delete a K3s worker node
- [x] Update a K3s cluster
- [x] Update a K3s master node
- [x] Update a K3s worker node
- [ ] ~~Validation of the Master Node on which the Worker Node is being created~~ 
  - (Not possible due to the lack of state management)


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Library

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

1. The binary will be available in the `$HOME/go/bin` directory

## Using the CLI

Substitute the placeholders with the desired values for your cluster.

```yml
cluster:
    version: ${K3S_VERSION}
    token: ${K3S_TOKEN}

master_nodes:
    - name: ${MASTER_NODE_NAME}
      connection:
          host: ${MASTER_NODE_HOST}
          port: ${MASTER_NODE_PORT}
          user: ${MASTER_NODE_USER}
          password: ${MASTER_NODE_PASSWORD}
      options:
        - ${MASTER_NODE_OPTION_1}
        - ${MASTER_NODE_OPTION_2}

worker_nodes:
    - name: ${WORKER_NODE_NAME}
      server_address: ${MASTER_SERVER_ADDRESS}
      connection:
          host: ${WORKER_NODE_HOST}
          port: ${WORKER_NODE_PORT}
          user: ${WORKER_NODE_USER}
          password: ${WORKER_NODE_PASSWORD}
      options:
        - ${WORKER_NODE_OPTION_1}
        - ${WORKER_NODE_OPTION_2}
```

Run the following command to create a cluster:

```shell
yoshik3s -config ${CONFIG_FILE}
```

Run the following command to delete a cluster:

```shell
yoshik3s -config ${CONFIG_FILE} -delete
```
