# Provider SDK

> This libray help you building a Clevercloud addon provider
> more informations on the [documentation](https://developers.clever-cloud.com/doc/extend/add-ons-api/)

## Setup a new addon provider

### Write the manifest

You need to write on a file the next lines and adapt it

```json
// manifest.json
{
  "id": "provider-id",
  "name": "Provider Name",
  "api": {
    "config_vars": [ "ADDON_NAME_MY_VAR" ], // A list of environment your addon will expose to applications (must be prefixed by `PROVIDER_ID`)
    "regions": [ "par" ], // Be sure the regions you put already exists on CleverCloud
    "password": "44ca82ddf8d4e74d52494ce2895152ee", // some random generated secret
    "sso_salt": "fcb5b3add85d65e1dddda87a115b429f", // idem
    "production": {
      "base_url": "https://yourservice.com/clevercloud/resources", // the endpoint we can contact your provider
      "sso_url": "https://yourservice.com/clevercloud/sso/login"
    },
    "test": {
      "base_url": "https://localhost:9000/clevercloud/resources", // same as production
      "sso_url": "https://localhost:9000/clevercloud/sso/login"
    }
  }
}
```

### Register your provider

On your [Console](https://console.clever-cloud.com/), choose an organisation and go to `create` > `an addon provider`.

Send the manifest file you just wrote.

Under `Plans` tab, add at least 1 plan for your addon.


> You now have an addon provider in ALPHA state, this means, only your organisation can order it.

### Code your provider

Create a new Golang project, and add the SDK as dependency

```sh
mkdir dummyprovider
cd dummyprovider
go mod init dummyprovider

go get go.clever-cloud.dev/provider

```

It's time to write some logic in your `main()`, here is an [example](./example/main.go).

### Run your provider on CleverCloud

```sh
clever-tools create --type go -r par -o ORGANISATION_ID dummy-provider
clever-tools env set CC_GO_BUILD_TOOL gobuild

git init
git add .
git commit -m 'init'
clever-tools deploy
```

> Don't forget to add a vhost (domain) on your app

## Addon provider client

As an addon provider you can perform several actions on the CleverCloud API

```golang
c := client.New(cfg)

c.ListAddons(ctx context.Context) ([]Addon, error)
c.GetAddon(ctx context.Context, addonID string) (*AddonInfo, error)
c.UpdateEnvironment(ctx context.Context, addonID string, environment map[string]string) error

```

## Useful commands

### Show addon providers informations

```sh
$ clever addon providers
$ clever addon providers show X
```