import unittest
from parentnode import ParentNode
from leafnode import LeafNode

class TestParentNode(unittest.TestCase):
    def test_to_html_no_props(self):
        node = ParentNode(
            "p",
            [
                LeafNode("b", "Bold text"),
                LeafNode(None, "Normal text"),
                LeafNode("i", "italic text"),
                LeafNode(None, "Normal text"),
            ],
        )
        expected_output = "<p><b>Bold text</b>Normal text<i>italic text</i>Normal text</p>"
        self.assertEqual(node.to_html(), expected_output)
    
    def test_to_html_props(self):
        node = ParentNode(
            "a",
            [
                LeafNode("b", "Bold text"),
                LeafNode(None, "Normal text"),
            ],
            {"href": "https://example.dev"} 
        )
        expected_output = '<a href="https://example.dev"><b>Bold text</b>Normal text</a>'
        self.assertEqual(node.to_html(), expected_output)
    
    def test_to_html_multiple_parents(self):
        node = ParentNode(
            "a",
            [
                LeafNode("b", "Bold text"),
                LeafNode(None, "Normal text"),
            ],
            {"href": "https://example.dev"} 
        )
        
        parent_node = ParentNode("p", [node])

        expected_output = '<p><a href="https://example.dev"><b>Bold text</b>Normal text</a></p>'
        self.assertEqual(parent_node.to_html(), expected_output)
 


if __name__ == "__main__":
    unittest.main()
