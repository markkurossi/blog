# Being Social Media Friendly

After the first published post last Friday, I started checking how
article links would appear on social media platforms. And they did
appear ugly! The templates and the publishing system were missing
various HTML
[metadata](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta)
elements so that the posts would be missing all preview images and
information. Another problem was the site URL structure. The first
version produced all outputs under one output directory. When the
number of posts increases and they have additional assets like images,
it will likely create naming conflicts in the output directory. This
article describes how I fixed these limitations, i.e., how the system
became more social media friendly.

## Article Metadata

The article `settings.toml` file has a new `[meta]` category for
specifying the metadata attributes:

    [meta]
    Title = "Article metadata title."
    Description = "Article metadata description."

The following settings values are currently supported:

`Title`
: Defines an optional metadata title value. The document title
defaults to the article title if you omit this value.

`Description`
: Defines an optional article description. If you do not set this
value, the system will not generate the description meta-attributes.

With these settings, all template files have now been updated to
contain the [Twitter
Cards](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/markup)
and [Open Graph
Sharing](https://developers.facebook.com/docs/sharing/webmasters/)
markup tags:

    <meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@markkurossi">
    <meta name="twitter:title" content="{{.MetaTitle}}">

    <meta name="og:type" content="blog">
    <meta name="og:title" content="{{.MetaTitle}}">

The Facebook [Sharing
Debugger](https://developers.facebook.com/tools/debug/) proposes that
the `og:image` property should be explicitly defined, even if the
value could be inferred from other tags. The image URL must be
absolute, so the default template now contains the hard-coded URL
values pointing to my site. These could be moved to the site-level
configuration file when that is implemented. With this change, there
could be different images for different categories, and the
configuration file mapping could do the mapping from tags to
URLs. That said, the `mtr` template files now define the
`twitter:image` and `og:image` URLs (thanks to
[Iconmonstr](https://iconmonstr.com/) for the icons):

    <meta name="twitter:image" content="https://.../iconmonstr-file-22-240.png">
    <meta name="og:image" content="https://.../iconmonstr-file-22-240.png">

And as mentioned earlier, the metadata description tags are
optional. They are omitted with the [Go
Template](https://pkg.go.dev/text/template) if-actions when not
provided:

    {{if .MetaDescription}}
    <meta name="description" content="{{.MetaDescription}}">
    <meta name="twitter:description" content="{{.MetaDescription}}">
    <meta name="og:description" content="{{.MetaDescription}}">
    {{end}}

## URL Structure

The initial version of the blog system put all output files in the one
output directory. This model doesn't work well when there are more
articles with various additional assets; there is a very likely risk
for name conflicts between different assets and files. So how to fix
this? The first idea was to put all articles in a separate directory,
named by the article title. This is a safe assumption since the
article titles should be unique. However, it means that the default
document would be `index.html` for each article. This would be ok, and
many static site generators use this approach (or they can be
configured this way). However, I wanted to keep the "index" file to be
named after the document title, i.e., it would have the Website
[Slug](https://developer.mozilla.org/en-US/docs/Glossary/Slug) encoded
in the final path element. So I ended up using the following naming
algorithm:

 - Article assets are stored under a directory named by the article
   publishing date: `2021-12-24`.
 - If there are multiple articles published on the same day, they are
   sorted by their publishing time, and the article asset directory is
   appended with the order prefix: "-_nth_." The articles have a
   canonical order (publishing date must be unique, the file
   modification dates order drafts), so there won't be conflicts.
 - The index file inside the directory is named by the article title
   (being-social-media-friendly.html for this article).

This algorithm defines unique URLs for all articles. The article
assets are stored in an article-specific directory, so there are no
naming conflicts. This model makes it also easy to reference article
assets as the references are the same relative ones in the source and
output directories.
