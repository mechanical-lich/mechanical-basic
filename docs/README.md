# Mechanical Basic Documentation

This folder contains the documentation site for Mechanical Basic, built with Jekyll and the Just the Docs theme.

## GitHub Pages

The site is automatically deployed to GitHub Pages from this folder.

## Local Development

To test the documentation locally:

1. Install Ruby and Bundler (if not already installed)
2. Install dependencies:
   ```bash
   cd docs
   bundle install
   ```
3. Serve the site locally:
   ```bash
   bundle exec jekyll serve
   ```
4. Open http://localhost:4000/mechanical-basic/ in your browser

## Adding Pages

1. Create a new `.md` file in this folder
2. Add front matter with navigation order:
   ```yaml
   ---
   layout: default
   title: Your Page Title
   nav_order: 6
   ---
   ```
3. Write your content in Markdown

Pages will automatically appear in the sidebar navigation based on `nav_order`.

## Theme Documentation

This site uses the Just the Docs theme. For more configuration options, see:
https://just-the-docs.github.io/just-the-docs/
