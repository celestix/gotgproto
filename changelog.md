# GoTGProto Changelog

- Updated GoTD to v0.81.0
- Deprecated `ctx.ForwardMessage` (Use `ctx.ForwardMessages` now)
- Replaced BigCache with AnimeKaizoku/Cacher for caching peers
- Fixed high memory usage by gotgproto (uses around 5MBs now, earlier it was 100+ MBs)
- Rewrote GoTGProto client, it should be more handy to create a new client now
- Added a new `dispatcher.Dispatcher` inteface
- Renamed `dispatcher.CustomDispatcher` to `dispatcher.NativeDispatcher`
- Optimised command and message handlers
- Added new `types.Message`, which is a union of `tg.Message`, `tg.MessageService`, `tg.MessageEmpty`
- `ext.Update.EffectiveMessage` is of type `*types.Message`
- Added a new optional field in ClientOpts, named `AutoFetchReply` (setting this field to true will automatically cast ReplyToMessage field)