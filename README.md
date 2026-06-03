# 💬 Prompts

[![Release](https://github.com/zeiss/prompts/actions/workflows/main.yml/badge.svg)](https://github.com/zeiss/prompts/actions/workflows/main.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeiss/prompts.svg)](https://pkg.go.dev/github.com/zeiss/prompts)
[![Go Report Card](https://goreportcard.com/badge/github.com/zeiss/prompts)](https://goreportcard.com/report/github.com/zeiss/prompts)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

A teeny-tiny experimental package to prompt for answers in [Ollama](https://ollama.com/) and other OpenAI-compatible APIs.

## Use

The idea is to have a simple, minimalistic package that enables you to prompt for answers without the overhead of a full-fledged client library. The focus is on composable packages and a consistent API surface that can be used across different providers.

## Supported Schemas

⚠️ The focus is on the modern Response API pioneered by [OpenAI](https://openai.com).

| Provider | Response API (compact) | Chat Completion API | Streams
|---|---|---|---|
| [Ollama](https://ollama.com/) | ✅ | 🛑 | ✅ |

## Docs

You can find the documentation hosted on [godoc.org](https://godoc.org/github.com/zeiss/prompts).

## Examples

The examples are located in the [examples](/examples) directory.

## License

[Apache 2.0](/LICENSE)
