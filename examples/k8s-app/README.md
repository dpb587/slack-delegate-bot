# slack-delegate-bot/k8s-app

An example for deploying the bot to Kubernetes.

## Preparation

Create a directory for bot configuration ([docs](../../README.md#Configuration)). See [cloudfoundry/slack-interrupt-bot](https://github.com/cloudfoundry/slack-interrupt-bot) for some real-world examples.

```
CONFIG_DIR=$HOME/workspace/slack-delegate-bot-config
mkdir -p $CONFIG_DIR
```

Create a `.env` file with the required `SLACK_TOKEN` environment variable and any other needed variables:

```yaml
SLACK_TOKEN=xoxb-...snip...
PAIRIST_PASSWORD_$NAME=...snip...
```

## Kubernetes

This example uses the following resources:

 * `slack-delegate-bot-env` - Secret with runtime environment variables (e.g. from `.env`)
 * `slack-delegate-bot-config` - ConfigMap with the bot configuration files (e.g. from `$CONFIG_DIR`)
 * `slack-delegate-bot` - Deployment managing the bot

You can create the resources with the following commands:

```
kubectl create secret generic slack-delegate-bot-env --from-env-file=.env
kubectl create configmap slack-delegate-bot-config --from-file=$CONFIG_DIR
kubectl apply -f deployment.yml
```

## Docker

Start the bot:

```console
docker run --rm --env-file=.env --volume $CONFIG_DIR:/config \
  docker.pkg.github.com/dpb587/slack-delegate-bot/slack-delegate-bot \
  --config=/config/*.yml \
  --config=/config/default.delegatebot \
  run
```
