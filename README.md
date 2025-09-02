# Autonas

Automatic script(s) to deploy (&amp; update) docker compose files, that uses [Homepage](https://gethomepage.dev/) as a dashboard.

The main idea is to represent all the deployed versions of the services in this git repo, which will allow an easy redployment and update of the servcies.
Minimal custom configuration should be contained in a single file `config.env`.

## Requirements

a Linux environement with `bash`, `git` and `docker compose` installed

## What does it do

The script `deploy-auto-nas.sh` is meant to be run in a cron job, it will :

1. read the configuration files (`config.env` will override the default file)
2. update the files (&amp; versions) using `git pull`
3. copy the services folder into `SERVICES_PATH`
4. generate the configuraion file `services.yaml` for `Homepage`
5. run `docker compose up` on all the activated services

## How to use

start by cloning the repo

```bash
git clone https://github.com/omar-kada/autonas.git
```

then copy `config.example.env` to `config.env` and fill the configuration variables (adjust PULL & STOP as needed)

at last, run `deploy-auto-nas.sh`

```bash
./deploy-auto-nas.sh
```

## Global configuration

- **PULL** : (default = 1) when set to 0 it will disable pulling from the git repo
- **STOP** : (default = 0) when set to 0 it will disable stopping docker containers before redeploy
- **HOST** : hostname (needed for Homepage configuration for example)
- **SERVICES_PATH** : directory that will contain the services compose files and the generated .env variables file with any configuraiton included
- **DATA_PATH** : directory where all the containers data will be stored
- **SERVICES** : comma-seperated string of enabled services that will be deployed

## Service-specific configuraiton

to add service specific configuration, you just need to add a section `[service_name]` and below it all the configuration needed
the main properties that will be used for each service are

- **PORT** : the port where the service will be exposed
- **VERSION** : the version of the image
- **DESCRIPTION** : optional text that will be displayed in `Homepage` (defaults to service name)
- **ICON** : optional icon name for `Homepage` (defaults to service name)

any other configuration will be copied into the `.env` file related to each service

## Example of `config.env`

```ini
[global]
HOST=my-nas

SERVICES_PATH=/mnt/pool/autonas/services
DATA_PATH=/mnt/pool/autonas/data

SERVICES=homepage,immich,dockge

[homepage]
PORT=3210

[immich]
UPLOAD_LOCATION=/mnt/pool/autonas/data/syncthing/immich/library
DB_PASSWORD=MyStrongPassword

```
