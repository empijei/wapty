# Initial stage
This stage will be the first stage for wapty, before this is finished wapty will likely have unstable APIs and won't be really usable.

1. finish Repeat tests
1. use templates for UI
1. Add UI to repeater
1. Add saving functionality
1. Add scoping
1. Add history filtering/sorting
1. releasing the intercept should forward all pending requests
1. Add internal router
1. ignore recursive connect
1. Add intercept checker in the right spots
1. Allow the user to change the destination endpoint
1. Serve the certificates on a specific fake host/path
1. Add req ID to editor and reject unexpected requests
1. Add Intruder
1. [UI] Send the whole status on ui connect
1. [UI] Sanitize metadata
1. [UI] show already pending request/response upon connection
1. [UI] error log
1. [UI] auto-open ui in browser on launch
1. [UI] monospace textareas
1. Look for fixmes and todos in the code
1. Provide a ui to the decoder

The following is just some general polishing before calling this a proper project
1. Improve README
1. Handle panics within the package
1. Move all constant strings to actual constants
1. All the deferred closes if err!= nil send that, otherwise propagate the new one
1. Doc comment should be a complete sentence that starts with the name being declared.
1. general code polish, doc and and testing

# Moving to Release
This is meant to be mostly an improvement, adding features that are less used in burpsuite but are still there and should end up in wapty before it is called a proper replacement for burp

1. Add AutoEdit
1. Add cURL converter
1. Default to bare sockets on error
1. profile the code, try to find limit-cases
1. Add Spider (remember to add timeouts §8.10)
1. Add Scanner
1. Add recursive intruder with flows

# Release
1. Have penetration testers use wapty for a while, collect feedback
1. Implement fixes, add suggestions to a feature list
1. Advertise and publish the project on a broader scale

# Improvements
This section contains the features that burpsuite lacks but that will make this project different :)

These features will probably be implemented along with the ones in the other stages.

1. Add Mocksy
1. Add pre-engagement analysis/recon
1. Add a Pathfinder feature to spider that allows to backtrace how a certain URL was discovered
1. Add a Plugin manager / Make plugin behave as package testing, just plug the stuff
1. Add a SQLmap invoker
1. Add SAML, JWT decoder/editor
1. Add fuzzing payloads generator

# Misc:
These are the feature I still don't know if are worth adding

(PRs are welcome :D )

1. Add Content-Length override
1. Add Beautifier
1. Decompress HTTP2 instead of disabling it
1. [UI] Make operations unblocking and detect ui freezes/deaths. If channel is full and not being read, kill the client.
