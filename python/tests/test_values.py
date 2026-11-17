import unittest

import nltl_values


class ColorTest(unittest.TestCase):
    def test_colors(self):
        self.assertEqual(nltl_values.color("primary"), "#DFA4E3")
        self.assertEqual(nltl_values.color("secondary"), "#879B87")
        self.assertEqual(nltl_values.color("bg"), "#141414")
        self.assertEqual(nltl_values.color("accent-glow"), "#879B8759")
        self.assertIn("border-strong", nltl_values.colors())

    def test_every_color_is_hex(self):
        for path, hex_ in nltl_values.colors().items():
            self.assertRegex(hex_, r"^#([0-9a-fA-F]{6}|[0-9a-fA-F]{8})$", path)

    def test_alias_resolves_to_target(self):
        self.assertEqual(nltl_values.color("accent"), nltl_values.color("secondary"))

    def test_unknown_color(self):
        with self.assertRaises(KeyError):
            nltl_values.color("nope")


class FontTest(unittest.TestCase):
    def test_fonts(self):
        self.assertEqual(set(nltl_values.fonts()), {"sans", "mono"})
        sans = nltl_values.font("sans")
        self.assertEqual(sans[0], "Helvetica Neue")
        self.assertEqual(sans[-1], "sans-serif")
        self.assertEqual(nltl_values.font("mono")[-1], "monospace")

    def test_unknown_font(self):
        with self.assertRaises(KeyError):
            nltl_values.font("nope")


def area(points):
    n = len(points)
    return abs(sum(points[i][0] * points[(i + 1) % n][1] - points[(i + 1) % n][0] * points[i][1] for i in range(n))) / 2


class ShapeTest(unittest.TestCase):
    def test_shapes(self):
        self.assertEqual(nltl_values.shape_names(), ["face", "space"])
        face = nltl_values.shape("face")
        self.assertTrue(face.closed)
        self.assertTrue(face.description)
        self.assertEqual(face.points, ((0.0, 0.0), (0.75, 0.25), (0.75, 1.0), (0.0, 1.0)))

    def test_face_and_space_fill_unit_square(self):
        total = area(nltl_values.shape("face").points) + area(nltl_values.shape("space").points)
        self.assertAlmostEqual(total, 1.0)

    def test_unknown_shape(self):
        with self.assertRaises(KeyError):
            nltl_values.shape("nope")


if __name__ == "__main__":
    unittest.main()
