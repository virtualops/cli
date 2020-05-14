# Breeze CLI

<div style="display: flex; align-items: center; justify-content: center">
    <img src="https://duez0tpxkp9od.cloudfront.net/ec8f2fa2-f1d4-41a6-8873-98aa863bb875/svg/logo-new.svg" alt="Breeze Logo" />
</div>

## Introduction

The Breeze CLI lets you run any application in a local Kubernetes
cluster, without having to configure any Kubernetes manifests.

Currently, only a Laravel preset exists, which will automatically
handle building a Docker image for your Laravel app and deploying it
to a local cluster using [Tilt](https://tilt.dev).

## Installation

Until Breeze is released as a stable version, we recommend always
downloading the latest binary from our [`latest`](https://github.com/virtualops/cli/releases/tag/latest) release,
which is automatically published on every new commit to the master branch.

Download the binary for your platform, and move it to `/usr/local/bin/breeze`.

## Usage

If you've never used Kubernetes before, make sure you have Kubernetes locally.
You can get Kubernetes by [downloading Docker Desktop](https://www.docker.com/products/docker-desktop),
and enabling the Kubernetes cluster from the settings.

Once you have Kubernetes active, run

```bash
breeze cluster setup
```

This should download required dependencies and setup the cluster ready for use.
 
### Laravel

To use Breeze with Laravel, run

```bash
breeze init
breeze dev
```

`breeze init` will create a `breeze.yaml` file in your project, which contains
information about your build.

`breeze dev` will start Tilt with a generated Docker image, optimised for Laravel.
Note that `breeze dev` also creates a `.breezedev` directory. This directory is
only used temporarily while `breeze dev` is running, and should not be deleted.

## Where's the Windows release?

What's a windows?
