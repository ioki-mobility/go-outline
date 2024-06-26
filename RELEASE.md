# Release

We use the [GoReleaser(Action)](https://github.com/goreleaser/goreleaser-action) to perform a new release.
The GoReleaser configuration can be found at [.goreleaser.yml](.goreleaser.yml).

To trigger a new release, go to the Actions tab and select the [`Release` workflow](https://github.com/ioki-mobility/go-outline/actions/workflows/release.yml).
Click on the `Run workflow` drop-down menu and enter a meaningful `tag_name`.
We follow [Gos' version number convention](https://go.dev/doc/modules/version-numbers) for the `tag_name`.
So it should start with `v`, followed by a valid semver version.
Run the workflow from the `main` branch and select `Run workflow`.

That's it 🎉

If the workflow finished successfully, you should see a new release in the [Releases section](https://github.com/ioki-mobility/go-outline/releases).
