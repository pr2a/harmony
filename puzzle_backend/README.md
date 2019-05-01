# REST API for hexie

- `/enter?email=xxxx`
  - It responds `{balance, address, level}`
- `/finish?new_level=xxx&account=xxx&key=xxx`
  - It responds `{msg, level}`

# App Engine with golang

- Install [SDK](https://cloud.google.com/sdk/docs/) and followings

```
gcloud components update
gcloud components install app-engine-go
```

- To deploy in local, run following:

```
dev_appserver.py app.yaml
```

App engine project id: benchmark-209420
