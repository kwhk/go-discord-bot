# go-discord-bot

## How to run locally

The following directory paths assume that you are in the project's root directory.
1. Install all modules
    ```go
    go mod download
    ```
1. Run bot:
    ```go
    go run ./cmd/go-discord-bot
    ```

## Possible features to implement

1. Rename voice channel to be the name and how long it's been open for
1. Reimplement @here when voice channel is active
1. Stat tracker that tracks:
    1. last online (instead of using command you can create role which shows timestamp)
    1. message count
    1. voice channels that were accessed and length of time
        1. can be used to show who was active on which voice channel
    1. create ephemeral roles that adjust according to the stats tracked
1. create a collage/gallery of all images posted with filters