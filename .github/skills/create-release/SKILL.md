---
name: create-release
description: Prepare a TibiaData API Go release or publish it on request after the changelog PR is merged, including changelog formatting, author mentions, new contributors, comparison links, and release workflow verification.
---

# Create a release

## Authorization

- Preparation does not authorize creating a PR, pushing tags, publishing a release, or deploying. Publication is on request after the changelog PR merges, not automatic on merge.
- An explicit request to publish a particular release authorizes publication and its existing publication-triggered workflows. Before publishing, re-read `.github/workflows/build.yml`, `documentation.yml`, and `dockerhub.yml`: publication pushes images and dispatches production-related updates.
- Never move an existing release tag or edit an existing release without authorization.
- Only stable releases are supported for now. Do not create prereleases, use `--prerelease`, or publish prerelease-version tags such as `vX.Y.Z-rc.1`. The current `release: published` workflows also run for prereleases and dispatch production-targeted updates; prerelease support requires a separately approved workflow change first.

## Establish the release scope

1. Inspect the working tree and preserve unrelated changes. Confirm the requested stable version, release date, and intended target branch/commit. Ask about unresolved decisions; do not silently choose another version.
2. Inspect GitHub releases and tags to identify the previous release for the intended release line. Compare that tag with the intended target and collect all changes in the range. Do not assume that every recently merged PR belongs in this release.
3. Use the authenticated `gh` CLI to inspect merged PR titles, authors, merge commits, and dates. Paginate results and include direct commits without PRs when relevant, citing their commits rather than fabricating PR references.
4. Reconcile the commit range with the PR list so no notable change is omitted or counted twice. Use it as evidence for curated notes, not as a requirement to list every commit. Separate user-facing changes from dependency updates.

## Prepare the changelog and notes

1. Follow `.github/instructions/changelog.instructions.md` and its Keep a Changelog 1.1.0 baseline. Curate a new dated entry during release preparation from the verified change range, add its bottom comparison reference, and preserve older entries. Do not add an Unreleased section or require ongoing per-PR changelog updates.
2. Summarize significant changes accurately. For a major version, complete the additional major-release review below.
3. Include significant improvements shared by supported API versions as well as version-specific changes.
4. Prepare GitHub release notes using the format below.
5. Inspect directly related version references and documentation and update only those required by the release. The Docker build injects release metadata through build arguments; do not invent a source-code version bump.
6. Review the diff for complete change coverage, correct dependency ranges, and PR/author attribution. Confirm the new `[X.Y.Z]` heading has exactly one bottom comparison reference with the verified previous and new tags; preparation is incomplete without it.
7. Check CI for the intended release commit. If code changed during preparation, run relevant tests; the main module's race-test command is `go test -race ./...`. Account for separate Go modules when they changed, and report network/toolchain limitations. Changelog-only preparation does not require Go tests.
8. Present the version, intended target, notes, and blockers for review.

## Additional review for major releases

- Treat a major release as a coordinated publication, not just a larger change list. Confirm whether it introduces a new API version; do not infer runtime behavior from the package version alone.
- Require an introductory level-three section before the change categories in both `CHANGELOG.md` and the GitHub release body, such as `### The release of API v6`. Summarize the major highlights and compatibility implications. Link the published announcement when available; if it is not ready, surface that as a publication decision rather than inventing a URL.
- Review endpoint availability, response fields, canonical URLs/redirects, HTTP status codes, source/parsing changes, and migration requirements. Distinguish changes specific to the new API version from shared improvements. Past v5 highlights included Fansite API data, achievement grades, status codes, caching, logging, and telemetry; verify the new release's actual highlights rather than repeating that list.
- Confirm which API versions remain supported, the default version, and any deprecation decisions. Preparation must not silently expose a future version or deprecate an existing one.
- Identify readiness and required updates across application documentation, Swagger/OpenAPI, the announcement, and deployment/gateway routing. Report blockers and obtain decisions for unresolved migration or rollout behavior; do not create or publish announcements or change other repositories without authorization.
- Before publication, confirm the introductory text is approved and these readiness decisions are resolved. Minor and patch releases use the normal change categories without requiring this major-release introduction.

## GitHub release notes format

- Extract only the requested version's entry, stopping at the next release heading. Omit the changelog file heading, version/date heading, and unrelated versions. Preserve its introduction, categories, curated changes, and PR references. Prefer a short introduction and a change list over long prose.
- Use the level-three change categories and ordering from the changelog instructions, including Deprecated, Removed, and Security when applicable. The published [v5.0.0 release](https://github.com/tibiadata/tibiadata-api-go/releases/tag/v5.0.0) is a layout example, not a restriction to its three categories or Conventional Commit wording. Keep any introductory heading at level three; do not wrap categories in a generated "What's Changed" section.
- Use `-` for change bullets and compact same-repository PR references: `- User-facing description (#123 by @username)`. For consolidated updates, use `(#123, #124 by @dependabot)`. Convert explicit links to PRs in this repository to these auto-linked references; retain explicit URLs for cross-repository PRs.
- Convert author profile links to plain `@username` mentions, including bots and all authors of consolidated entries. For example, `by [username](https://github.com/username)` becomes `by @username`. Do not put mentions inside code spans or retain profile-link syntax in release notes. The repository changelog keeps its existing author-link format.
- Check for first-time contributors using GitHub's release-note generation endpoint, `POST repos/tibiadata/tibiadata-api-go/releases/generate-notes`, with `tag_name`, the confirmed `target_commitish`, and `previous_tag_name` set to the verified previous release tag. This generates preview text only; it does not publish a release. Use its New Contributors section as evidence, not its generated change list as a replacement for the curated changelog.
- Include New Contributors only when verified first-time contributors exist, using the heading and bullet format shown below. If generation is unavailable or inconclusive, verify against repository contribution history; do not silently omit the check or guess that someone is new.
- End every release body with the Full Changelog footer shown below, using the URL from the new version's bottom reference definition in `CHANGELOG.md`.
- Resolve any reference-style links copied from the entry into standalone links (or include their required definitions), so the published body has no dangling references.

Use this template, replacing every placeholder with verified release data. Include only applicable change categories and repeat the contributor bullet for each verified first-time contributor; omit New Contributors when there are none. Never publish placeholders:

```markdown
### Added

- <User-facing change description> (#<PR_NUMBER> by @<AUTHOR>)

## New Contributors
* @<NEW_CONTRIBUTOR> made their first contribution in #<FIRST_PR_NUMBER>

**Full Changelog**: https://github.com/tibiadata/tibiadata-api-go/compare/<PREVIOUS_TAG>...<NEW_TAG>
```

## Publish on request after the changelog PR is merged

1. Verify the changelog PR is merged into the intended release branch using GitHub's live PR state. Identify the resulting merge/squash/rebase commit and confirm the intended target includes the merged changelog. Do not tag an unmerged feature branch or assume the current worktree/default-branch tip is the intended release commit. If additional commits landed after preparation, resolve the intended target before publishing.
2. Read `CHANGELOG.md` from that exact remote target commit, not a local copy. Confirm the requested version, release date, bottom comparison reference, previous release tag, and CI status. For a publication-only request, do not regenerate or rewrite the merged changelog; report discrepancies before proceeding.
3. Check for an existing release and remote tag. If the release exists, report its URL rather than creating a duplicate. If a tag exists, resolve annotated tags to their commit and require it to match the confirmed target.
4. Build the release body using the format above and save it in session artifacts. Verify it contains only the requested version's changes, author mentions throughout, New Contributors when applicable, and exactly one correct Full Changelog footer.
5. Publish with the authenticated `gh release create`, explicitly supplying `--repo tibiadata/tibiadata-api-go`, the `vX.Y.Z` tag, `--title`, and `--notes-file` pointing to the curated body. If the remote tag does not exist, pass `--target` with the full confirmed commit SHA so GitHub creates the tag at that commit. If the tag already exists and matches, use `--verify-tag`. Check the release line before choosing latest status; use `--latest=false` for a maintenance release that must not replace the latest stable line. Do not use `--generate-notes`: publish the curated body, not GitHub's generated change list.
6. Verify the resulting release is published with the correct tag/target and complete notes, then perform the workflow verification below.

## Verify after publication

1. Confirm the workflows triggered on `release: published`. Do not use a manual build as a substitute; a manual build on main follows the dev path and does not dispatch the release-only Helm chart update.
2. Verify both linux/amd64 and linux/arm64 builds, the multi-architecture manifest, release-version tags in Docker Hub and GHCR, and the signing/attestation steps.
3. Verify the Argo CD dispatch includes the manifest `docker_digest`, `release_version`, and production `subdomain` (`api`). Verify the Helm chart dispatch receives `release_version`.
4. Verify the documentation workflow attaches `swagger.json` with the release tag as its version and dispatches the docs update. Check the Docker Hub description workflow as well.
5. Where access permits, check the resulting downstream runs/updates in `tibiadata/tibiadata-argocd-app-of-apps`, `tibiadata/tibiadata-helm-charts`, and `tibiadata/tibiadata-api-docs`. A successful dispatch alone is not proof of a successful deployment. Do not manually modify downstream repositories without authorization.
6. Report the actual GitHub release URL and distinguish publication, image builds, downstream dispatches, and confirmed deployment. State any failures or inaccessible verification explicitly.
