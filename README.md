# Applier-CLI

Applier-CLI is a tool for initializing and working with [OpenShift-Applier](https://github.com/redhat-cop/openshift-applier) inventories.

[![asciicast](https://asciinema.org/a/240606.svg)](https://asciinema.org/a/240606)

## Installation

OS X & Linux:

```sh
go get github.com/redhat-cop/applier-cli
```

## Commands

All commands here are subcommands of the main `applier-cli` binary.

---

### init

Scaffolds an empty OpenShift-Applier project with the following structure:

```
[current directory]
├── apply.yml
├── files
├── inventory
│   ├── group_vars
│   │   └── all.yml
│   ├── hosts
│   └── host_vars
│       └── localhost.yml
├── params
├── roles
│   └── [large OpenShift-Applier role with many files]
├── requirements.yml
└── templates
```

#### Usage

This command does not take any input

#### Example

```sh
↪ applier-cli init
Initialized empty openshift-applier inventory
Successfully installed ansible-galaxy requirements
```

#### Notes

:information_source: Applier-CLI uses the GitHub API to determine the latest release of OpenShift-Applier, and writes that version to the `requirements.yml` file.

:information_source: Applier-CLI attempts to invoke `ansible-galaxy` to download the OpenShift-Applier role into the `roles` directory. If `ansible-galaxy` is unavailable, the `roles` directory will not be available.

---

### add

Adds resources to the current OpenShift-Applier inventory. This command can source resources from files on the local machine, or from an existing OpenShift cluster that the user is authenticated to via the `oc` CLI. Resources which are of a kind other than `Template` can optionally be templatized. Template resources are added to the `templates` directory, and a corresponding parameters file is initialized in the `params` directory. Resources which are not templates are simply added to the `files` directory.

#### Usage

This command requires either the `--from-cluster` (`-c`) or the `--from-file` (`-f`) flags, indicating where to find the resource to be added to the inventory.

This command can accept the `--make-template` (`-t`) flag, which will wrap any non-template resource in a template before adding it to the inventory.

Optionally, the `--edit` (`-e`) flag can be passed, which opens the newly added resource in your default editor immediately after adding it to the inventory.

#### Example

```sh
↪ applier-cli add --from-cluster deployment/test-deployment --make-template --edit
<editor is invoked at this point>
Template added to the current inventory.
```

```sh
↪ applier-cli add --from-file ~/project/some_resource.yml
File added to the current inventory.
```

#### Notes

:information_source: If using the `--edit` flag, the path to your default editor should either be set as the environment variable `EDITOR`, or set as the `editor` key in a `$HOME/.applier-cli.yaml` file. If neither of these are set, the editor default to `vi`.

:information_source: The `--edit` flag is particularly useful when used in combination with the `--make-template` flag, as the customization of parameters is likely desired when `--make-template` is used.

---

### get-latest-version

Displays the latest version of the OpenShift-Applier, as fetched by the GitHub Release API.

#### Usage

This command does not take any input

#### Example

```sh
↪ applier-cli get-latest-version
v2.0.8
```

#### Notes

:information_source: This command does not touch the fileysstem and does not require you to be in an OpenShift-Applier inventory to operate.

---

### run

Runs the current inventory against the cluster that the local `oc` client is logged in to.

#### Usage

With no flags set, `run` executes `ansible-playbook` on your local machine. There is also a `--docker` flag, which uses the OpenShift-Applier Docker container to run the playbook.

#### Example

```sh
↪ applier-cli run
<ansible-playbook output>
```

```sh
↪ applier-cli run --docker
<docker output... which should also look like an ansible playbook>
```

#### Notes

:information_source: When the `--docker` flag is used, this command attempts to determine the current state of SELinux on the local machine. If SELinux is `enforcing`, this command adds`:z` to the Docker volume mount to ensure that filesystem permissions are correct.

---

## Release History

* 0.0.1
  * Work in progress

## Meta

Distributed under the Apache License, v2.0. See ``LICENSE`` for more information.

## Contributing

1. Fork it (<https://github.com/redhat-cop/applier-cli/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
