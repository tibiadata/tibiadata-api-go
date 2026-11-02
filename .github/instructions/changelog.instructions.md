---
applyTo: "CHANGELOG.md"
---

# Changelog formatting

- Follow [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/): curate notable changes for humans, not a dump of commits or merged PR titles. Use [Semantic Versioning](https://semver.org/) for release versions.
- Preserve the existing `# Changelog` heading and historical release entries.
- Curate the changelog during release preparation, collecting notable changes since the verified previous release tag into one dated entry in the release PR. Do not require per-PR changelog updates or add an Unreleased section/reference. Unreleased is an optional Keep a Changelog practice that this repository does not use. List releases newest first.
- Use `## [X.Y.Z] - YYYY-MM-DD`: no `v` prefix in the heading. Use the agreed release date; ask if it is not known rather than inventing a publication date.
- Every new release heading must have exactly one matching reference definition in the comparison-link block at the bottom of the file. Add it above the existing definitions (newest first), preserving historical links:

  ```markdown
  [X.Y.Z]: https://github.com/tibiadata/tibiadata-api-go/compare/vPREVIOUS...vX.Y.Z
  ```

- Replace the placeholders with the verified previous release tag for the intended release line and the new release tag. Before finishing, confirm the heading's reference resolves to this definition and the comparison endpoints match the release's actual change range. A release entry without its bottom comparison link is incomplete.
- Use `### Added`, `### Changed`, `### Deprecated`, `### Removed`, `### Fixed`, and `### Security`, in that order, omitting empty categories. Classify by the actual change, not the commit prefix.
- Added covers new features; Changed covers changes in existing functionality; Deprecated covers features scheduled for removal; Removed covers removed features; Fixed covers bug fixes; Security covers vulnerability fixes. Explicitly describe breaking changes, deprecations, removals, and migration implications.
- Use one bullet per notable change, with a clear user-facing description and links to the relevant PRs and authors. Summarize related commits/PRs together when they describe one change; do not require Conventional Commit prefixes in changelog prose:

  ```markdown
  - Correct achievement grades in character responses ([#123](https://github.com/tibiadata/tibiadata-api-go/pull/123) by [username](https://github.com/username))
  ```

- Verify titles, PR numbers, and authors from GitHub. Do not invent entries or copy example placeholders into the changelog.
- The examples above describe the repository changelog. When turning an entry into GitHub release notes, follow the release skill: use plain `@username` mentions for authors, include verified new contributors when present, and end with an explicit Full Changelog comparison link. Keep historical changelog attribution unchanged.
- Use PR titles as evidence, not mandatory wording. Exclude routine merge/formatting noise, but include every important change. Do not retroactively reformat historical entries unless requested or correcting an omission/error.
- Consolidate repeated dependency bumps for the same dependency within the release into one bullet, using the earliest old version and final new version and linking every contributing PR in chronological order:

  ```markdown
  - deps(go): bump `golang.org/x/net` from 0.58.0 to 0.61.0 ([#728](https://github.com/tibiadata/tibiadata-api-go/pull/728), [#766](https://github.com/tibiadata/tibiadata-api-go/pull/766) by [dependabot](https://github.com/dependabot))
  ```

- Keep dependency updates in different modules/directories separate when their version ranges differ; retain the directory qualifier. Use backticks around Go module paths.
- Preserve grouped Dependabot update descriptions and combine repeated updates to the same group without losing PR references. If contributors differ, retain attribution for each contributing PR rather than assigning all changes to one author.
- In Changed, list non-dependency changes first, then dependency updates grouped by CI, Docker, and Go; order dependencies consistently by name within each group.
- For a major release, require a short introductory `###` section before the change categories (for example, `### The release of API v6`). Describe the key changes and link a published announcement if one exists; do not invent an announcement URL. Follow the release skill's additional major-release review.
- Distinguish API-version-specific behavior from improvements shared by supported API versions. Verify claims against code and merged PRs; a package release number is not itself an API compatibility change.
- Retain withdrawn releases in the history and mark their heading `## [X.Y.Z] - YYYY-MM-DD [YANKED]`; do not silently remove them.
