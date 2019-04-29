# Firebase deployment and local environment

## Deployment

Steps to deploy the app

- You need to install firebase tools by

```
npm install -g firebase-tools
```

- You need to login to the firebase project. If you are not granted to the project, ask the owner (@minhdoan).

```
firebase login
```

- You will have to ask questions related to the configuration to pick the project.

- To deploy the app, you need to run

```
firebase deploy
```

## Local environment

To debug the app in local environment, you need to run

```
firebase serve --only functions
```

Then it will shows the REST API in your console.
