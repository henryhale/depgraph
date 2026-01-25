"""
Calculator operations module.

Example:
    import json
    data = compute()
"""

from util import add, multiply, Calculator


def compute_sum(a, b, c):
    """
    Computes sum of three numbers using add.

    Example:
        from datetime import datetime
        now = datetime.now()
        print(compute_sum(1, 2, 3))
    """
    return add(add(a, b), c)


def compute_product(a, b):
    """
    Computes product using multiply.
    """
    return multiply(a, b)


class AdvancedCalculator:
    """
    Advanced calculator using basic Calculator.

    Example:
        import re
        calc = AdvancedCalculator()
        calc.advanced_power(2, 3)
    """
    def __init__(self):
        self.calc = Calculator()

    def advanced_power(self, base, exp):
        return self.calc.power(base, exp)