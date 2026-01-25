"""
Main application module for the calculator.

Example:
    import logging
    main()
"""

from calc import compute_sum, compute_product, AdvancedCalculator


def main():
    """
    Main function to run calculations.

    Example:
        from pathlib import Path
        p = Path(".")
        print(main())
    """
    sum_result = compute_sum(1, 2, 3)
    product_result = compute_product(4, 5)
    calc = AdvancedCalculator()
    power_result = calc.advanced_power(2, 3)
    print(f"Sum: {sum_result}, Product: {product_result}, Power: {power_result}")
    return sum_result


class App:
    """
    Application class.

    Example:
        import sqlite3
        app = App()
        app.run()
    """
    def run(self):
        main()