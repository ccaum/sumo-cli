# Sumo Logic CLI

`sumo` is a CLI tool to manage content and interact with the Sumo Logic platform.

**Note this tool is currently in beta. Issues and pull requests are welcome!**

## Getting Started

1. Download the release for your platform from the [Releases page](releases)

2. Rename the file and make executable with `mv sumo-<platform> sumo && chmod +x sumo`

3. Configure your access credentials in `$HOME/.sumo-cli` (See Getting Access Credentials below)

4. Run `./sumo --help` to learn about what the command can do


## Why this exists

This command line tool aims to provide the following:

- A way to manage Sumo Logic content (dashboards, saved searches, etc.) as code in version control
- A familiar interface for perform common tasks such as searches and exporting content (coming soon!)
- A way to automate data and reporting workflows through scripting (coming soon!)

### Automating workflows

With the sumo CLI, you can automate the acceptance of new versions from
applications in the Sumo Logic app catalog, including the automatic application
of organization and team level modifications to dashboards and saved searches.

### Managing Sumo Logic application content as code

Sumo Logic application content can be fully managed as code, including the ability
to manage customizations of applications from the Sumo Logic App Catalog. The content,
such as dashboards, folders, and saved searches, are composed of reusable components,
such as panels and filter variables. Each component can be defined once and referenced
in multiple places.

**Composability**

Applications managed as code with this tool are fully composable. A dashboard panel
can be defined once and reused in multiple dashboards. A single dashboard can be referenced
in multiple folders. If any component is updated, every other component that references
the updated component will get the update.

**Idempotency**

The application code is fully idempotent and be applied multiple times to your Sumo Logic
account, always ensuring what's in your Sumo Logic account is always the desired state.
This allows you to apply your application code on a schedule. If the application is found
to be out of the desired state, this tool will bring it back into its correct configuration.

## How it works

### How to manage application content

#### Importing content from Sumo Logic

#### Applying app modifications in overlays

#### GitOps - Automating development workflows in GitHub


## Roadmap

- [ ] Download dashboards as PNG or PDF
- [ ] Run searches and streaming results to stdout
- [ ] Provide automatic installation of common Sumo Logic apps like Kubernetes and GitHub
- [ ] Manage collectors as code
- [ ] Manage FERs as code
- [ ] Manage parsers as code
- [ ] Manage monitors as code