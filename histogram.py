import re
import numpy as np
from collections import Counter
import matplotlib.pyplot as plt
import math

def extract_sizes(log_content):
    # Regular expression to match sizes (handles MiB, KiB, etc.)
    pattern = r'(\d+(?:\.\d+)?)\s*(:MiB|KiB)'
    matches = re.findall(pattern, log_content)
    
    # Convert all sizes to same unit (KiB)
    sizes = []
    for (size, unit) in matches:
        size = float(size)
        if unit == 'MiB':
            size = size * 1024
        sizes.append(size)
    
    return sizes

def analyze_distribution(sizes):
    # Create two subplots
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))
    
    # Plot 1: Histogram with logarithmic scale
    ax1.hist(sizes, bins=50, edgecolor='black')
    ax1.set_xscale('log')
    ax1.set_title('Distribution of Package Sizes (Log Scale)')
    ax1.set_xlabel('Size (KiB)')
    ax1.set_ylabel('Frequency')
    ax1.grid(True)
    
    # Plot 2: Histogram of first digits with Benford's Law comparison
    first_digits = [int(str(float(size)).replace('.', '')[0]) for size in sizes]
    digit_counts = Counter(first_digits)
    
    # Calculate percentages
    total = len(first_digits)
    actual_freq = {d: (digit_counts.get(d, 0) / total) * 100 for d in range(1, 10)}
    
    # Benford's Law expected frequencies
    benford = {
        1: 30.1, 2: 17.6, 3: 12.5, 4: 9.7, 5: 7.9,
        6: 6.7, 7: 5.8, 8: 5.1, 9: 4.6
    }
    
    # Prepare data for plotting
    digits = list(range(1, 10))
    actual_values = [actual_freq.get(d, 0) for d in digits]
    expected_values = [benford[d] for d in digits]
    
    # Create bar plot
    x = np.arange(len(digits))
    width = 0.35
    
    ax2.bar(x - width/2, actual_values, width, label='Actual', color='skyblue')
    ax2.bar(x + width/2, expected_values, width, label='Expected (Benford\'s Law)', color='lightgreen')
    
    ax2.set_xlabel('First Digit')
    ax2.set_ylabel('Frequency (%)')
    ax2.set_title('First Digit Distribution vs Benford\'s Law')
    ax2.set_xticks(x)
    ax2.set_xticklabels(digits)
    ax2.legend()
    ax2.grid(True)
    
    # Print numerical comparison
    print("\nNumerical Comparison:")
    print("Digit | Actual % | Expected %")
    print("-" * 30)
    for d in digits:
        print(f"{d:5d} | {actual_freq.get(d, 0):8.1f} | {benford[d]:8.1f}")
    
    # Adjust layout and display
    plt.tight_layout()
    plt.show()

# Read the file and analyze
with open('/home/jan/Desktop/benford/pacman.log', 'r') as file:
    content = file.read()
    sizes = extract_sizes(content)
    analyze_distribution(sizes)
