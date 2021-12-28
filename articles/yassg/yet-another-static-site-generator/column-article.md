[//]: # (-*- markdown -*-)

# Yet Another Static Site Generator

## Abstract

The world needs Yet Another Static Site
Generator--[YASSG](https://github.com/markkurossi/blog)! Or most
probably not, but I wanted to write some blog-like texts and didn't
(quickly) find any tools that would support my all requirements:

 - Zero dependencies: I don't want to install tons of tools and
   dependencies. In the spirit of [Go](https://go.dev/), I want a
   single binary that will generate all outputs.
 - Source code must be ASCII text to be easily handled with version
   control tools like Git.
 - Support for several data visualizations:
   - 2D ASCII graphics converted to 2D line drawings.
   - Histograms and bar charts.
   - [Tabulating](https://github.com/markkurossi/tabulate) text.
 - Simple templates: document fragments are generated from Markdown
   files and placed in templates. The templates define all font and
   CSS definitions.
 - Automatic word index generation to create cross-references between
   posts.

The framework will be opinionated, supporting convention over
configuration, use external social media platforms for comments
(without any internal commenting system), and tuned for my personal
preferences instead of targeting general-purpose usefulness.

## High-Level Design

The system assumes that all blog article input materials are stored in
the following directory hierarchy:

    articles/
             index/
                   settings.toml
                   column-article.md
             yassg/yet-another-static-site-generator/
                                                     settings.toml
                                                     column-article.md
             mpc/
                 settings.toml
                 column-article.md
             ...

The `index` is the only pre-defined article name. It will create the
blog site `index.html` page. All other article sub-directories
(`yassg`, `mpc` above) can be freely named, and they will produce a
separate blog post. All sub-directories between the article root
(`articles` above) and the article directory (which contains the
Markdown files) will be tags for the article. For example, all
articles in this series will have the `yassg` tag.

The markdown files under the article sub-directories are named
matching template variables. The default template `mtr` defines one
expansion variable `{{.ColumnArticle}}`. Its content will be generated
from the `column-article.md` input file. The blog system automatically
converts from the dashed markdown file name to the camel case.

The `settings.toml` file specifies the article settings:

    [article]
    Title = "YASSG - Yet Another Static Site Generator"
    Tags = ["yassg", "go"]

The following settings values are currently supported:

 - `Title`: specifies the HTML document title
 - `Tags`: select blog post tags. They are used to construct a
   category tag cloud for navigation

## Roadmap

This blog category (yassg) follows the program development. The
current version is very primitive. It can create simple static pages
without any advanced features, but it can produce the correct HTML
output you can see on this page. The following sections describe the
high-level development items in expected implementation order.

### 2D Graphics and Inline Data

The blog platform will use Markdown code blocks for inline data and
custom data processing. The main idea is to keep raw data, formatted
output, and the article text in the same document. The initial idea is
to mark data blocks with `data` info:

    ```data:csv
    Year,Base,Delta
    2019,100,7
    2020,120,12
    2021,140,15
    ```

This data block can be referenced with `data:`_nth_ URL where _nth_
specifies the 0-based data block index in the document:

    ```iql -t uclight
    select * from data:0
    ```

This block would be replaced with the output from the
[iql](https://github.com/markkurossi/iql) command:

```
┌──────┬──────┬───────┐
│ Year │ Base │ Delta │
├──────┼──────┼───────┤
│ 2019 │  100 │     7 │
│ 2020 │  120 │    12 │
│ 2021 │  140 │    15 │
└──────┴──────┴───────┘
```

### Word Index

Automated index generation will extract words from the blog texts and
create an inverted index from words to articles. The indexing system
could use part-of-speech (POS) tagging to index only nouns and foreign
words or make a full-text index (excluding stopwords).

The word index would also allow full-text searching of blog
articles. The search logic could be implemented in JavaScript,
maintaining the idea of a CDN-friendly static site.

### ChangeLog and Publishing Automation

If Git manages the blog article version control, the blog command can
create automated released and updated tags and changelog of article
content changes. Note that since the article content (Markdown files)
is separate from the publishing toolchain and layout, the actual
output can and will change even if the article data does not change.

The blog commenting requires that the article link to the announcement
Tweet (or other social media platform posts). The publish automation
might automate this by creating a Twitter post about the article,
getting its URL, updating the comment link in the article, and finally
committing to the new updated version.
