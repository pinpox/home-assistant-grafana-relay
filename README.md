# home-assistant-grafana-relay

 
 
<p align="center">
 <img height="400" src="https://user-images.githubusercontent.com/1719781/138470470-d96ed6b8-0a07-44ef-8af3-7feb7e0f01f2.png">
   <a href="https://grafana.com">Grafana</a> ❤️ <a href="https://www.home-assistant.io">Home-assistant</a>
</p>

Listens for alerts via webhooks from [Grafana](https://grafana.com) and relays
as notifications using [Home-assistant](https://www.home-assistant.io/)

## Configuration

Configuration is done via environment variables. 

| Variable         | Description                 | Example                                             |
|------------------|-----------------------------|-----------------------------------------------------|
| `AUTH_TOKEN`     | Home-assistant auth token   | `LONG_LIVED_ACCESS_TOKEN`                           |
| `HM_SERVICE_URI` | Home-assistant API endpoint | `http://home.domain.tld/api/services/notify/notify` |
| `LISTEN_PORT`    | Port to listen on           | `12000`                                             |
| `LISTEN_HOST`    | Adress to listen on         | `localhost`                                         |

### Home-assistant

Generate a [long-lived access
token](https://developers.home-assistant.io/docs/auth_api/#long-lived-access-token)
in home-assistant. Tokens can be created in the profile section
(`https://home.domain.tld/profile`).

### Grafana

Create a new notification channel in grafana under
`https://mydomain.tld/alerting/notification/new` of type `webhook`. Make sure to
check the `Include Image` checkbox. As URL enter where this relay service will
be listening, e.g. `localhost:12000`.

### Example

In production you might want to write systemd service. A minimal working setup
looks like this.

```
export AUTH_TOKEN="LONG_LIVED_ACCESS_TOKEN"
export HM_SERVICE_URI="http://home.domain.tld/api/services/notify/notify"
export LISTEN_PORT="12000"
export LISTEN_HOST="localhost"
./home-assistant-grafana-relay
```

## NixOS module

A flake.nix file is included for compatibility with [Nix
Flakes](https://nixos.wiki/wiki/Flakes) for those that wish to use it as a
module. A bare-minimum flake.nix would be as follows:

```nix
{
  description = "NixOS configuration";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    ha-relay.url = "github:pinpox/home-assistant-grafana-relay";
  };

  outputs = { home-manager, nixpkgs, ... }: {
    nixosConfigurations = {
      hostname = nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
			./configuration.nix
			ha-relay.nixosModules.ha-relay
			{
				pinpox.services.home-assistant-grafana-relay = {
					enable = true;
					listenHost = "localhost";
					listenPort = "12000";
					haUri = "https://home.domain.com/api/services/notify/notify";
					envFile = "/var/secrets/ha-envfile";
				};
			}
        ];
      };
    };
  };
}
```

The envfile can be used to provide the token without putting it in the nix
store.

```env
AUTH_TOKEN="LONG_LIVED_ACCESS_TOKEN"
```

