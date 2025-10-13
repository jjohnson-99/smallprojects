import os
import pathlib
import shutil
from block_markdown import markdown_to_html_node

def copy_tree(src, dst):
    if os.path.exists(dst):
        shutil.rmtree(dst)
    os.mkdir(dst)

    src_contents = os.listdir(src)
    for content in src_contents:
        content_path = os.path.join(src, content)
        if os.path.isfile(content_path):
            shutil.copy(content_path, dst)
        else:
            new_dst_path = os.path.join(dst, content)
            copy_tree(content_path, new_dst_path)

def extract_title(markdown):
    if not markdown.startswith("# "):
        raise Exception("Markdown file does not have a header")

    return markdown.split("\n", 1)[0][2:].strip()

def generate_page(from_path, template_path, dest_path):
    print(f"Generating page from {os.path.basename(from_path)} "
          f"to {os.path.basename(dest_path)} "
          f"using {os.path.basename(template_path)}")
    with open(from_path) as f:
        markdown_text = f.read()

    with open(template_path) as f:
        template = f.read()

    html_string = markdown_to_html_node(markdown_text).to_html()
    title = extract_title(markdown_text)
    template = template.replace("{{ Title }}", title)
    template = template.replace("{{ Content }}", html_string)

    dest_dir = os.path.dirname(dest_path)
    os.makedirs(dest_dir, exist_ok = True)
    with open(dest_path, "w") as f:
        f.write(template)

def generate_pages_recursive(dir_path_content, template_path, dest_dir_path):
    markdown_files = list(pathlib.Path(dir_path_content).glob('*.md'))
    for file in markdown_files:
        dst_file = os.path.join(dest_dir_path, f"{file.stem}.html")
        generate_page(file, template_path, dst_file)

    subdirs = list(pathlib.Path(dir_path_content).glob('*/'))
    for dir in subdirs:
        dst_dir = os.path.join(dest_dir_path, dir.stem)
        generate_pages_recursive(dir, template_path, dst_dir)

def main():
    dirname = os.path.dirname(__file__)

    static = os.path.join(dirname, "../static") 
    public = os.path.join(dirname, "../public")
    copy_tree(static, public)

    content_path = os.path.join(dirname, "../content") #/index.md")
    template_path = os.path.join(dirname, "../template.html")
    index_path = os.path.join(public, "index.html")
    #generate_page(content_path, template_path, index_path)

    generate_pages_recursive(content_path, template_path, public)

if __name__ == "__main__":
    main()
