# Autonas

AutoNAS is a lightweight gitops automation tool that handles **Docker compose** stacks deployment.  
Its purpose is to make deploying and updating a self-hosted environment **simple**, **fast**, and **reproducible**. without the need to have heavy tooling or dependencies.

## Requirements

a Linux environement with `docker` installed

## Features

1. Easy service customization through a simple configuration file
2. Minimal, clean structure — no unnecessary scaffolding
3. No special syntax, compose stacks configuration are kept in their original form
4. Works on any Docker-capable system

## Getting started


1. Create a configuration repo containing all your compose stacks with this structure (for example : [AutoNAS Config](https://github.com/omar-kada/autonas-config))
```
services/
├── service1/
|   ├── compose.yaml
|   └── .env        
└── service2/
    ├── compose.yaml
    └── .env        
```

2. Copy  the `compose.yaml` file your system and fill the needed variables (make sure to read the comments about each variable), here are the main ones : 

```yaml
SERVICES_DIR : where the stack configuration will be stored
CONFIG_FILES : list of configuration files (last has more priority)
CONFIG_REPO : repository containing the stack definition
CONFIG_BRANCH: branch used to pull from the repo
CRON_PERIOD: cron schedule of when the periodic deployement will be executed
```

3. Create a `config.yaml` file and define the services you want to deploy :

```yaml
ENV_VAR: value # will be available in all services

services:
  service1: 
    ENV_VAR: override value # will override global value for this service
    SERVICE_SPECIFIC_VAR: another_value
  
  service2:
    Disabled: true # if Disabled, service will not be deployed
```

4. Run the stack using : 
```bash
docker compose up -d
```

