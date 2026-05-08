# Create ZIP

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-create-zip?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-create-zip/releases)

Creates a ZIP from the given file/dir to the given destination.

<details>
<summary>Description</summary>

It stores symlinks as symlinks (it will not copy the file which the symlink is pointing to).

### Configuring the Step

1. In the **Source directory** path, provide the directory you want to compress.          
2. In the **Target directory path** input, you can select where to output the compressed file. 

### Troubleshooting

If the **Source directory** path does not exist, there will be an error.
If the folder structure belonging to the destination does not exist, you will be notified with a warning in the log and then the Step will create it.
If a ZIP exists on the specified destination, the user will be notified with a warning in the log and then the previous file will be overwritten.

### Related Steps

 - [Deploy to Bitrise.io](https://www.bitrise.io/integrations/steps/deploy-to-bitrise-io)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `source_path` | The path of the directory which you want to compress with zip. | required |  |
| `destination` | The path where you want to move the compressed file.  Can be a direcory or the archive path.  The `.zip` extension will be added automatically if it was omitted.  | required |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-create-zip/pulls) and [issues](https://github.com/bitrise-steplib/steps-create-zip/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
