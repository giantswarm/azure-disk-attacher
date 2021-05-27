[![CircleCI](https://circleci.com/gh/giantswarm/azure-disk-attacher.svg?style=shield)](https://circleci.com/gh/giantswarm/azure-disk-attacher)

# Azure disk attacher

This is a tool that allows to attach managed disks to Azure VMSS instances without invalidating the instance's `LatestModelApplied` field.

This is needed to circumvent a bug in `az` cmd line tool.

See https://github.com/Azure/azure-cli/issues/18259
