# TerrAPI

[![Current Tag](https://img.shields.io/github/v/tag/thomas-illiet/terrapi?sort=semver)](https://github.com/thomas-illiet/terrapi) [![Build Status](https://github.com/thomas-illiet/terrapi/actions/workflows/general.yml/badge.svg)](https://github.com/thomas-illiet/terrapi/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/thomas-illiet/terrapi.svg)](https://pkg.go.dev/github.com/thomas-illiet/terrapi) [![Go Report Card](https://goreportcard.com/badge/github.com/thomas-illiet/terrapi)](https://goreportcard.com/report/github.com/thomas-illiet/terrapi) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/d2bc4877341f4c7fbf9b4fa62b8d0484)](https://www.codacy.com/gh/thomas-illiet/terrapi/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=thomas-illiet/terrapi&amp;utm_campaign=Badge_Grade)

TerrAPI is a RESTful API service designed to simplify and automate your Terraform deployments. With a playful nod to "therapy," TerrAPI helps ease the stress of managing complex infrastructure by providing an easy-to-use API interface to trigger, manage, and monitor Terraform deployments.

This project allows developers and DevOps teams to integrate Terraform workflows into their existing applications and services effortlessly. TerrAPI handles the heavy lifting, so you can focus on what matters — building scalable, resilient infrastructure. Whether you're deploying infrastructure on AWS, GCP, or other cloud providers, TerrAPI offers a seamless experience for infrastructure-as-code automation.

Let TerrAPI be the "therapy" your infrastructure needs!

## Development

Make sure you have a working Go environment, for further reference or a guide
take a look at the [install instructions][golang]. This project requires
Go >= v1.17, at least that's the version we are using.

```console
git clone https://github.com/thomas-illiet/terrapi.git
cd terrapi

make generate build

./bin/terrapi -h
```

## Security

If you find a security issue please contact
[contact@thomas-illiet.fr](mailto:contact@thomas-illiet.fr) first.

## Contributing

Fork -> Patch -> Push -> Pull Request

## Authors

- [Thomas ILLIET](https://github.com/thomas-illiet)

## License

Apache-2.0

## Copyright

```console
Copyright (c) 2024 Thomas ILLIET <contact@thomas-illiet.fr>
```
