## Guide for developers of this application

## [Required Tools]

This project runs a real Kubernetes cluster (kind) for local development. Its an integration environment, so it requires
some real services to be running.

- [Go](https://go.dev/)
- [Air (Live Reload for Go)](https://github.com/cosmtrek/air)
- [Yarn](https://yarnpkg.com/)
- [Node](https://nodejs.org/en)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [helm](https://helm.sh/)
- [Kind](https://kind.sigs.k8s.io/)
- [Podman](https://podman.io/docs/installation)

## Setting up `podman` for the first time on MacOS

1. Run `brew install podman` 
2. Run `podman machine init`
3. Run `podman machine set --rootful` (required to run the postgres statefulset)
4. Run `podman machine start`


## Setting up an Oauth provider:
For local development, use of any oauth2 provider will work. (Zitadel, Okta, AzureEntraID etc)
- You will need an oauth provider in order to run this application as it relies on oauth to handle authentication 
  and to authorize users based on their ABAC designated permissions and groups. 

  You will need an OAuth2.0 client setup as a PKCE AuthCode flow client with a client ID. We recommend using a free tier provider like 
  [Zitadel](https://zitadel.com/) for local development. 

## Dev workflow

- First, start the client side of the app:
  - `cd ./client && yarn install && yarn start`
  - The client application will be running at http://localhost:3000
- Second, start the API server:
  - replace the `placeholders` in the local-secret.env file with real AWS access credentials (if using reports feature) -- You can skip this step if you are not working on the reports feature. 
  - execute sh `./localdev/start-kind.sh kind.yaml`
  - execute `air`from the root of this project
  - Navigate to http://localhost:3000 to view application from the front-end live reload server
  - View `hotload-backend.sh` for details on what `air` executes whenever you save server side code. Air does _not_ watch code in the `client` directory, but you can enable reload/rebuild of the client code directly. You should use `yarn start` from the client directory to handle hot reloads of the client code.
