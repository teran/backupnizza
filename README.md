# backupnizza

![Test & Build status](https://github.com/teran/backupnizza/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/teran/backupnizza)](https://goreportcard.com/report/github.com/teran/backupnizza)
[![Go Reference](https://pkg.go.dev/badge/github.com/teran/backupnizza.svg)](https://pkg.go.dev/github.com/teran/backupnizza)

A scheduler initially written to run restic but actually can run anything

## What it is and what it for

backupnizza is initially designed and developed to for a particular purpose:
to run [rclone](https://github.com/rclone/rclone) and [restic](https://github.com/restic/restic)
as a part of a single backup pipeline periodically and the most secure manner
possible. To achieve that backupnizza uses some techniques to securely pass
secrets to them. However it can be used with any scripts and tools as a cron-like
scheduler.

The overall approach used is to create PoC and evaluate it to MVP so it's not
designed as a full end-user product with amount of documentation, friendly
UIs and so on - it's rough and need some experience to work with.

## How it works

backupnizza relies on the following workflow:

* The user needs something (which is called "task") to run periodically
* Each task allows to get secrets from the environment or, even better,
    run some command to get the secret
* Secrets are initially stored in any keychain application like macOS Keychain,
    1Password, etc.
* User doesn't want to authenticate the task each time it runs according
    to schedule, assuming the authentication at system start is secure enough.

So backupnizza uses backup.yaml (or JSON - it doesn't matter) to read what
to run, how often to run and where to grab secrets for it reading them just once
at start, stores in-memory until the time to use it will come.

When task runs there are two options to get the secret:

* get it via environment variable
* get it via [secretbox](https://github.com/teran/secretbox) CLI

First case simply performs substitution to the task environment - here's nothing
special.

Second case accesses the scheduler daemon via CLI and retrieves the secret from it
with one-time token authentication so nobody else can reuse the token. If the token
won't be used within the time specified in `max_token_ttl` - it revokes automatically
and cannot be used.

## What it all for

Both of [rclone](https://github.com/rclone/rclone) and [restic](https://github.com/restic/restic)
allows to accept secrets via environment variables but this reveals the case of
possibility to get them from both of Linux and macOS systems while rclone or restic
are running even without any special privileges while 1Password or keychain have
no options to authenticate particular app once and access to secrets while the
system (macOS) is locked.

## Security stuff

### Known issues

backupnizza is designed with security in mind but it doesn't mean it's not exposed
to some kind of attacks, here just some obvious by design:

* Cold boot attack - could allow to dump the whole memory including secrets
    from backupnizza (just like for almost any other application)
* There's a short period of time in rclone and restic before they're appear in
    process list (with one-time token) and they're actually redeem the token,
    typically <100ms. However it's technically possible to intercept the token
    from process list or process' environment variables and retrieve the secret
    from backupnizza. This risk has very low probability because of very short token
    TTL and the actual time between token is generated and redeemed but it's
    still present.
* Third party application binary spoofing - backupnizza doesn't verify the binaries
    it runs so it's technically possible to replace them with any malicious ones.

### Configuration tips

#### Set immutable flag to all of the files

All the binaries and configuration (which contain exact commands to run) must
be secure to make the whole process secure. So adding some third-party tool
to the task list could cause secrets exposure after backupnizza is
restarted (since configuration is read once at start).
