"""
Math utilities module.

Example usage:
    import os
    from math import sqrt
    result = add(1, 2)
"""

import os
import sys as system

from math import sqrt, pi
from collections import defaultdict as dd


def add(a, b):
    """
    Adds two numbers.

    Example:
        from collections import defaultdict
        d = defaultdict(int)
        print(add(1, 2))  # 3
    """
    return a + b


def multiply(a, b):
    """
    Multiplies two numbers.
    """
    return a * b


class Calculator:
    """
    A simple calculator class.

    Example:
        import sys
        calc = Calculator()
        print(calc.multiply(3, 4))  # 12
    """
    def power(self, base, exp):
        return base ** exp