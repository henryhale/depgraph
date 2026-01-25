"""
This is a test module for dependency graph analysis.

Example usage:
    import os
    from math import sqrt
    result = add(1, 2)
"""

def add(a, b):
    """
    Adds two numbers.

    Example:
        from collections import defaultdict
        d = defaultdict(int)
        print(add(1, 2))  # 3
    """
    return a + b

class Calculator:
    """
    A simple calculator class.

    Example:
        import sys
        calc = Calculator()
        print(calc.multiply(3, 4))  # 12
    """
    def multiply(self, x, y):
        return x * y