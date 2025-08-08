# Changelog

## [1.0.0](https://github.com/marouanesouiri/yada/compare/v0.2.1...v1.0.0) (2025-08-08)


### ⚠ BREAKING CHANGES

* Method signatures changed from returning *Call[T] to returning (T, error) directly.
* existing enum constant names have changed; update references accordingly.

### Features

* add all Discord channel types with clean, consistent documentation ([436584b](https://github.com/marouanesouiri/yada/commit/436584b359e66ff893f1c1b8397e92e7d2de50a1))
* add attachment.go with Attachment struct and methods ([48ddc06](https://github.com/marouanesouiri/yada/commit/48ddc064da1d0c99a1ddc70835c803db2a6c6363))
* add Builder method and field management helpers ([a8d1b93](https://github.com/marouanesouiri/yada/commit/a8d1b9343a38744dc1db60ca840cc6d55eb66ec1))
* add ChannelType.Is helper method ([aeab0bd](https://github.com/marouanesouiri/yada/commit/aeab0bde0494ef7bc697de898798e93e2bec3af0))
* add Cluster type for high-level Discord bot management ([83f5fef](https://github.com/marouanesouiri/yada/commit/83f5fef9b4b995c7b1c5fa553e2a0433a527b41d))
* add comprehensive Discord User types with full documentation ([3eca19e](https://github.com/marouanesouiri/yada/commit/3eca19ed037d042dc7bc5e163baab4735b0ce7f6))
* add DisplayAvatarURL, DisplayAvatarURLWith, DisplayName, and Mention helpers ([c84732b](https://github.com/marouanesouiri/yada/commit/c84732b65607f287b98135a11c4e6b8425b03a94))
* add Embed struct and builder for Discord embeds ([0e9c150](https://github.com/marouanesouiri/yada/commit/0e9c150c5f515b6901bcc289694694424b2cb3eb))
* add Emoji struct and Mention() helper for custom emojis ([1f58003](https://github.com/marouanesouiri/yada/commit/1f5800308e5359b7638e101f9764d048979d9c9a))
* add GatewayIntent constants, opcode definitions, and close event codes ([4535d54](https://github.com/marouanesouiri/yada/commit/4535d5468d7be16de9c293bed93af9126d1a5462))
* add JSON marshal/unmarshal for Snowflake type ([26cef92](https://github.com/marouanesouiri/yada/commit/26cef928e2f7a202738a5c0f1b292e0f436b6c75))
* add JumpURL() method; reorder channel.go symbols for clarity ([87231cd](https://github.com/marouanesouiri/yada/commit/87231cd3fc53d75ceb01d4f98fbbc62a02debb13))
* add MessageCreateEvent and MessageDeleteEvent structs with JSON unmarshal logic ([f27b5ac](https://github.com/marouanesouiri/yada/commit/f27b5acb54f38686a95cf8f0246118c4a974e397))
* add methods to get emoji image URLs ([e3a637d](https://github.com/marouanesouiri/yada/commit/e3a637d6d72429bdf6ed9f0c8e509a1949587b68))
* add NewImageFile helper to load image files as base64 data URIs ([5465dc1](https://github.com/marouanesouiri/yada/commit/5465dc1956c922fb360758978da3174d6c75dfa3))
* add Permissions type with constants and bitmask operations ([7a587e9](https://github.com/marouanesouiri/yada/commit/7a587e90339dea5d3561d19343de6caff2143e1d))
* add ReadyEvent struct with JSON unmarshal method ([caf4cb9](https://github.com/marouanesouiri/yada/commit/caf4cb976434302e7e868b24fd0410bfd976f8c2))
* add readyHandlers for READY event handling ([b3c5b9b](https://github.com/marouanesouiri/yada/commit/b3c5b9b7a9eb6f819982bc5fb1efaf8ab69694ac))
* add Role, RoleColors, RoleTags, RoleFlags structs with helpers ([2ca4bf9](https://github.com/marouanesouiri/yada/commit/2ca4bf9dc70bf080659c0cb64162ee79c1a8ef6e))
* add shard management with websocket connection, heartbeat, and identify rate limiting ([2485ea0](https://github.com/marouanesouiri/yada/commit/2485ea0f1dd6d1d37fc166c71104ecfad09f65da))
* add simple dynamic worker pool ([7145d16](https://github.com/marouanesouiri/yada/commit/7145d1673ce0444db1af4887c79aed472a908c75))
* add Sticker and StickerPack ([02e537e](https://github.com/marouanesouiri/yada/commit/02e537ecfa5d5c9802ff28ea5fa98e481e35d5ab))
* add StickerPackBannerURL helper with format and size options ([c4b2eab](https://github.com/marouanesouiri/yada/commit/c4b2eabb61f96ab9d56a5cd47d5144aedd3eab73))
* add support for X-Audit-Log-Reason header ([2be50d2](https://github.com/marouanesouiri/yada/commit/2be50d23a1080da6a8dcd279c4d8cc1d6e7a746e))
* add unified image URL generation helpers for all Discord endpoints ([7a26706](https://github.com/marouanesouiri/yada/commit/7a2670661cb801f2b9e9c03e10399c4d9b927cab))
* add UnSet method to Snowflake for zero-value checks ([c0ed56e](https://github.com/marouanesouiri/yada/commit/c0ed56e8fc3aa7f859b0c7f7a80c08b01f9b941d))
* add user endpoints with clear docs and usage examples ([29efef6](https://github.com/marouanesouiri/yada/commit/29efef6d528cbc35452347729cdf835a4c2e7fb7))
* implement event dispatcher with handler registration ([ad4a19d](https://github.com/marouanesouiri/yada/commit/ad4a19d2ab0cdef0440de75ce5dda8a106176862))
* implement MESSAGE_CREATE and MESSAGE_DELETE handlers managers ([e650cac](https://github.com/marouanesouiri/yada/commit/e650cace61a85060f7480961fbd04dbd713ebfb5))
* implement Shutdown for requester, restApi, and Shard ([d6c3b76](https://github.com/marouanesouiri/yada/commit/d6c3b76e78a852176caaae80d854a85a8ab56946))
* improve restapi.go docs, fix pointer usage in unmarshal, add ModifySelfUser with params and JSON marshalling ([a7f9eef](https://github.com/marouanesouiri/yada/commit/a7f9eef3aba1477b7c0adca8375598f60d46d738))
* use worker pool for async event handling ([facdbfd](https://github.com/marouanesouiri/yada/commit/facdbfde124458bd2991cb26363b3ba3ea59a551))


### Bug Fixes

* **ci:** correct branch and token in release workflow ([994b9e6](https://github.com/marouanesouiri/yada/commit/994b9e69a251a0c7c291c821f578cdd627360759))
* correctly pass event structs by pointer to sonic.Unmarshal and fix logs messages ([055e06a](https://github.com/marouanesouiri/yada/commit/055e06a2404da583af34ce7c5ac28e5002433069))
* handle 401 Unauthorized responses with error and log ([68048c3](https://github.com/marouanesouiri/yada/commit/68048c3fca8ec5ddd0adfef9a46a6d6f49a6bce2))
* make bucket keys include major param for per-resource rate limits ([44e84fb](https://github.com/marouanesouiri/yada/commit/44e84fb82056754962447d08113ab5d277b86582))
* pass correct payload data to dispatcher in readLoop ([083e309](https://github.com/marouanesouiri/yada/commit/083e309dd20f3176f4c2d2a4dcf8d36a3fcf28e0))
* prevent recursive String() call by casting Color to int64 in Sprintf ([3a911a3](https://github.com/marouanesouiri/yada/commit/3a911a3c27338b6ba20a7c4cb7c25f6c3915ddfc))
* replace authNotRequired with authWithToken for correct requester calls ([499f33a](https://github.com/marouanesouiri/yada/commit/499f33aba0c08124d40d5f92511507703a3af222))
* return unicode emoji name in Mention if ID is zero ([86bdc73](https://github.com/marouanesouiri/yada/commit/86bdc7384ea7a50b9f06eead59e2278ea875dcc8))


### Code Refactoring

* rename all enum constants to camel-case style ([df1aede](https://github.com/marouanesouiri/yada/commit/df1aede922ad5c85a269f1ea28596cd8f50510d7))
* rename Get* methods to Fetch* and simplify REST calls ([165b63d](https://github.com/marouanesouiri/yada/commit/165b63dc7283a4913305a0df47f18317b781d89b))
