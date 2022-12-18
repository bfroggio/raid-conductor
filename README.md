Note: I no longer stream to Twitch. You're welcome to try using this project or to modify the code yourself, but I am unable to offer support if you get stuck.

# Raid Conductor

At the end of a stream, it's tiring and time consuming to have to pick who to raid with your current viewers. Raid Conductor simplifies that process by running an automated search over a pre-defined list of streamers. During this search, Raid Conductor prioritizes streamers based on your preferences, game dislikes, and language selection.

## Usage

1. Go to the [Twitch Developers](https://dev.twitch.tv/console/apps/create) page and create a Twitch application (note down your client ID and secret)
1. Create a `config.toml` file as outlined in the "Config File" section below
1. Launch `raid-conductor.exe` by double clicking on it (or launching it from a Stream Deck action)

## Config File

`config.toml` should be in the same directory as the `raid-conductor.exe` binary. There are a few properties that can go in your `config.toml` file. Properties listed below are optional unless otherwise noted. See `sample-config.toml` for an example with fake configuration values.

- `twitch_username` (required): The username for your Twitch account/channel. Raid Conductor can't read your channel's chat messages without this.
- `twitch_bot_username` (required): The username for the Raid Conductor chat bot. Can be the same as `twitch_username`. This user needs "Channel Editor" permissions to be able to start raids.
- `twitch_bot_secret`: The OAUTH token for the `twitch_bot_username` account. Needed to post messages to your Twitch chat (which is how raids are started).
- `twitch_client_secret` (required): The client secret for your Twitch application. Raid Conductor can't read talk to Twitch without this.
- `twitch_client_id` (required): The client ID for your Twitch application. Raid Conductor can't read talk to Twitch without this.
- `twitch_client_secret` (required): The client secret for your Twitch application. Raid Conductor can't read talk to Twitch without this.
- `priority_streamers` (required): The list of streamers you want to prioritize raiding.
- `backup_streamers` (required): The list of streamers to check after the prioritized list.
- `games_blacklist` (required): The list of games you don't want to send your viewers over to.
