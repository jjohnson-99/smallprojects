from textnode import text_node_to_html_node
from parentnode import ParentNode
from inline_markdown import text_to_textnodes

def markdown_to_blocks(markdown):
    return [section.strip() for section in markdown.split("\n\n") if section != '']

def block_to_block_type(markdown_block):
    max_num_headers = 6
    headings = ['#'*i + ' ' for i in range(1,max_num_headers+1)]
    for i in range(1, min(max_num_headers+2, len(markdown_block)+1)):
        if markdown_block[0:i] in headings:
            return "heading"

    if len(markdown_block) >= 6 and markdown_block[0:3] == '```' and markdown_block[-3:] == '```':
        return "code"
    
    markdown_block = markdown_block.strip()
    lines = markdown_block.split('\n')
    if all([line[0] == '>' for line in lines]):
        return "quote"
    if len(markdown_block) > 1 and all([line[0:2] in ['* ', '- '] for line in lines]):
        return "unordered_list"
    if len(markdown_block) > 2 and all([line[0:3] in [f'{i+1}. '] for (i, line) in enumerate(lines)]):
        return "ordered_list"
    else:
        return "paragraph"

def block_tag(text):
    block_type = block_to_block_type(text)
    if block_type == "heading":
        i = 0
        while text[i] == '#':
            i += 1
        return f"h{i}"
    if block_type == "code":
        return "pre"
    if block_type == "quote":
        return "blockquote"
    if block_type == "unordered_list":
        return "ul"
    if block_type == "ordered_list":
        return "ol"
    else:
        return "p"

def separate_list_items(text):
    text = text.strip()
    return [line.split(' ',1)[1] for line in text.split('\n')]
 
def trim_markdown_type_indicator(text):
    block_type = block_to_block_type(text)
    if block_type == "heading":
        return text.split(' ',1)[1]
    if block_type == "code":
        return text[2:-2]
    if block_type == "paragraph":
        return text
    if block_type == "quote":
        return '\n'.join([line[1:] for line in text.split('\n')]).rstrip()
    else:
        return text

def text_to_children(text):
    trimmed_text = trim_markdown_type_indicator(text)
    child_textnodes = text_to_textnodes(trimmed_text)
    child_htmlnodes = []
    for child in child_textnodes:
        child_htmlnodes.append(text_node_to_html_node(child))

    return child_htmlnodes
 
def block_to_parentnode(text):
    return ParentNode(block_tag(text), text_to_children(text))

def helper_list_items(text):
    child_htmlnodes = []
    for line in separate_list_items(text):
        child_htmlnodes.append(ParentNode("li", text_to_children(line)))

    return child_htmlnodes

def block_to_html(text):
    block_type = block_to_block_type(text)
    if block_type in ["ordered_list", "unordered_list"]:
        children = helper_list_items(text)
        return ParentNode(block_tag(text), children)
    else:
        return block_to_parentnode(text)

def markdown_to_html_node(text):
    htmlnodes = []
    blocks = markdown_to_blocks(text)

    for block in blocks:
        htmlnodes.append(block_to_html(block))
    return ParentNode("div", htmlnodes)




