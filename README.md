# Sumo Logic CLI

`sumo` is a CLI tool to manage content and interact with the Sumo Logic platform.

**Note this tool is currently in beta. Issues and pull requests are welcome!**

## Getting Started

1. Download the release for your platform from the [Releases page](releases)

2. Rename the file and make executable with `mv sumo-<platform> sumo && chmod +x sumo`

3. Configure your access credentials in `$HOME/.sumo-cli` (See [Getting Access Credentials below](#getting-access-credentials))

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

### Getting Access Credentials
1) Create an access ID/key pair by following [the documentation](https://help.sumologic.com/Manage/Security/Access-Keys#manage-your-access-keys-on-preferences-page)
2) Get your account's deployment region code by following [these instructions](https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-and-Firewall-Security#how-can-i-determine-which-endpoint-i-should-use)
3) Save your credentials in `$HOME/.sumo-cli` with the following format:

```
access-id: <your access ID>
access-key: <your access key>
deployment: <your deployment region>
```

### How to manage application content

#### Importing content from Sumo Logic
Import a content folder from your Sumo Logic account with the following:
`sumo app download-folder <folder ID> | sumo app import`

This will import your content to a `base` directory.

#### Overwriting base content
Individual component resources such as folders, dashboards, panels, saved-searches, and variables can be modified through overlays. An overlay is a place to put content modifications that will be merged with the parent overlay. 

There are three layers, a base layer and two overlay layers
- base
- middle
- final

Upstream content, such an app catalog release or content from another git repository, goes in `base`. Any components that have modification data in the `middle` overlay will overwrite the same component in the `base` layer. Likewise, component modifications in the `final` layer will overwrite anything in the `middle` layer. Note that if you modify a component in the `final` layer and that component is not defined at all in the `middle` layer, the modification effectively applies to the component in the `base` layer.


#### Applying app modifications in overlays

[TBD]

#### Performing and deploying a build
When it's time to push content to Sumo Logic, you can create a build with the following command:
`sumo app build`

This will create a file called `build.json` that contains all of the folders, dashboards, and saved searches defined in your application, including modifications defined in overlays.

Deploy the build to Sumo Logic with:
`sumo app push -d <parent folder ID> --overwrite`

The `--overwrite` flag will force any existing content in the parent folder to be replaced with the content defined in the `build.json` file. This can be run on a schedule to perform desired state reconsiliation in order to ensure our production content always matches the source of truth: the code.

#### GitOps - Automating development workflows in GitHub


## Roadmap

- [ ] Download dashboards as PNG or PDF
- [ ] Run searches and streaming results to stdout
- [ ] Provide automatic installation of common Sumo Logic apps like Kubernetes and GitHub
- [ ] Manage collectors as code
- [ ] Manage FERs as code
- [ ] Manage parsers as code
- [ ] Manage monitors as code
