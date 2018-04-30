<div align="center">
    <h2>Actionizer</h2>
    <p align="center">
        <p>Add fun to your organization by randomly assigning funny tasks to people</p>
    </p>
</div>

<p align="center">
    <img src=".github/screenshot.png">
</p>

## Contents

* [Installation](#installation)
* [Usage](#usage)
* [Todo](#todo)
* [License](#license)

## Installation

Before installing `actionizer` you need to have mongodb installed alongside with its golang driver:
```shell
go get gopkg.in/mgo.v2
```

Once done install `actionizer` by running:
```shell
go install github.com/think-it-labs/actionizer
```

## Usage

```
$ actionizer -h
Usage of actionizer:
  -config string
        Configuration file (default "actionizer.json")
```

Actionzer read its configuration from a json file (`actionizer.json` by default) that has this format

```json
{
    "database": {
        "host": "localhost",
        "name": "actionizer",
        "user": "",
        "password": ""
    },
    "http_listen": "127.0.0.1",
    "http_port": 8080,
    "action_duration": "2w"
}
```

The value of `action_duration` field should be in form `{AMOUT}{UNIT}`

Unit can be: `s (seconds)`, `m (minutes)`, `h (hours)`, `d (days)` or `w (weeks)`.

### Config file Example

`database.user` and `database.password` can contains database user credentials if needed.

The value of `action_duration` here is "4d" which stands for every 4 days.

```json
{
    "database": {
        "host": "localhost",
        "name": "actionizer",
        "user": "sample_user",
        "password": "sample_password"
    },
    "http_listen": "127.0.0.1",
    "http_port": 8080,
    "action_duration": "4d"
}
```

## Todo

- [ ] Add subcommand that can add a user or an action 
- [ ] Slack integration
- [ ] Make it possible to mark a `task` as done (OKs)
- [ ] Consider not assigning actions that can only be done locally to remote persons

## License

This repository has been released under the [MIT License](LICENSE)

------------------
Made with ♥ by [Think.iT](http://www.think-it.io/).