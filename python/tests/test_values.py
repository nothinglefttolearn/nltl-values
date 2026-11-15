import unittest

import nltl_values


class ColorTest(unittest.TestCase):
    def test_example_color(self):
        self.assertEqual(nltl_values.color("primary"), "#DFA4E3")
        self.assertEqual(nltl_values.color("secondary"), "#879B87")
        self.assertEqual(set(nltl_values.colors()), {"primary", "secondary"})

    def test_unknown_color(self):
        with self.assertRaises(KeyError):
            nltl_values.color("nope")


class ShapeTest(unittest.TestCase):
    def test_example_shape(self):
        arrow = nltl_values.shape("arrow")
        self.assertEqual(arrow.name, "arrow")
        self.assertTrue(arrow.closed)
        self.assertEqual(len(arrow.points), 7)
        self.assertEqual(arrow.points[0], (0.0, 3.0))
        self.assertIn("arrow", nltl_values.shape_names())

    def test_unknown_shape(self):
        with self.assertRaises(KeyError):
            nltl_values.shape("nope")


if __name__ == "__main__":
    unittest.main()
