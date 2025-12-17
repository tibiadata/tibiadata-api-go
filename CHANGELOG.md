# Changelog

## [4.7.0] - 2025-12-17

### Added

* feat: add support for new highscore categories ([#563](https://github.com/tibiadata/tibiadata-api-go/pull/563) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

* chore(ci): add concurrency to build workflow ([#578](https://github.com/tibiadata/tibiadata-api-go/pull/578) by [tobiasehlert](https://github.com/tobiasehlert))
* deps(ci): bump actions/checkout from 5.0.0 to 6.0.1 ([#550](https://github.com/tibiadata/tibiadata-api-go/pull/550), [#555](https://github.com/tibiadata/tibiadata-api-go/pull/555), [#561](https://github.com/tibiadata/tibiadata-api-go/pull/561) by [dependabot](https://github.com/dependabot))
* deps(ci): bump actions/download-artifact from 6.0.0 to 7.0.0 ([#573](https://github.com/tibiadata/tibiadata-api-go/pull/573) by [dependabot](https://github.com/dependabot))
* deps(ci): bump actions/setup-go from 6.0.0 to 6.1.0 ([#557](https://github.com/tibiadata/tibiadata-api-go/pull/557) by [dependabot](https://github.com/dependabot))
* deps(ci): bump actions/upload-artifact from 5.0.0 to 6.0.0 ([#575](https://github.com/tibiadata/tibiadata-api-go/pull/575) by [dependabot](https://github.com/dependabot))
* deps(ci): bump codecov/codecov-action from 5.5.1 to 5.5.2 ([#571](https://github.com/tibiadata/tibiadata-api-go/pull/571) by [dependabot](https://github.com/dependabot))
* deps(ci): bump docker/metadata-action from 5.8.0 to 5.10.0 ([#543](https://github.com/tibiadata/tibiadata-api-go/pull/543), [#559](https://github.com/tibiadata/tibiadata-api-go/pull/559) by [dependabot](https://github.com/dependabot))
* deps(ci): bump docker/setup-qemu-action from 3.6.0 to 3.7.0 ([#544](https://github.com/tibiadata/tibiadata-api-go/pull/544) by [dependabot](https://github.com/dependabot))
* deps(ci): bump github/codeql-action from 4.31.0 to 4.31.8 ([#541](https://github.com/tibiadata/tibiadata-api-go/pull/541), [#548](https://github.com/tibiadata/tibiadata-api-go/pull/548), [#558](https://github.com/tibiadata/tibiadata-api-go/pull/558), [#560](https://github.com/tibiadata/tibiadata-api-go/pull/560), [#574](https://github.com/tibiadata/tibiadata-api-go/pull/574) by [dependabot](https://github.com/dependabot))
* deps(ci): bump peter-evans/repository-dispatch from 4.0.0 to 4.0.1 ([#549](https://github.com/tibiadata/tibiadata-api-go/pull/549) by [dependabot](https://github.com/dependabot))
* deps(ci): bump svenstaro/upload-release-action from 2.11.2 to 2.11.3 ([#556](https://github.com/tibiadata/tibiadata-api-go/pull/556) by [dependabot](https://github.com/dependabot))
* deps(ci): bump tailscale/github-action from 4.0.3 to 4.1.0 ([#542](https://github.com/tibiadata/tibiadata-api-go/pull/542) by [dependabot](https://github.com/dependabot))
* deps(docker): bump alpine from 3.22.2 to 3.23.0 ([#567](https://github.com/tibiadata/tibiadata-api-go/pull/567) by [dependabot](https://github.com/dependabot))
* deps(docker): bump golang from 1.25.3 to 1.25.5 ([#545](https://github.com/tibiadata/tibiadata-api-go/pull/545), [#566](https://github.com/tibiadata/tibiadata-api-go/pull/566) by [dependabot](https://github.com/dependabot))
* deps(go): bump github.com/go-resty/resty/v2 from 2.16.5 to 2.17.1 ([#553](https://github.com/tibiadata/tibiadata-api-go/pull/553), [#554](https://github.com/tibiadata/tibiadata-api-go/pull/554), [#576](https://github.com/tibiadata/tibiadata-api-go/pull/576), [#577](https://github.com/tibiadata/tibiadata-api-go/pull/577) by [dependabot](https://github.com/dependabot))
* deps(go): bump github.com/PuerkitoBio/goquery from 1.10.3 to 1.11.0 ([#551](https://github.com/tibiadata/tibiadata-api-go/pull/551) by [dependabot](https://github.com/dependabot))
* deps(go): bump github.com/quic-go/quic-go from 0.55.0 to 0.57.0 ([#572](https://github.com/tibiadata/tibiadata-api-go/pull/572) by [dependabot](https://github.com/dependabot))
* deps(go): bump golang.org/x/crypto from 0.43.0 to 0.45.0 ([#552](https://github.com/tibiadata/tibiadata-api-go/pull/552) by [dependabot](https://github.com/dependabot))
* deps(go): bump golang.org/x/net from 0.46.0 to 0.48.0 ([#547](https://github.com/tibiadata/tibiadata-api-go/pull/547), [#569](https://github.com/tibiadata/tibiadata-api-go/pull/569) by [dependabot](https://github.com/dependabot))
* deps(go): bump golang.org/x/text from 0.30.0 to 0.31.0 ([#546](https://github.com/tibiadata/tibiadata-api-go/pull/546) by [dependabot](https://github.com/dependabot))

### Fixed

* ci(fix): temporary removing cosign verification ([#564](https://github.com/tibiadata/tibiadata-api-go/pull/564) by [tobiasehlert](https://github.com/tobiasehlert))


## [4.6.0] - 2025-10-29

### Changed

* chore: rewrite for pure utf-8 content handling ([#511](https://github.com/tibiadata/tibiadata-api-go/pull/511) by [tobiasehlert](https://github.com/tobiasehlert))
* chore: refactor TibiaKillstatistics for better maintainability ([#513](https://github.com/tibiadata/tibiadata-api-go/pull/513) by [tobiasehlert](https://github.com/tobiasehlert))
* chore: use switches instead and silent TibiaDataStringToInteger error ([#514](https://github.com/tibiadata/tibiadata-api-go/pull/514) by [tobiasehlert](https://github.com/tobiasehlert))
* chore: use splitseq instead of split ([#515](https://github.com/tibiadata/tibiadata-api-go/pull/515) by [tobiasehlert](https://github.com/tobiasehlert))
* deps(ci): bump actions/download-artifact from 5.0.0 to 6.0.0 ([#536](https://github.com/tibiadata/tibiadata-api-go/pull/536) by [dependabot](https://github.com/dependabot))
* deps(ci): bump actions/upload-artifact from 4.6.2 to 5.0.0 ([#535](https://github.com/tibiadata/tibiadata-api-go/pull/535) by [dependabot](https://github.com/dependabot))
* deps(ci): bump docker/login-action from 3.5.0 to 3.6.0 ([#520](https://github.com/tibiadata/tibiadata-api-go/pull/520) by [dependabot](https://github.com/dependabot))
* deps(ci): bump github/codeql-action from 3.30.3 to 4.31.0 ([#516](https://github.com/tibiadata/tibiadata-api-go/pull/516), [#518](https://github.com/tibiadata/tibiadata-api-go/pull/518), [#524](https://github.com/tibiadata/tibiadata-api-go/pull/524), [#529](https://github.com/tibiadata/tibiadata-api-go/pull/529), [#533](https://github.com/tibiadata/tibiadata-api-go/pull/533), [#540](https://github.com/tibiadata/tibiadata-api-go/pull/540) by [dependabot](https://github.com/dependabot))
* deps(ci): bump peter-evans/dockerhub-description from 4.0.2 to 5.0.0 ([#522](https://github.com/tibiadata/tibiadata-api-go/pull/522) by [dependabot](https://github.com/dependabot))
* deps(ci): bump peter-evans/repository-dispatch from 3.0.0 to 4.0.0 ([#521](https://github.com/tibiadata/tibiadata-api-go/pull/521) by [dependabot](https://github.com/dependabot))
* deps(ci): bump sigstore/cosign-installer from 3.10.0 to 4.0.0 ([#534](https://github.com/tibiadata/tibiadata-api-go/pull/534) by [dependabot](https://github.com/dependabot))
* deps(ci): bump tailscale/github-action from 3.2.3 to 4.0.3 ([#517](https://github.com/tibiadata/tibiadata-api-go/pull/517), [#532](https://github.com/tibiadata/tibiadata-api-go/pull/532), [#537](https://github.com/tibiadata/tibiadata-api-go/pull/537) by [dependabot](https://github.com/dependabot))
* deps(docker): bump alpine from 3.22.1 to 3.22.2 ([#527](https://github.com/tibiadata/tibiadata-api-go/pull/527) by [dependabot](https://github.com/dependabot))
* deps(docker): bump golang from 1.25.1 to 1.25.3 ([#528](https://github.com/tibiadata/tibiadata-api-go/pull/528), [#539](https://github.com/tibiadata/tibiadata-api-go/pull/539) by [dependabot](https://github.com/dependabot))
* deps(go): bump github.com/gin-contrib/gzip from 1.2.3 to 1.2.5 ([#530](https://github.com/tibiadata/tibiadata-api-go/pull/530), [#538](https://github.com/tibiadata/tibiadata-api-go/pull/538) by [dependabot](https://github.com/dependabot))
* deps(go): bump github.com/quic-go/quic-go from 0.54.0 to 0.54.1 ([#526](https://github.com/tibiadata/tibiadata-api-go/pull/526) by [dependabot](https://github.com/dependabot))
* deps(go): bump golang.org/x/text from 0.29.0 to 0.30.0 ([#525](https://github.com/tibiadata/tibiadata-api-go/pull/525) by [dependabot](https://github.com/dependabot))

### Fixed

* fix(ci): specify slug for codecov ([#519](https://github.com/tibiadata/tibiadata-api-go/pull/519) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.5.0] - 2025-09-22

### Changed

- ci: use multiple runners for build ([#508](https://github.com/tibiadata/tibiadata-api-go/pull/508) by [tobiasehlert](https://github.com/tobiasehlert))
- ci: use sha instead of tag on steps ([#487](https://github.com/tibiadata/tibiadata-api-go/pull/487) by [tobiasehlert](https://github.com/tobiasehlert))
- lowercase GitHub organisation name ([#507](https://github.com/tibiadata/tibiadata-api-go/pull/507) by [tobiasehlert](https://github.com/tobiasehlert))
- bump actions/setup-go from 5.5.0 to 6.0.0 ([#495](https://github.com/tibiadata/tibiadata-api-go/pull/495) by [dependabot](https://github.com/dependabot))
- bump codecov/codecov-action from 5.4.3 to 5.5.1 ([#488](https://github.com/tibiadata/tibiadata-api-go/pull/488), [#496](https://github.com/tibiadata/tibiadata-api-go/pull/496) by [dependabot](https://github.com/dependabot))
- bump github/codeql-action from 3.29.10 to 3.30.3 ([#489](https://github.com/tibiadata/tibiadata-api-go/pull/489), [#494](https://github.com/tibiadata/tibiadata-api-go/pull/494), [#497](https://github.com/tibiadata/tibiadata-api-go/pull/497), [#500](https://github.com/tibiadata/tibiadata-api-go/pull/500), [#501](https://github.com/tibiadata/tibiadata-api-go/pull/501) by [dependabot](https://github.com/dependabot))
- bump sigstore/cosign-installer from 3.9.2 to 3.10.0 ([#503](https://github.com/tibiadata/tibiadata-api-go/pull/503) by [dependabot](https://github.com/dependabot))
- bump alpine from 3.21.3 to 3.22.1 ([#475](https://github.com/tibiadata/tibiadata-api-go/pull/475), [#482](https://github.com/tibiadata/tibiadata-api-go/pull/482) by [dependabot](https://github.com/dependabot))
- bump golang from 1.24.2 to 1.25.1 ([#477](https://github.com/tibiadata/tibiadata-api-go/pull/477), [#481](https://github.com/tibiadata/tibiadata-api-go/pull/481), [#484](https://github.com/tibiadata/tibiadata-api-go/pull/484), [#486](https://github.com/tibiadata/tibiadata-api-go/pull/486), [#502](https://github.com/tibiadata/tibiadata-api-go/pull/502) by [dependabot](https://github.com/dependabot))
- bump github.com/gin-gonic/gin from 1.10.0 to 1.11.0 ([#474](https://github.com/tibiadata/tibiadata-api-go/pull/474), [#510](https://github.com/tibiadata/tibiadata-api-go/pull/510) by [dependabot](https://github.com/dependabot))
- bump github.com/PuerkitoBio/goquery from 1.10.2 to 1.10.3 ([#468](https://github.com/tibiadata/tibiadata-api-go/pull/468) by [dependabot](https://github.com/dependabot))
- bump github.com/stretchr/testify from 1.10.0 to 1.11.1 ([#492](https://github.com/tibiadata/tibiadata-api-go/pull/492), [#493](https://github.com/tibiadata/tibiadata-api-go/pull/493) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/net from 0.36.0 to 0.38.0 ([#469](https://github.com/tibiadata/tibiadata-api-go/pull/469) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/text from 0.24.0 to 0.29.0 ([#472](https://github.com/tibiadata/tibiadata-api-go/pull/472), [#476](https://github.com/tibiadata/tibiadata-api-go/pull/476), [#480](https://github.com/tibiadata/tibiadata-api-go/pull/480), [#483](https://github.com/tibiadata/tibiadata-api-go/pull/483), [#498](https://github.com/tibiadata/tibiadata-api-go/pull/498) by [dependabot](https://github.com/dependabot))

### Fixed

- ci: use dockerhub variable instead ([#509](https://github.com/tibiadata/tibiadata-api-go/pull/509) by [tobiasehlert](https://github.com/tobiasehlert))
- character death entry with lowercased of-creature ([#505](https://github.com/tibiadata/tibiadata-api-go/pull/505) by [tobiasehlert](https://github.com/tobiasehlert))
- remove iso-encoding for umlaut issue ([#506](https://github.com/tibiadata/tibiadata-api-go/pull/506) by [tobiasehlert](https://github.com/tobiasehlert))
- use 555 for binary ([#499](https://github.com/tibiadata/tibiadata-api-go/pull/499) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.4.0] - 2025-04-08

### Added

- support new vocation monk ([#451](https://github.com/tibiadata/tibiadata-api-go/pull/451) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- bump github.com/gin-contrib/gzip from 1.2.2 to 1.2.3 ([#466](https://github.com/tibiadata/tibiadata-api-go/pull/466) by [dependabot](https://github.com/dependabot))
- bump golang from 1.24.1 to 1.24.2 ([#465](https://github.com/tibiadata/tibiadata-api-go/pull/465) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/text from 0.23.0 to 0.24.0 ([#467](https://github.com/tibiadata/tibiadata-api-go/pull/467) by [dependabot](https://github.com/dependabot))

### Fixed

- omit character deaths_truncated if false ([#464](https://github.com/tibiadata/tibiadata-api-go/pull/464) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.6] - 2025-04-01

### Changed

- bump golang from 1.24.0 to 1.24.1 ([#459](https://github.com/tibiadata/tibiadata-api-go/pull/459) by [dependabot](https://github.com/dependabot))
- bump golang.org/x/net from 0.35.0 to 0.36.0 ([#460](https://github.com/tibiadata/tibiadata-api-go/pull/460), [#461](https://github.com/tibiadata/tibiadata-api-go/pull/461) by [dependabot](https://github.com/dependabot))

### Fixed

- character death html header change ([#463](https://github.com/tibiadata/tibiadata-api-go/pull/463) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.5] - 2025-03-09

### Changed

- bump golang.org/x/text from 0.22.0 to 0.23.0 ([#457](https://github.com/tibiadata/tibiadata-api-go/pull/457) by [dependabot](https://github.com/dependabot))

### Fixed

- character-death parsing when only being assisted ([#454](https://github.com/tibiadata/tibiadata-api-go/pull/454) by [tobiasehlert](https://github.com/tobiasehlert))
- improve character death parsing and removing creature prefix ([#456](https://github.com/tibiadata/tibiadata-api-go/pull/456) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.4] - 2025-02-26

### Changed

- Bump golang.org/x/text from 0.21.0 to 0.22.0 ([#441](https://github.com/tibiadata/tibiadata-api-go/pull/441) by [dependabot](https://github.com/dependabot))
- Bump github.com/PuerkitoBio/goquery from 1.10.1 to 1.10.2 ([#443](https://github.com/tibiadata/tibiadata-api-go/pull/443) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.23.5 to 1.24.0 ([#444](https://github.com/tibiadata/tibiadata-api-go/pull/444) by [dependabot](https://github.com/dependabot))
- Bump alpine from 3.21.2 to 3.21.3 ([#445](https://github.com/tibiadata/tibiadata-api-go/pull/445) by [dependabot](https://github.com/dependabot))

### Fixed

- Character name can contain comma ([#447](https://github.com/tibiadata/tibiadata-api-go/pull/447) by [tobiasehlert](https://github.com/tobiasehlert))
- url_api field format on news endpoints ([#452](https://github.com/tibiadata/tibiadata-api-go/pull/452) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.3] - 2025-02-04

### Changed

- Improving validation-tests ([#434](https://github.com/tibiadata/tibiadata-api-go/pull/434) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump alpine from 3.21.0 to 3.21.2 ([#429](https://github.com/tibiadata/tibiadata-api-go/pull/429) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-contrib/gzip from 1.1.0 to 1.2.2 ([#430](https://github.com/tibiadata/tibiadata-api-go/pull/430), [#437](https://github.com/tibiadata/tibiadata-api-go/pull/437) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.16.2 to 2.16.5 ([#427](https://github.com/tibiadata/tibiadata-api-go/pull/427), [#428](https://github.com/tibiadata/tibiadata-api-go/pull/428), [#432](https://github.com/tibiadata/tibiadata-api-go/pull/432), [#433](https://github.com/tibiadata/tibiadata-api-go/pull/433), [#435](https://github.com/tibiadata/tibiadata-api-go/pull/435), [#436](https://github.com/tibiadata/tibiadata-api-go/pull/436) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.23.4 to 1.23.5 ([#431](https://github.com/tibiadata/tibiadata-api-go/pull/431) by [dependabot](https://github.com/dependabot))

### Fixed

- Parsing of boosted creature that is fictitious ([#440](https://github.com/tibiadata/tibiadata-api-go/pull/440) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.2] - 2024-12-29

### Changed

- Bump tailscale/github-action from 2 to 3 ([#422](https://github.com/tibiadata/tibiadata-api-go/pull/422) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-contrib/gzip from 1.0.1 to 1.1.0 ([#423](https://github.com/tibiadata/tibiadata-api-go/pull/423) by [dependabot](https://github.com/dependabot))
- Bump github.com/PuerkitoBio/goquery from 1.10.0 to 1.10.1 ([#425](https://github.com/tibiadata/tibiadata-api-go/pull/425) by [dependabot](https://github.com/dependabot))

### Fixed

- Change guild rank length check ([#426](https://github.com/tibiadata/tibiadata-api-go/pull/426) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.3.1] - 2024-12-12

### Added

- Add root-endpoint to container ([#420](https://github.com/tibiadata/tibiadata-api-go/pull/420) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 ([#415](https://github.com/tibiadata/tibiadata-api-go/pull/415), [#416](https://github.com/tibiadata/tibiadata-api-go/pull/416) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.16.1 to 2.16.2 ([#413](https://github.com/tibiadata/tibiadata-api-go/pull/413), [#414](https://github.com/tibiadata/tibiadata-api-go/pull/414) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.20.0 to 0.21.0 ([#417](https://github.com/tibiadata/tibiadata-api-go/pull/417) by [dependabot](https://github.com/dependabot))
- Bump alpine from 3.20.3 to 3.21.0 ([#419](https://github.com/tibiadata/tibiadata-api-go/pull/419) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.23.3 to 1.23.4 ([#418](https://github.com/tibiadata/tibiadata-api-go/pull/418) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/crypto from 0.27.0 to 0.31.0 ([#421](https://github.com/tibiadata/tibiadata-api-go/pull/421) by [dependabot](https://github.com/dependabot))

## [4.3.0] - 2024-11-22

### Added

- Add tibia_urls field into information-section ([#398](https://github.com/tibiadata/tibiadata-api-go/pull/398) by [Kai-Animator](https://github.com/Kai-Animator))

### Changed

- Bump codecov/codecov-action from 4 to 5 ([#406](https://github.com/tibiadata/tibiadata-api-go/pull/406) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.15.3 to 2.16.1 ([#404](https://github.com/tibiadata/tibiadata-api-go/pull/404), [#405](https://github.com/tibiadata/tibiadata-api-go/pull/405), [#407](https://github.com/tibiadata/tibiadata-api-go/pull/407), [#408](https://github.com/tibiadata/tibiadata-api-go/pull/408) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.23.2 to 1.23.3 ([#403](https://github.com/tibiadata/tibiadata-api-go/pull/403) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.19.0 to 0.20.0 ([#401](https://github.com/tibiadata/tibiadata-api-go/pull/401) by [dependabot](https://github.com/dependabot))
- Update dependabot configuration ([#399](https://github.com/tibiadata/tibiadata-api-go/pull/399) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- Empty boostablebosses list ([#411](https://github.com/tibiadata/tibiadata-api-go/pull/411) by [tobiasehlert](https://github.com/tobiasehlert))
- Failing creature endpoints ([#412](https://github.com/tibiadata/tibiadata-api-go/pull/412) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.2.3] - 2024-10-15

### Changed

- Bump alpine from 3.20.1 to 3.20.3 ([#375](https://github.com/tibiadata/tibiadata-api-go/pull/375), [#386](https://github.com/tibiadata/tibiadata-api-go/pull/386) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.13.1 to 2.15.3 ([#376](https://github.com/tibiadata/tibiadata-api-go/pull/376), [#377](https://github.com/tibiadata/tibiadata-api-go/pull/377), [#388](https://github.com/tibiadata/tibiadata-api-go/pull/388), [#389](https://github.com/tibiadata/tibiadata-api-go/pull/389), [#390](https://github.com/tibiadata/tibiadata-api-go/pull/390), [#391](https://github.com/tibiadata/tibiadata-api-go/pull/391), [#393](https://github.com/tibiadata/tibiadata-api-go/pull/393), [#392](https://github.com/tibiadata/tibiadata-api-go/pull/392), [#394](https://github.com/tibiadata/tibiadata-api-go/pull/394), [#395](https://github.com/tibiadata/tibiadata-api-go/pull/395) by [dependabot](https://github.com/dependabot))
- Bump github.com/PuerkitoBio/goquery from 1.9.2 to 1.10.0 ([#384](https://github.com/tibiadata/tibiadata-api-go/pull/384) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.22.5 to 1.23.2 ([#379](https://github.com/tibiadata/tibiadata-api-go/pull/379), [#382](https://github.com/tibiadata/tibiadata-api-go/pull/382), [#385](https://github.com/tibiadata/tibiadata-api-go/pull/385), [#396](https://github.com/tibiadata/tibiadata-api-go/pull/396) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.16.0 to 0.19.0 ([#378](https://github.com/tibiadata/tibiadata-api-go/pull/378), [#383](https://github.com/tibiadata/tibiadata-api-go/pull/383), [#397](https://github.com/tibiadata/tibiadata-api-go/pull/397) by [dependabot](https://github.com/dependabot))
- ci: remove GitHub action Tailscale version input ([#387](https://github.com/tibiadata/tibiadata-api-go/pull/387) by [tobiasehlert](https://github.com/tobiasehlert))
- chore: decrease Dockerfile layers ([#380](https://github.com/tibiadata/tibiadata-api-go/pull/380) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.2.2] - 2024-07-24

### Fixed

- Boosted creature parsing on overview page ([#374](https://github.com/tibiadata/tibiadata-api-go/pull/374) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.2.1] - 2024-07-22

### Fixed

- Boosted creature parsing on overview page ([#372](https://github.com/tibiadata/tibiadata-api-go/pull/372) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.2.0] - 2024-07-11

### Added

- Add restriction mode ([#370](https://github.com/tibiadata/tibiadata-api-go/pull/370) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- ci: rewrite GitHub action workflows ([#346](https://github.com/tibiadata/tibiadata-api-go/pull/346) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump alpine from 3.19.1 to 3.20.1 ([#363](https://github.com/tibiadata/tibiadata-api-go/pull/363), [#367](https://github.com/tibiadata/tibiadata-api-go/pull/367) by [dependabot](https://github.com/dependabot))
- Bump docker/build-push-action from 5 to 6 ([#366](https://github.com/tibiadata/tibiadata-api-go/pull/366) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-contrib/gzip from 1.0.0 to 1.0.1 ([#355](https://github.com/tibiadata/tibiadata-api-go/pull/355) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-gonic/gin from 1.9.1 to 1.10.0 ([#357](https://github.com/tibiadata/tibiadata-api-go/pull/357) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.12.0 to 2.13.1 ([#358](https://github.com/tibiadata/tibiadata-api-go/pull/358), [#359](https://github.com/tibiadata/tibiadata-api-go/pull/359), [#362](https://github.com/tibiadata/tibiadata-api-go/pull/362), [#360](https://github.com/tibiadata/tibiadata-api-go/pull/360) by [dependabot](https://github.com/dependabot))
- Bump github.com/PuerkitoBio/goquery from 1.9.1 to 1.9.2 ([#354](https://github.com/tibiadata/tibiadata-api-go/pull/354) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.22.1 to 1.22.5 ([#342](https://github.com/tibiadata/tibiadata-api-go/pull/342), [#361](https://github.com/tibiadata/tibiadata-api-go/pull/361), [#365](https://github.com/tibiadata/tibiadata-api-go/pull/365), [#368](https://github.com/tibiadata/tibiadata-api-go/pull/368) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/net from 0.22.0 to 0.23.0 ([#353](https://github.com/tibiadata/tibiadata-api-go/pull/353), [#351](https://github.com/tibiadata/tibiadata-api-go/pull/351), [#352](https://github.com/tibiadata/tibiadata-api-go/pull/352) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.14.0 to 0.16.0 ([#356](https://github.com/tibiadata/tibiadata-api-go/pull/356), [#364](https://github.com/tibiadata/tibiadata-api-go/pull/364) by [dependabot](https://github.com/dependabot))

### Fixed

- Align world endpoint transfer_type to regular by default ([#348](https://github.com/tibiadata/tibiadata-api-go/pull/348) by [tobiasehlert](https://github.com/tobiasehlert))
- New world record_date field omit if empty ([#350](https://github.com/tibiadata/tibiadata-api-go/pull/350) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.1.2] - 2024-04-02

### Changed

- Bump github.com/go-resty/resty/v2 from 2.11.0 to 2.12.0 ([#334](https://github.com/tibiadata/tibiadata-api-go/pull/334), [#335](https://github.com/tibiadata/tibiadata-api-go/pull/335) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-contrib/gzip from 0.0.6 to 1.0.0 ([#336](https://github.com/tibiadata/tibiadata-api-go/pull/336) by [dependabot](https://github.com/dependabot))
- ci: switch from custom proxy to tailscale vpn ([#339](https://github.com/tibiadata/tibiadata-api-go/pull/339) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- Character killer parser by rewrite of death parsing logic ([#338](https://github.com/tibiadata/tibiadata-api-go/pull/338) by [phenpessoa](https://github.com/phenpessoa))

## [4.1.1] - 2024-03-15

### Fixed

- Fix permission issues with Dockerfile nonroot implementation ([#333](https://github.com/tibiadata/tibiadata-api-go/pull/333) by [tobiasehlert](https://github.com/tobiasehlert))

## [4.1.0] - 2024-03-14

### Changed

- Update Dockerfile to specific version and use of nonroot user ([#322](https://github.com/tibiadata/tibiadata-api-go/pull/322) by [tobiasehlert](https://github.com/tobiasehlert))
- Update go mods ([#332](https://github.com/tibiadata/tibiadata-api-go/pull/332) by [tobiasehlert](https://github.com/tobiasehlert))
- Using constant instead of reusing similar strings ([#321](https://github.com/tibiadata/tibiadata-api-go/pull/321) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump actions/cache from 3 to 4 ([#313](https://github.com/tibiadata/tibiadata-api-go/pull/313) by [dependabot](https://github.com/dependabot))
- Bump codecov/codecov-action from 3 to 4 ([#317](https://github.com/tibiadata/tibiadata-api-go/pull/317) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.10.0 to 2.11.0 ([#308](https://github.com/tibiadata/tibiadata-api-go/pull/308), [#307](https://github.com/tibiadata/tibiadata-api-go/pull/307), [#310](https://github.com/tibiadata/tibiadata-api-go/pull/310) by [dependabot](https://github.com/dependabot))
- Bump github.com/PuerkitoBio/goquery from 1.8.1 to 1.9.1 ([#326](https://github.com/tibiadata/tibiadata-api-go/pull/326) by [dependabot](https://github.com/dependabot))
- Bump github.com/stretchr/testify from 1.8.4 to 1.9.0 ([#327](https://github.com/tibiadata/tibiadata-api-go/pull/327), [#328](https://github.com/tibiadata/tibiadata-api-go/pull/328) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.21.5 to 1.22.1 ([#312](https://github.com/tibiadata/tibiadata-api-go/pull/312), [#324](https://github.com/tibiadata/tibiadata-api-go/pull/324), [#330](https://github.com/tibiadata/tibiadata-api-go/pull/330) by [dependabot](https://github.com/dependabot))
- Bump google.golang.org/protobuf from 1.32.0 to 1.33.0 ([#331](https://github.com/tibiadata/tibiadata-api-go/pull/331) by [dependabot](https://github.com/dependabot))
- Bump peter-evans/dockerhub-description from 3 to 4 ([#316](https://github.com/tibiadata/tibiadata-api-go/pull/316) by [dependabot](https://github.com/dependabot))
- Bump peter-evans/repository-dispatch from 2 to 3 ([#315](https://github.com/tibiadata/tibiadata-api-go/pull/315) by [dependabot](https://github.com/dependabot))
- Bump reproducible-containers/buildkit-cache-dance from 2.1.3 to 2.1.4 ([#314](https://github.com/tibiadata/tibiadata-api-go/pull/314) by [dependabot](https://github.com/dependabot))

### Fixed

- Never logged in parsing on characters endpoint ([#320](https://github.com/tibiadata/tibiadata-api-go/pull/320) by [tobiasehlert](https://github.com/tobiasehlert))
- Properly parse players with 'and' in their names ([#323](https://github.com/tibiadata/tibiadata-api-go/pull/323) by [phenpessoa](https://github.com/phenpessoa))

## [4.0.2] - 2023-12-25

### Changed

- Bump github/codeql-action from 2 to 3 ([#302](https://github.com/tibiadata/tibiadata-api-go/pull/302) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/crypto from 0.14.0 to 0.17.0 ([#304](https://github.com/tibiadata/tibiadata-api-go/pull/304) by [dependabot](https://github.com/dependabot))

### Fixed

- Validation: consider +, - and ' as a valid separators for character name words ([#306](https://github.com/tibiadata/tibiadata-api-go/pull/306) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

## [4.0.1] - 2023-12-12

### Changed

- Bump golang from 1.21.4 to 1.21.5 ([#299](https://github.com/tibiadata/tibiadata-api-go/pull/299) by [dependabot](https://github.com/dependabot))
- Bump actions/setup-go from 4 to 5 ([#298](https://github.com/tibiadata/tibiadata-api-go/pull/298) by [dependabot](https://github.com/dependabot))

### Fixed

- Webserver/validation: improve error handling on non OK status codes ([#301](https://github.com/tibiadata/tibiadata-api-go/pull/301) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

## [4.0.0] - 2023-12-01

### The release of API v4

Larger rewrite of application with implementation of new error handling, changed data structures and more.

Head over to [tibiadata.com](https://tibiadata.com/2023/12/tibiadata-api-v4-released/) for more information.

### Added

- Add error handling and some more ([#119](https://github.com/tibiadata/tibiadata-api-go/pull/119) by [tobiasehlert](https://github.com/tobiasehlert))
- Add new error ErrorMaintenanceMode ([#186](https://github.com/tibiadata/tibiadata-api-go/pull/186) by [tobiasehlert](https://github.com/tobiasehlert))
- Add parsing of healed by damage for creatures ([#182](https://github.com/tibiadata/tibiadata-api-go/pull/182) by [tobiasehlert](https://github.com/tobiasehlert))
- Add possibility to set trusted proxies in gin by env ([#200](https://github.com/tibiadata/tibiadata-api-go/pull/200) by [tobiasehlert](https://github.com/tobiasehlert))
- Add test for guild at war ([#266](https://github.com/tibiadata/tibiadata-api-go/pull/266) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Extend special creatures list for death parsing ([#184](https://github.com/tibiadata/tibiadata-api-go/pull/184) by [tobiasehlert](https://github.com/tobiasehlert))
- Extending test for tibiadata utils ([#267](https://github.com/tibiadata/tibiadata-api-go/pull/267) by [tobiasehlert](https://github.com/tobiasehlert))
- Reduce heap allocations ([#221](https://github.com/tibiadata/tibiadata-api-go/pull/221) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Replacing strings.Title with cases.Title ([#183](https://github.com/tibiadata/tibiadata-api-go/pull/183) by [tobiasehlert](https://github.com/tobiasehlert))
- Rewrite boostable bosses parser ([#233](https://github.com/tibiadata/tibiadata-api-go/pull/233) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Rewrite TibiaDataDate function ([#240](https://github.com/tibiadata/tibiadata-api-go/pull/240) by [tobiasehlert](https://github.com/tobiasehlert))
- Rewrite TibiaCharactersCharacter by removing various regex ([#228](https://github.com/tibiadata/tibiadata-api-go/pull/228), [#227](https://github.com/tibiadata/tibiadata-api-go/pull/227), [#226](https://github.com/tibiadata/tibiadata-api-go/pull/226), [#225](https://github.com/tibiadata/tibiadata-api-go/pull/225), [#224](https://github.com/tibiadata/tibiadata-api-go/pull/224), [#223](https://github.com/tibiadata/tibiadata-api-go/pull/223), [#229](https://github.com/tibiadata/tibiadata-api-go/pull/229) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Update boostable bosses parsers to latest tibia's html ([#261](https://github.com/tibiadata/tibiadata-api-go/pull/261) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Update of isEnvExist to return false if empty ([#205](https://github.com/tibiadata/tibiadata-api-go/pull/205) by [tobiasehlert](https://github.com/tobiasehlert))
- Update of test for webserver and main ([#269](https://github.com/tibiadata/tibiadata-api-go/pull/269) by [tobiasehlert](https://github.com/tobiasehlert))
- Workflow setup go version by go.mod ([#281](https://github.com/tibiadata/tibiadata-api-go/pull/281) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump docker/build-push-action from 3 to 5 ([#185](https://github.com/tibiadata/tibiadata-api-go/pull/185), [#250](https://github.com/tibiadata/tibiadata-api-go/pull/250) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.4.0 to 0.6.0 ([#170](https://github.com/tibiadata/tibiadata-api-go/pull/170), ([#173](https://github.com/tibiadata/tibiadata-api-go/pull/173) by [dependabot](https://github.com/dependabot))
- Bump actions/checkout from 3 to 4 ([#248](https://github.com/tibiadata/tibiadata-api-go/pull/248) by [dependabot](https://github.com/dependabot))
- Bump actions/setup-go from 3 to 4 ([#204](https://github.com/tibiadata/tibiadata-api-go/pull/204) by [dependabot](https://github.com/dependabot))
- Bump docker/login-action from 2 to 3 ([#251](https://github.com/tibiadata/tibiadata-api-go/pull/251) by [dependabot](https://github.com/dependabot))
- Bump docker/metadata-action from 4 to 5 ([#252](https://github.com/tibiadata/tibiadata-api-go/pull/252) by [dependabot](https://github.com/dependabot))
- Bump docker/setup-buildx-action from 2 to 3 ([#253](https://github.com/tibiadata/tibiadata-api-go/pull/253) by [dependabot](https://github.com/dependabot))
- Bump docker/setup-qemu-action from 2 to 3 ([#254](https://github.com/tibiadata/tibiadata-api-go/pull/254) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-gonic/gin from 1.8.2 to 1.9.1 ([#197](https://github.com/tibiadata/tibiadata-api-go/pull/197), [#214](https://github.com/tibiadata/tibiadata-api-go/pull/214) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.7.0 to 2.10.0 ([#256](https://github.com/tibiadata/tibiadata-api-go/pull/256), [#260](https://github.com/tibiadata/tibiadata-api-go/pull/260), [#282](https://github.com/tibiadata/tibiadata-api-go/pull/282) by [dependabot](https://github.com/dependabot))
- Bump github.com/go-resty/resty/v2 from 2.7.0 to 2.10.0 in /src/tibiamapping ([#255](https://github.com/tibiadata/tibiadata-api-go/pull/255), [#259](https://github.com/tibiadata/tibiadata-api-go/pull/259), [#284](https://github.com/tibiadata/tibiadata-api-go/pull/284) by [dependabot](https://github.com/dependabot))
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.4 ([#199](https://github.com/tibiadata/tibiadata-api-go/pull/199), [#211](https://github.com/tibiadata/tibiadata-api-go/pull/211), [#213](https://github.com/tibiadata/tibiadata-api-go/pull/213) by [dependabot](https://github.com/dependabot))
- Bump github.com/stretchr/testify from 1.8.1 to 1.8.4 in /src/validation ([#247](https://github.com/tibiadata/tibiadata-api-go/pull/247) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.19.5 to 1.21.4 ([#187](https://github.com/tibiadata/tibiadata-api-go/pull/187), [#196](https://github.com/tibiadata/tibiadata-api-go/pull/196), [#202](https://github.com/tibiadata/tibiadata-api-go/pull/202), [#207](https://github.com/tibiadata/tibiadata-api-go/pull/207), [#208](https://github.com/tibiadata/tibiadata-api-go/pull/208), [#216](https://github.com/tibiadata/tibiadata-api-go/pull/216), [#238](https://github.com/tibiadata/tibiadata-api-go/pull/238), [#243](https://github.com/tibiadata/tibiadata-api-go/pull/243), [#245](https://github.com/tibiadata/tibiadata-api-go/pull/245), [#249](https://github.com/tibiadata/tibiadata-api-go/pull/249), [#283](https://github.com/tibiadata/tibiadata-api-go/pull/283), [#293](https://github.com/tibiadata/tibiadata-api-go/pull/293) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.7.0 to 0.14.0 ([#201](https://github.com/tibiadata/tibiadata-api-go/pull/201), [#206](https://github.com/tibiadata/tibiadata-api-go/pull/206), [#218](https://github.com/tibiadata/tibiadata-api-go/pull/218), [#237](https://github.com/tibiadata/tibiadata-api-go/pull/237), [#242](https://github.com/tibiadata/tibiadata-api-go/pull/242), [#246](https://github.com/tibiadata/tibiadata-api-go/pull/246), [#289](https://github.com/tibiadata/tibiadata-api-go/pull/289) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/net from 0.16.0 to 0.17.0 ([#279](https://github.com/tibiadata/tibiadata-api-go/pull/279) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/net from 0.16.0 to 0.17.0 in /src/tibiamapping ([#278](https://github.com/tibiadata/tibiadata-api-go/pull/278) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/net from 0.16.0 to 0.17.0 in /src/validation ([#280](https://github.com/tibiadata/tibiadata-api-go/pull/280) by [dependabot](https://github.com/dependabot))
- Bump reproducible-containers/buildkit-cache-dance from 2.1.2 to 2.1.3 ([#286](https://github.com/tibiadata/tibiadata-api-go/pull/286) by [dependabot](https://github.com/dependabot))

### Fixed

- Fix character deaths where timestamp contains CEST ([#265](https://github.com/tibiadata/tibiadata-api-go/pull/265) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix character name max length ([#194](https://github.com/tibiadata/tibiadata-api-go/pull/194) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix character title that contains Grade ([#264](https://github.com/tibiadata/tibiadata-api-go/pull/264) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix characters with too many deaths cause crash ([#272](https://github.com/tibiadata/tibiadata-api-go/pull/272) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix highscores pages bug ([#210](https://github.com/tibiadata/tibiadata-api-go/pull/210) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix max length of character name ([#275](https://github.com/tibiadata/tibiadata-api-go/pull/275) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix not parsing all news ticker properly ([#291](https://github.com/tibiadata/tibiadata-api-go/pull/291) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix shorter character names ([#288](https://github.com/tibiadata/tibiadata-api-go/pull/288) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix test for tibiaHighscores with page param ([#235](https://github.com/tibiadata/tibiadata-api-go/pull/235) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix TibiaHousesOverview: format Ab'Dendriel properly ([#297](https://github.com/tibiadata/tibiadata-api-go/pull/297) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix truncated deaths ([#276](https://github.com/tibiadata/tibiadata-api-go/pull/276) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix validation: replace '+' with ' ' in TownExists ([#295](https://github.com/tibiadata/tibiadata-api-go/pull/295) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix webserver: pass normalized creature race to TibiaCreaturesCreatureIm… ([#236](https://github.com/tibiadata/tibiadata-api-go/pull/236) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Fix webserver: return on error and check if location is not nil ([#257](https://github.com/tibiadata/tibiadata-api-go/pull/257) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

## [3.7.5] - 2023-11-10

### Fixed

- Fix not parsing all news ticker properly ([#292](https://github.com/tibiadata/tibiadata-api-go/pull/292) by [phenpessoa](https://github.com/phenpessoa))

## [3.7.4] - 2023-09-29

### Fixed

- Fix invalid json v3 ([#258](https://github.com/tibiadata/tibiadata-api-go/pull/258) by [phenpessoa](https://github.com/phenpessoa))

## [3.7.3] - 2023-08-01

### Fixed

- Rewrite TibiaDataDateV3 function ([#241](https://github.com/tibiadata/tibiadata-api-go/pull/241) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.7.2] - 2023-05-22

### Fixed

- Removing highscore limit of 23 pages ([#210](https://github.com/tibiadata/tibiadata-api-go/pull/210) by [phenpessoa](https://github.com/phenpessoa))

## [3.7.1] - 2023-03-06

### Fixed

- Removing unintended logging of trustedProxies var
- Adding back latest tagging of builds

## [3.7.0] - 2023-03-06

### Added

- Add parsing of healed by damage for creatures ([#182](https://github.com/tibiadata/tibiadata-api-go/pull/182) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Bump golang from 1.19.5 to 1.20.1
- Extend special creatures list for death parsing ([#184](https://github.com/tibiadata/tibiadata-api-go/pull/184) by [tobiasehlert](https://github.com/tobiasehlert))
- Updating workflows and go mods

## [3.6.0] - 2023-01-25

### Added

- Feature: adding support for k8s health endpoints ([#147](https://github.com/tibiadata/tibiadata-api-go/pull/147) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Switching to assert.Empty from assert.Equal empty ([#181](https://github.com/tibiadata/tibiadata-api-go/pull/181) by [tobiasehlert](https://github.com/tobiasehlert))
- Adjusting highscore span to support page up to 23 ([#168](https://github.com/tibiadata/tibiadata-api-go/pull/168) by [tobiasehlert](https://github.com/tobiasehlert))
- Add more spell tests ([#179](https://github.com/tibiadata/tibiadata-api-go/pull/179) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump github.com/stretchr/testify from 1.8.0 to 1.8.1 ([#166](https://github.com/tibiadata/tibiadata-api-go/pull/166) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.4.0 to 0.6.0 ([#170](https://github.com/tibiadata/tibiadata-api-go/pull/170), ([#173](https://github.com/tibiadata/tibiadata-api-go/pull/173) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.19.2 to 1.19.5 ([#169](https://github.com/tibiadata/tibiadata-api-go/pull/169), [#171](https://github.com/tibiadata/tibiadata-api-go/pull/171), [#174](https://github.com/tibiadata/tibiadata-api-go/pull/174) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-gonic/gin from 1.8.1 to 1.8.2 ([#172](https://github.com/tibiadata/tibiadata-api-go/pull/172) by [dependabot](https://github.com/dependabot))

### Fixed

- Fix issue parsing summon mana cost ([#175](https://github.com/tibiadata/tibiadata-api-go/pull/175) by [soul4soul](https://github.com/soul4soul))
- Fix account badge parsing ([#178](https://github.com/tibiadata/tibiadata-api-go/pull/178) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix guild date format ([#180](https://github.com/tibiadata/tibiadata-api-go/pull/180) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.5.1] - 2022-10-19

### Fixed

- Automated documentation with swaggo failed due to formatting

## [3.5.0] - 2022-10-19

### Added

- Fix highscores by implementing pagination ([#164](https://github.com/tibiadata/tibiadata-api-go/pull/164), [#165](https://github.com/tibiadata/tibiadata-api-go/pull/165) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Bump golang from 1.19.0 to 1.19.2 ([#158](https://github.com/tibiadata/tibiadata-api-go/pull/158), [#159](https://github.com/tibiadata/tibiadata-api-go/pull/159) by [dependabot](https://github.com/dependabot))
- Bump golang.org/x/text from 0.3.7 to 0.4.0 ([#160](https://github.com/tibiadata/tibiadata-api-go/pull/160), [#161](https://github.com/tibiadata/tibiadata-api-go/pull/161) by [dependabot](https://github.com/dependabot))

## [3.4.1] - 2022-08-30

### Changed

- Bump golang from 1.18.4 to 1.19.0 ([#156](https://github.com/tibiadata/tibiadata-api-go/pull/156) by [dependabot](https://github.com/dependabot))

### Fixed

- Fix handling of special characters in boosted boss/creature name ([#157](https://github.com/tibiadata/tibiadata-api-go/pull/157) by [tiagomartines11](https://github.com/tiagomartines11))

## [3.4.0] - 2022-08-02

### Added

- Adding proxy protocol support to replace default https ([#154](https://github.com/tibiadata/tibiadata-api-go/pull/154) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- Parsing of spells list ([#155](https://github.com/tibiadata/tibiadata-api-go/pull/155) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.3.1] - 2022-07-29

### Changed

- Added bosspoints as additional HighscoreCategory ([#153](https://github.com/tibiadata/tibiadata-api-go/pull/153) by [Tholdrim](https://github.com/Tholdrim))

## [3.3.0] - 2022-07-27

### Added

- Add support to new boostable bosses page ([#152](https://github.com/tibiadata/tibiadata-api-go/pull/152) by [tiagomartines11](https://github.com/tiagomartines11))

### Changed

- Removing **v** in container image tag ([#144](https://github.com/tibiadata/tibiadata-api-go/pull/144) by [tobiasehlert](https://github.com/tobiasehlert))
- Updating `go build` step in dockerfile ([#148](https://github.com/tibiadata/tibiadata-api-go/pull/148) by [tobiasehlert](https://github.com/tobiasehlert))
- Some go mod and workflow build updates ([#137](https://github.com/tibiadata/tibiadata-api-go/pull/137) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump github.com/gin-gonic/gin from 1.7.7 to 1.8.1 ([#136](https://github.com/tibiadata/tibiadata-api-go/pull/136), [#140](https://github.com/tibiadata/tibiadata-api-go/pull/140) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.18.1 to 1.18.4 ([#135](https://github.com/tibiadata/tibiadata-api-go/pull/135), [#138](https://github.com/tibiadata/tibiadata-api-go/pull/138), [#149](https://github.com/tibiadata/tibiadata-api-go/pull/149) by [dependabot](https://github.com/dependabot))
- Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 ([#139](https://github.com/tibiadata/tibiadata-api-go/pull/139), [#141](https://github.com/tibiadata/tibiadata-api-go/pull/141), [#142](https://github.com/tibiadata/tibiadata-api-go/pull/142), [#143](https://github.com/tibiadata/tibiadata-api-go/pull/143) by [dependabot](https://github.com/dependabot))
- Bump github.com/gin-contrib/gzip from 0.0.5 to 0.0.6 ([#145](https://github.com/tibiadata/tibiadata-api-go/pull/145) by [dependabot](https://github.com/dependabot))

### Fixed

- Parsing issue when guild description contains founded string ([#151](https://github.com/tibiadata/tibiadata-api-go/pull/151) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.2.2] - 2022-05-12

### Changed

- Bump various workflow versions ([#126](https://github.com/tibiadata/tibiadata-api-go/pull/126), [#127](https://github.com/tibiadata/tibiadata-api-go/pull/127), [#128](https://github.com/tibiadata/tibiadata-api-go/pull/128), [#129](https://github.com/tibiadata/tibiadata-api-go/pull/129), [#130](https://github.com/tibiadata/tibiadata-api-go/pull/130), [#131](https://github.com/tibiadata/tibiadata-api-go/pull/131), [#132](https://github.com/tibiadata/tibiadata-api-go/pull/132) by [dependabot](https://github.com/dependabot))

### Fixed

- Fix traded-string appearing in account characters info ([#133](https://github.com/tibiadata/tibiadata-api-go/pull/133) by [tiagomartines11](https://github.com/tiagomartines11))

## [3.2.1] - 2022-04-26

### Changed

- Bump github/codeql-action from 1 to 2 ([#122](https://github.com/tibiadata/tibiadata-api-go/pull/122) by [dependabot](https://github.com/dependabot))

### Fixed

- Trimming suffix on guild rank ([#121](https://github.com/tibiadata/tibiadata-api-go/pull/121) by [tobiasehlert](https://github.com/tobiasehlert))
- Change go version in go.mod ([#124](https://github.com/tibiadata/tibiadata-api-go/pull/124) by [sergot](https://github.com/sergot))
- Adjusting parsing of fansite page ([#125](https://github.com/tibiadata/tibiadata-api-go/pull/125) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.2.0] - 2022-04-24

### Changed

- Adding support for finished auctions ([#115](https://github.com/tibiadata/tibiadata-api-go/pull/115), [#118](https://github.com/tibiadata/tibiadata-api-go/pull/118) by [tobiasehlert](https://github.com/tobiasehlert))
- Enhancing deaths parsing of characters ([#116](https://github.com/tibiadata/tibiadata-api-go/pull/116) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump golang from 1.17.8 to 1.18.1 ([#108](https://github.com/tibiadata/tibiadata-api-go/pull/108), [#113](https://github.com/tibiadata/tibiadata-api-go/pull/113) by [dependabot](https://github.com/dependabot))
- Bump various workflow versions ([#101](https://github.com/tibiadata/tibiadata-api-go/pull/101), [#102](https://github.com/tibiadata/tibiadata-api-go/pull/102), [#109](https://github.com/tibiadata/tibiadata-api-go/pull/109), [#110](https://github.com/tibiadata/tibiadata-api-go/pull/110), [#111](https://github.com/tibiadata/tibiadata-api-go/pull/111), [#112](https://github.com/tibiadata/tibiadata-api-go/pull/112) by [dependabot](https://github.com/dependabot))

### Fixed

- Fix regex to handle one bed in house ([#105](https://github.com/tibiadata/tibiadata-api-go/pull/105) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix encoding of apostrophes in multiple places ([#106](https://github.com/tibiadata/tibiadata-api-go/pull/106) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix guild description to contain guild details ([#107](https://github.com/tibiadata/tibiadata-api-go/pull/107) by [tobiasehlert](https://github.com/tobiasehlert))
- Fix characters missing marriage and some more tests ([#117](https://github.com/tibiadata/tibiadata-api-go/pull/117) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.1.1] - 2022-03-15

### Fixed

- Adding sanitize of nbsp string in death reason of players ([#100](https://github.com/tibiadata/tibiadata-api-go/pull/100) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.1.0] - 2022-03-10

### Changed

- Renaming Tibiadata to TibiaData ([#90](https://github.com/tibiadata/tibiadata-api-go/pull/90) by [tobiasehlert](https://github.com/tobiasehlert))
- Removing injection of houseType from assets ([#95](https://github.com/tibiadata/tibiadata-api-go/pull/95) by [tobiasehlert](https://github.com/tobiasehlert))
- Removing return of loyalty title on no title ([#98](https://github.com/tibiadata/tibiadata-api-go/pull/98) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump docker/login-action from 1.14.0 to 1.14.1 ([#87](https://github.com/tibiadata/tibiadata-api-go/pull/87) by [dependabot](https://github.com/dependabot))
- Bump actions/checkout from 2 to 3 ([#88](https://github.com/tibiadata/tibiadata-api-go/pull/88) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.17.7 to 1.17.8 ([#96](https://github.com/tibiadata/tibiadata-api-go/pull/96) by [dependabot](https://github.com/dependabot))

### Fixed

- Environment function logic fix ([#91](https://github.com/tibiadata/tibiadata-api-go/pull/91) by [tobiasehlert](https://github.com/tobiasehlert))
- Stop using ioutil as it is deprecated ([#92](https://github.com/tibiadata/tibiadata-api-go/pull/92) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Adding sanitize of nbsp string in death section of players ([#99](https://github.com/tibiadata/tibiadata-api-go/pull/99) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.0.0] - 2022-03-01

### The release of API v3

Bumping version to _v3.0.0_ so that GitHub releases match API version.

Head over to [tibiadata.com](https://tibiadata.com/2022/03/tibiadata-api-v3-released/) for more information.

## [0.6.2] - 2022-02-28

### Fixed

- Adding if to prevent panic due to regex ([#84](https://github.com/tibiadata/tibiadata-api-go/pull/84) by [tobiasehlert](https://github.com/tobiasehlert))

## [0.6.1] - 2022-02-27

### Changed

- Switching assets download URL to assets.tibiadata.com
- Bump docker/build-push-action from 2.8.0 to 2.9.0 ([#80](https://github.com/tibiadata/tibiadata-api-go/pull/80) by [dependabot](https://github.com/dependabot))
- Bump golang from 1.17.6 to 1.17.7 ([#81](https://github.com/tibiadata/tibiadata-api-go/pull/81) by [dependabot](https://github.com/dependabot))
- Bump docker/login-action from 1.12.0 to 1.13.0 ([#82](https://github.com/tibiadata/tibiadata-api-go/pull/82) by [dependabot](https://github.com/dependabot))

## [0.6.0] - 2022-01-31

### Added

- Addition of some more tests ([#76](https://github.com/tibiadata/tibiadata-api-go/pull/76) by [tobiasehlert](https://github.com/tobiasehlert))
- Implement graceful shutdown ([#78](https://github.com/tibiadata/tibiadata-api-go/pull/78) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Implement 404 page not found ([#77](https://github.com/tibiadata/tibiadata-api-go/pull/77) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

### Changed

- Switching to http status codes ([#79](https://github.com/tibiadata/tibiadata-api-go/pull/79) by [tobiasehlert](https://github.com/tobiasehlert))
- Cleanup of characters deathlist ([#75](https://github.com/tibiadata/tibiadata-api-go/pull/75) by [tobiasehlert](https://github.com/tobiasehlert))
- Adjustment of endpoint namings ([#74](https://github.com/tibiadata/tibiadata-api-go/pull/74) by [tobiasehlert](https://github.com/tobiasehlert))

## [0.5.1] - 2022-01-28

### Added

- Adding gzip compression middleware ([#73](https://github.com/tibiadata/tibiadata-api-go/pull/73) by [tobiasehlert](https://github.com/tobiasehlert))
- Swagger documentation annotations ([#67](https://github.com/tibiadata/tibiadata-api-go/pull/67) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Add some 'fake' unit tests to up coverage on webserver.go ([#71](https://github.com/tibiadata/tibiadata-api-go/pull/71) by [kamilon](https://github.com/kamilon))
- Add unit tests for House Overview and House APIs ([#72](https://github.com/tibiadata/tibiadata-api-go/pull/72) by [kamilon](https://github.com/kamilon))

## [0.5.0] - 2022-01-19

### Added

- Add unit tests for Spells Overview API, fix bugs in Spells Overview API ([#55](https://github.com/tibiadata/tibiadata-api-go/pull/55) by [kamilon](https://github.com/kamilon))
- Add unit tests for Spells API, fix various bugs in pulling of spell data ([#56](https://github.com/tibiadata/tibiadata-api-go/pull/56) by [kamilon](https://github.com/kamilon))
- Add unit tests for Kill Statistics API, move to goquery ([#57](https://github.com/tibiadata/tibiadata-api-go/pull/57) by [kamilon](https://github.com/kamilon))
- Add unit tests for Guild Overview and Guild API ([#58](https://github.com/tibiadata/tibiadata-api-go/pull/58) by [kamilon](https://github.com/kamilon))
- Add unit tests for Creatures Overview and Creature API, fix bug ([#59](https://github.com/tibiadata/tibiadata-api-go/pull/59) by [kamilon](https://github.com/kamilon))
- Add unit tests for Highscores API ([#61](https://github.com/tibiadata/tibiadata-api-go/pull/61) by [kamilon](https://github.com/kamilon))
- Add unit tests for Fansites API ([#62](https://github.com/tibiadata/tibiadata-api-go/pull/62) by [kamilon](https://github.com/kamilon))
- Add unit tests for News List API and News API ([#60](https://github.com/tibiadata/tibiadata-api-go/pull/60) by [kamilon](https://github.com/kamilon))
- Create and use HighscoreCategory enum ([#63](https://github.com/tibiadata/tibiadata-api-go/pull/63) by [kamilon](https://github.com/kamilon))
- Feature: Tibia Houses endpoints ([#26](https://github.com/tibiadata/tibiadata-api-go/pull/26) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Refactor gin server to reduce code duplication ([#64](https://github.com/tibiadata/tibiadata-api-go/pull/64) by [kamilon](https://github.com/kamilon))
- Cache new regex queries ([#66](https://github.com/tibiadata/tibiadata-api-go/pull/66) by [kamilon](https://github.com/kamilon))

### Fixed

- Fix race condition with TibiaDataRequestStruct ([#65](https://github.com/tibiadata/tibiadata-api-go/pull/65) by [kamilon](https://github.com/kamilon))
- Bump docker/build-push-action from 2.7.0 to 2.8.0 ([#68](https://github.com/tibiadata/tibiadata-api-go/pull/68) by [dependabot](https://github.com/dependabot))

## [0.4.0] - 2022-01-12

### Added

- Enhance data collection and debugging information ([#31](https://github.com/tibiadata/tibiadata-api-go/pull/31) by [tobiasehlert](https://github.com/tobiasehlert))
- Adding rate-limit detection and status code switching ([#36](https://github.com/tibiadata/tibiadata-api-go/pull/36) by [tobiasehlert](https://github.com/tobiasehlert))
- Adding TIBIADATA_HOST env ([#39](https://github.com/tibiadata/tibiadata-api-go/pull/39) by [tobiasehlert](https://github.com/tobiasehlert))
- Rewrite of characters and adding of unit testing ([#40](https://github.com/tibiadata/tibiadata-api-go/pull/40) by [kamilon](https://github.com/kamilon))
- Adding codecov.io ([#48](https://github.com/tibiadata/tibiadata-api-go/pull/48) by [tobiasehlert](https://github.com/tobiasehlert))
- Adding new Tibia News endpoints ([#32](https://github.com/tibiadata/tibiadata-api-go/pull/32) by [tobiasehlert](https://github.com/tobiasehlert))
- Worlds Overview unit tests ([#53](https://github.com/tibiadata/tibiadata-api-go/pull/53) by [kamilon](https://github.com/kamilon))
- World unit tests ([#54](https://github.com/tibiadata/tibiadata-api-go/pull/54) by [kamilon](https://github.com/kamilon))

### Changed

- Refactor to use pure goquery and no manual regex on guilds ([#29](https://github.com/tibiadata/tibiadata-api-go/pull/29) by [JorgeMag96](https://github.com/JorgeMag96))
- Performing go mod tidy ([#43](https://github.com/tibiadata/tibiadata-api-go/pull/43) by [tobiasehlert](https://github.com/tobiasehlert))
- Use new assert instead of passing t many times ([#45](https://github.com/tibiadata/tibiadata-api-go/pull/45) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Updating workflows ([#51](https://github.com/tibiadata/tibiadata-api-go/pull/51) by [tobiasehlert](https://github.com/tobiasehlert))
- Bump golang from 1.17.5 to 1.17.6 ([#50](https://github.com/tibiadata/tibiadata-api-go/pull/50) by [dependabot](https://github.com/dependabot))
- Reduce dependencies ([#44](https://github.com/tibiadata/tibiadata-api-go/pull/44) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Clean up webserver.go and a few utility funcs ([#37](https://github.com/tibiadata/tibiadata-api-go/pull/37), [#52](https://github.com/tibiadata/tibiadata-api-go/pull/52) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

### Fixed

- Maintenance mode detection with error message return ([#34](https://github.com/tibiadata/tibiadata-api-go/pull/34) by [tobiasehlert](https://github.com/tobiasehlert))
- Guilds name upper casing correction ([#35](https://github.com/tibiadata/tibiadata-api-go/pull/35) by [tobiasehlert](https://github.com/tobiasehlert))
- Handle timezone information during DateTime parsing ([#46](https://github.com/tibiadata/tibiadata-api-go/pull/46), [#49](https://github.com/tibiadata/tibiadata-api-go/pull/49) by [kamilon](https://github.com/kamilon), [tobiasehlert](https://github.com/tobiasehlert))

## [0.3.0] - 2022-01-05

### Added

- Addition of two guild endpoints ([#20](https://github.com/tibiadata/tibiadata-api-go/pull/20) by [tobiasehlert](https://github.com/tobiasehlert))
- Adding proxy support to replace default URL ([#25](https://github.com/tibiadata/tibiadata-api-go/pull/25) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed

- Clear lint errors ([#18](https://github.com/tibiadata/tibiadata-api-go/pull/18) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Idiomatic go ([#19](https://github.com/tibiadata/tibiadata-api-go/pull/19) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
- Removing duplicate function ([#27](https://github.com/tibiadata/tibiadata-api-go/pull/27) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- Highscores endpoints redirect ([#22](https://github.com/tibiadata/tibiadata-api-go/pull/22) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

## [0.2.0] - 2022-01-02

### Changed

- Sanitize certain data ([#14](https://github.com/tibiadata/tibiadata-api-go/pull/14) by [darrentaytay](https://github.com/darrentaytay))
- Removing some code smell from SonarCloud ([#10](https://github.com/tibiadata/tibiadata-api-go/pull/10) by [tobiasehlert](https://github.com/tobiasehlert))
- Moving vocation logic to separate function ([#17](https://github.com/tibiadata/tibiadata-api-go/pull/17) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed

- Fix decoding issues and save some allocations ([#15](https://github.com/tibiadata/tibiadata-api-go/pull/15), [#16](https://github.com/tibiadata/tibiadata-api-go/pull/16) by [kamilon](https://github.com/kamilon), [tobiasehlert](https://github.com/tobiasehlert))

## [0.1.1] - 2021-12-31

### Changed

- Bump actions/cache from 2.1.6 to 2.1.7 ([#2](https://github.com/tibiadata/tibiadata-api-go/pull/2) by [dependabot](https://github.com/dependabot))
- Bump docker/metadata-action from 3.6.0 to 3.6.2 ([#1](https://github.com/tibiadata/tibiadata-api-go/pull/1) by [dependabot](https://github.com/dependabot))
- Update README.md ([#3](https://github.com/tibiadata/tibiadata-api-go/pull/3), [#5](https://github.com/tibiadata/tibiadata-api-go/pull/5) by [tobiasehlert](https://github.com/tobiasehlert))
- Implementation of response handler ([#4](https://github.com/tibiadata/tibiadata-api-go/pull/4) by [tobiasehlert](https://github.com/tobiasehlert))
- Changing building and releasing ([#6](https://github.com/tibiadata/tibiadata-api-go/pull/6), [#7](https://github.com/tibiadata/tibiadata-api-go/pull/7) by [tobiasehlert](https://github.com/tobiasehlert))
- Shrink of TibiaFansitesV3 ([#8](https://github.com/tibiadata/tibiadata-api-go/pull/8) by [tobiasehlert](https://github.com/tobiasehlert))
- Updating build workflow with enhancements ([#11](https://github.com/tibiadata/tibiadata-api-go/pull/11) by [tobiasehlert](https://github.com/tobiasehlert))

## [0.1.0] - 2021-12-23

Initial commit

[4.7.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.6.0...v4.7.0
[4.6.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.5.0...v4.6.0
[4.5.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.4.0...v4.5.0
[4.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.6...v4.4.0
[4.3.6]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.5...v4.3.6
[4.3.5]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.4...v4.3.5
[4.3.4]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.3...v4.3.4
[4.3.3]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.2...v4.3.3
[4.3.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.1...v4.3.2
[4.3.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.3.0...v4.3.1
[4.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.2.3...v4.3.0
[4.2.3]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.2.2...v4.2.3
[4.2.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.2.1...v4.2.2
[4.2.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.2.0...v4.2.1
[4.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.1.2...v4.2.0
[4.1.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.1.1...v4.1.2
[4.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.1.0...v4.1.1
[4.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.0.2...v4.1.0
[4.0.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.0.1...v4.0.2
[4.0.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v4.0.0...v4.0.1
[4.0.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.5...v4.0.0
[3.7.5]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.4...v3.7.5
[3.7.4]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.3...v3.7.4
[3.7.3]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.2...v3.7.3
[3.7.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.1...v3.7.2
[3.7.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.7.0...v3.7.1
[3.7.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.6.0...v3.7.0
[3.6.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.5.1...v3.6.0
[3.5.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.5.0...v3.5.1
[3.5.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.4.1...v3.5.0
[3.4.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.3.1...v3.4.0
[3.3.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.3.0...v3.3.1
[3.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.2...v3.3.0
[3.2.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.1.1...v3.2.0
[3.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.0.0...v3.1.0
[3.0.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.2...v3.0.0
[0.6.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/30f328f...v0.1.0
