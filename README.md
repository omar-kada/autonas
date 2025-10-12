# Autonas

AutoNAS is a simple tool that allows to handle docker compose stacks deployements through a configuration in a git repo, that uses [Homepage](https://gethomepage.dev/) as a dashboard.

## Requirements

a Linux environement with `docker` installed

## What does it do

1. read the configuration files
2. copy the services folder into `SERVICES_PATH`
3. generate `.env` file for each compose stack
4. run `docker compose up` on all the activated services

## How to use

first the tool is installed

**WIP**

start by copying `config.example.yaml` to `config.yaml` and fill the configuration variables

at last, run `autonas`

```bash
autonas run -c config.default.yaml,config.yaml -r https://github.com/omar-kada/autonas-config
```

## Global configuration

- **AUTONAS_HOST** : hostname (needed for Homepage configuration for example)
- **SERVICES_PATH** : directory that will contain the services compose files and the generated .env variables file with any configuraiton included
- **DATA_PATH** : directory where all the containers data will be stored
- **enabled_services** : list of enabled services that will be deployed

## Service-specific configuraiton

to add service specific configuration, you just need to add a section `<service_name>:` and below it all the configuration needed
the main properties that will be used for each service are

- **PORT** : the port where the service will be exposed
- **VERSION** : the version of the image
- **DESCRIPTION** : optional text that will be displayed in `Homepage` (defaults to service name)
- **ICON** : optional icon name for `Homepage` (defaults to service name)

any other fields will be copied into the `.env` file related to each service

## Example of `config.yaml`

```yaml
AUTONAS_HOST: "<hostname>"
SERVICES_PATH: "/path/to/directory/where/services/are/stored"
DATA_PATH: "/path/to/data/folder"

enabled_services: # list of services to install
  - homepage
  - dockge

services:
  homepage:
    PORT: 1234
```
