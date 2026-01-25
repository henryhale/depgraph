"""
Module 1 for testing inter-file dependencies.

Example:
    import json
    data = compute()
"""

def compute():
    """
    Computes something.

    Example:
        from datetime import datetime
        now = datetime.now()
        print(compute())
    """
    return 42

class Processor:
    """
    Processes data.

    Example:
        import re
        p = Processor()
        p.process("data")
    """
    def process(self, data):
        return data.upper()