# Gitlab Webhook

Used to receive push events from Gitlab to trigger pipeline runs.

## Developing

### Workflow

Build the server:

```bash
run build
```

Run the server:

```bash
./gitlab-webhook
```

Start an ngrok tunnel in another terminal window:

```bash
ngrok http 9090
```

Make sure the integration in your test project is configured
to point at the ngrok HTTP url and that SSL verification is
disabled. For more info see the [docs](https://docs.gitlab.com/ee/user/project/integrations/webhooks.html).
