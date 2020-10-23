# Raid Conductor
[![Man Hours](https://img.shields.io/endpoint?url=https%3A%2F%2Fmh.jessemillar.com%2Fhours%3Frepo%3Dhttps%3A%2F%2Fgithub.com%2Fbfroggio%2Fraid-conductor.git)](https://jessemillar.com/r/man-hours)

At the end of a stream, it's tiring and time consuming to have to pick who to raid with your current viewers. Raid Conductor simplifies that process by running an automated search over a pre-defined list of streamers. During this search, Raid Conductor prioritizes streamers based on your preferences, game dislikes, and language selection.

## Usage

1. Go to the [Twitch Developers](https://dev.twitch.tv/console/apps/create) page and create a Twitch application (note down your client ID and secret)
1. Create a `config.toml` file as outlined in the "Config File" section below
1. Launch `raid-conductor.exe` by double clicking on it (or launching it from a Stream Deck action)

## Config File

`config.toml` should be in the same directory as the `raid-conductor.exe` binary. There are a few properties that can go in your `config.toml` file. Properties listed below are optional unless otherwise noted. See `sample-config.toml` for an example with fake configuration values.

- `twitch_client_id` (required): The client ID for your Twitch application. Raid Conductor can't read talk to Twitch without this.
- `twitch_client_secret` (required): The client secret for your Twitch application. Raid Conductor can't read talk to Twitch without this.
- `priority_streamers` (required): The list of streamers you want to prioritize raiding.
- `backup_streamers` (required): The list of streamers to check after the prioritized list.
- `games_blacklist` (required): The list of games you don't want to send your viewers over to.
