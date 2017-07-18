# md2medium

Write your blog posts in Markdown and publish them on Medium using `md2medium`.

![](screenshot3.png)

## Install

`md2medium` requires Go 1.8 and can be installed using:

```
go get -u github.com/cjoudrey/md2medium
```

## Features

**Syntax highlighting**

Medium does not support syntax highlighting of code blocks. Instead, you must
create gists and embed the gists in your article.

`md2medium` makes this chore much easier by automatically creating gists
for any code blocks in your markdown file:

![](screenshot1.png)

**Image caption**

When specifying an `alt` for an image in Markdown, `md2medium` will include it
under the image:

![](screenshot2.png)


## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/cjoudrey/md2markdown. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The package is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
