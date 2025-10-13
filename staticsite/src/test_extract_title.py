import unittest
from extract_title import *

class TestExtractTitle(unittest.TestCase):
    def test_extract_title(self):
        text = "# This is a header \n ## this a subheader \n this is text"
        expected_output = "This is a header"
        self.assertEqual(expected_output, extract_title(text))
