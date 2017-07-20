package main

import "testing"
import "github.com/stretchr/testify/assert"

const markdownWithMetadata = `---
title: My article
tags: [graphql, ruby]
---

# My article

Lorem ipsum dolor sit amet, consectetur adipiscing elit.`

const markdownWithoutMetadata = `
# My article

Lorem ipsum dolor sit amet, consectetur adipiscing elit.`

const markdownWithInvalidMetadata = `---
invalid yaml
---

# My article`

func TestExtractMarkdownMetadata(t *testing.T) {
	metadata, strippedMarkdown, err := ExtractMarkdownMetadata([]byte(markdownWithMetadata))

	assert.Nil(t, err)

	assert.Equal(t, []string{"graphql", "ruby"}, metadata.Tags)
	assert.Equal(t, markdownWithoutMetadata, string(strippedMarkdown))
}

func TestExtractMarkdownMetadataWithoutMetadata(t *testing.T) {
	metadata, strippedMarkdown, err := ExtractMarkdownMetadata([]byte(markdownWithoutMetadata))

	assert.Nil(t, err)
	assert.Nil(t, metadata.Tags)
	assert.Equal(t, markdownWithoutMetadata, string(strippedMarkdown))
}

func TestExtractMarkdownMetadataWithInvalidYAML(t *testing.T) {
	metadata, strippedMarkdown, err := ExtractMarkdownMetadata([]byte(markdownWithInvalidMetadata))

	assert.Nil(t, metadata)
	assert.Nil(t, strippedMarkdown)
	assert.NotNil(t, err)
}
