import unittest
from block_markdown import (
        markdown_to_blocks,
        block_to_block_type,
        block_tag,
        trim_markdown_type_indicator,
        separate_list_items,
        block_to_parentnode,
        helper_list_items,
        markdown_to_html_node
)
from parentnode import ParentNode
from textnode import TextNode


class TestBlockTag(unittest.TestCase):
    def test_tag_h1(self):
        text = "# This is a heading"
        tag = block_tag(text)
        self.assertEqual("h1", tag)
                
    def test_tag_h6(self):
        text = "###### This is a heading"
        tag = block_tag(text)
        self.assertEqual("h6", tag)
               
    def test_tag_code(self):
        text = "``````"
        tag = block_tag(text)
        self.assertEqual("pre", tag)

    def test_tag_paragraph(self):
        text = "This is text"
        tag = block_tag(text)
        self.assertEqual("p", tag)
                
class TestTrimMarkdown(unittest.TestCase):
    def test_trim_heading(self):
        text = "#### This is a heading"
        trimmed_text = trim_markdown_type_indicator(text)
        self.assertEqual("This is a heading", trimmed_text)
 
    def test_trim_empty_code(self):
        text = "``````"
        trimmed_text = trim_markdown_type_indicator(text)
        self.assertEqual("``", trimmed_text)
     
    def test_trim_quote(self):
        text = ">This is a quote\n>with multiple\n> lines \n "
        trimmed_text = trim_markdown_type_indicator(text)
        self.assertEqual("This is a quote\nwith multiple\n lines", trimmed_text)

    def test_trim_other(self):
        text = "* This is a list item"
        trimmed_text = trim_markdown_type_indicator(text)
        self.assertEqual("* This is a list item", trimmed_text)
        
class TestSeparateList(unittest.TestCase):
    def test_separate_list(self):
        text = "* list item 1\n- item 2\n3. item3"
        separate_list = separate_list_items(text)
        self.assertEqual(["list item 1", "item 2", "item3"], separate_list)

class TestMarkdownToBlock(unittest.TestCase):
    def test_all(self):
        markdown = "# This is a heading\n\nThis is a paragraph of text. It has some **bold** and *italic* words inside of it.\n\n\n\n\n* This is the first list item in a list block\n* This is a list item\n* This is another list item\n"
        blocks = markdown_to_blocks(markdown)
        self.assertListEqual(
            [
                "# This is a heading",
                "This is a paragraph of text. It has some **bold** and *italic* words inside of it.",
                "* This is the first list item in a list block\n* This is a list item\n* This is another list item"
            ],
            blocks
        )

class TestBlockToParentNode(unittest.TestCase):
    def test_heading_to_parent(self):
        text = "### Heading with **bold** and *italic*"
        parentnode = block_to_parentnode(text)
        expected_output = "<h3>Heading with <b>bold</b> and <i>italic</i></h3>"
        self.assertEqual(expected_output, parentnode.to_html())

    def test_paragraph_to_parent(self):
        text = "Paragraph with **bold** and *italic*"
        parentnode = block_to_parentnode(text)
        expected_output = "<p>Paragraph with <b>bold</b> and <i>italic</i></p>"
        self.assertEqual(expected_output, parentnode.to_html())

    def test_code_block_to_parent(self):
        text = "```print('hello world')```"
        parentnode = block_to_parentnode(text)
        expected_output = "<pre><code>print('hello world')</code></pre>"
        self.assertEqual(expected_output, parentnode.to_html())

    def test_quote_to_parent(self):
        text = "> this is a quote\n>words\n> more words"
        parentnode = block_to_parentnode(text)
        expected_output = "<blockquote> this is a quote\nwords\n more words</blockquote>"
        self.assertEqual(expected_output, parentnode.to_html())

    def test_unordered_list_to_parent(self):
        text = "* list item **bold**\n"
        tag = block_tag(text)
        list_item_nodes = helper_list_items(text)
        parentnode = ParentNode(tag, list_item_nodes)
        expected_output = "<ul><li>list item <b>bold</b></li></ul>"
        self.assertEqual(expected_output, parentnode.to_html())

    def test_unordered_list_to_parent(self):
        text = "* list item **bold**\n* list item *italic*\n- item3"
        tag = block_tag(text)
        list_item_nodes = helper_list_items(text)
        parentnode = ParentNode(tag, list_item_nodes)
        expected_output = ("<ul><li>list item <b>bold</b></li>"
                           "<li>list item <i>italic</i></li>"
                           "<li>item3</li></ul>"
                           )
        self.assertEqual(expected_output, parentnode.to_html())

    def test_ordered_list_to_parent(self):
        text = "1. list item **bold**\n2. list item *italic*\n3. item3"
        tag = block_tag(text)
        list_item_nodes = helper_list_items(text)
        parentnode = ParentNode(tag, list_item_nodes)
        expected_output = ("<ol><li>list item <b>bold</b></li>"
                           "<li>list item <i>italic</i></li>"
                           "<li>item3</li></ol>"
                           )
        self.assertEqual(expected_output, parentnode.to_html())

class TestMarkdownToHTML(unittest.TestCase):
    def test_markdown_to_html(self):
        text = ("# This is some Markdown\n\n this is text\n\n"
                "1. we have different types of lists\n2. item 2\n\n"
                "## This is a **subheader**\n\n"
                "```we can code too: print('hello world')```\n\n"
                "> I might need to\n> deal with images"
                )
        html_node = markdown_to_html_node(text)
        expected_output = ("<div><h1>This is some Markdown</h1>"
                           "<p>this is text</p>"
                           "<ol><li>we have different types of lists</li>"
                           "<li>item 2</li></ol>"
                           "<h2>This is a <b>subheader</b></h2>"
                           "<pre><code>we can code too: print('hello world')</code></pre>"
                           "<blockquote> I might need to\n deal with images</blockquote>"
                           "</div>"
                           )
        self.assertEqual(expected_output, html_node.to_html())

class TestBlockToBlockType(unittest.TestCase):
    def test_heading(self):
        markdown = "## This is a headering"
        block_type = block_to_block_type(markdown)
        self.assertEqual("heading", block_type)

    def test_another_heading(self):
        markdown = "###### This is a headering"
        block_type = block_to_block_type(markdown)
        self.assertEqual("heading", block_type)

    def test_incorrect_heading(self):
        markdown = "##1"
        block_type = block_to_block_type(markdown)
        self.assertEqual("paragraph", block_type)

    def test_code(self):
        markdown = "```This is code ```"
        block_type = block_to_block_type(markdown)
        self.assertEqual("code", block_type)

    def test_empty_code(self):
        markdown = "``````"
        block_type = block_to_block_type(markdown)
        self.assertEqual("code", block_type)

    def test_quotes(self):
        markdown = "> This is a quote\n> and another \n>and another\n>"
        block_type = block_to_block_type(markdown)
        self.assertEqual("quote", block_type)

    def test_unordered_list(self):
        markdown = "* item1 \n- item2\n- item3"
        block_type = block_to_block_type(markdown)
        self.assertEqual("unordered_list", block_type)

    def test_ordered_list(self):
        markdown = "1. item1 \n2. item2\n3. item3"
        block_type = block_to_block_type(markdown)
        self.assertEqual("ordered_list", block_type)
    
    def test_paragraph(self):
        markdown = "1.item1 \n2. item2\n3. item3"
        block_type = block_to_block_type(markdown)
        self.assertEqual("paragraph", block_type)

    def test_paragraph_text(self):
        markdown = "this is just text"
        block_type = block_to_block_type(markdown)
        self.assertEqual("paragraph", block_type)

if __name__ == "__main__":
    unittest.main()
