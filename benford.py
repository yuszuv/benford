#!/usr/bin/env python

import re
from collections import Counter
import matplotlib.pyplot as plt
import math

# Expected frequencies according to Benford's Law
benford = {
    1: 30.1,
    2: 17.6,
    3: 12.5,
    4: 9.7,
    5: 7.9,
    6: 6.7,
    7: 5.8,
    8: 5.1,
    9: 4.6
}

def extract_sizes(log_content):
    # Regular expression to match sizes (handles MiB, KiB, etc.)
    pattern = r'(\d+(?:\.\d+)?)\s*(?:MiB|KiB)'
    matches = re.findall(pattern, log_content)
    
    # Convert all sizes to same unit (KiB)
    sizes = []
    for size in matches:
        size = float(size)
        if 'MiB' in log_content.split(str(size))[1][:5]:  # Check if size is in MiB
            size = size * 1024  # Convert MiB to KiB
        sizes.append(size)
    
    return sizes

def get_first_digit(number):
    # Get first digit of a number
    return int(str(float(number)).replace('.', '')[0])

def analyze_benford(sizes):
    # Get first digits
    first_digits = [get_first_digit(size) for size in sizes]
    
    # Count frequencies
    counter = Counter(first_digits)
    total = len(first_digits)
    
    # Calculate actual percentages
    actual_freq = {digit: (count/total)*100 for digit, count in counter.items()}
    
    # Plot comparison
    plt.figure(figsize=(10, 6))
    digits = list(range(1, 10))
    
    # Plot expected frequencies
    plt.plot(digits, [benford[d] for d in digits], 'b-', label='Expected (Benford\'s Law)')
    
    # Plot actual frequencies
    actual_values = [actual_freq.get(d, 0) for d in digits]
    plt.plot(digits, actual_values, 'r--', label='Actual')
    
    plt.xlabel('First Digit')
    plt.ylabel('Frequency (%)')
    plt.title('Benford\'s Law Analysis of Package Sizes')
    plt.legend()
    plt.grid(True)
    plt.show()
    
    # Print numerical comparison
    print("\nNumerical Comparison:")
    print("Digit | Actual % | Expected %")
    print("-" * 30)
    for d in digits:
        print(f"{d:5d} | {actual_freq.get(d, 0):8.1f} | {benford[d]:8.1f}")

# Read the file and analyze
with open('/home/jan/Desktop/benford/pacman.log', 'r') as file:
    content = file.read()
    sizes = extract_sizes(content)
    analyze_benford(sizes)
