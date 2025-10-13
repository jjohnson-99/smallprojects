# static_site

This is a rudimentary script which crawls a source directory for markdown files, converts
the text into HTML, places the HTML text into a specified HTML template, then copies the resulting
file into a specified destination directory with the same file structure as the source. 

## Limitations
* Nested inline markdown is not supported. For example, in markdown 
**we can have *italicised text* inside of bold text.**
