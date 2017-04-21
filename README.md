# turbo-goprobe-example
Follow the steps below to try out this example GO probe.  This assumes you already have your Turbonomic OpsMgr installed.
1. Customize `turbo-server-conf.json` according to your Turbo OpsMgr.
2. Customize `target-conf.json` to your preference.
3. Run `go install ./...` to build and install.
4. Start the example probe: `./turbo-goprobe-example`
5. Confirm in your OpsMgr that a target has been created as specified in `target-conf.json`.
6. Inspect inventory in your OpsMgr that some example entities have been created as well.

More to come as how to write your own GO probe ...
