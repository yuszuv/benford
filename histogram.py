import re
import numpy as np
import matplotlib.pyplot as plt

def extract_sizes(log_content):
    # Regular expression to match sizes (handles MiB, KiB, etc.)
    pattern = r'(\d+(?:\.\d+)?)\s*(MiB|KiB)'
    matches = re.findall(pattern, log_content)
    
    # Convert all sizes to same unit (KiB)
    sizes = []
    for (size, unit) in matches:
        size = float(size)
        if unit == 'MiB':
            size = size * 1024
        sizes.append(size)
    
    return sizes

def create_leading_digit_bins(max_order=3):
    """
    Create bin edges for leading digit histogram
    max_order: maximum order of magnitude (e.g., 3 means up to 1000)
    """
    bins = []
    for order in range(max_order + 1):  # For each order of magnitude (1, 10, 100, ...)
        scale = 10 ** order
        for digit in range(1, 10):  # For leading digits 1-9
            bins.append(digit * scale)
    bins.insert(0, 0)  # Add 0 as the first bin edge
    return bins

def analyze_leading_digits(sizes):
    # Determine the maximum order of magnitude in the data
    max_order = int(np.log10(max(sizes))) + 1
    
    # Create bins based on leading digits
    bins = create_leading_digit_bins(max_order)
    
    # Create the histogram
    plt.figure(figsize=(15, 6))
    plt.hist(sizes, bins=bins, edgecolor='black')
    plt.xscale('log')  # Use log scale for x-axis
    
    # Customize the plot
    plt.title('Package Size Distribution by Leading Digit')
    plt.xlabel('Size (KiB)')
    plt.ylabel('Frequency')
    plt.grid(True, which="both", ls="-", alpha=0.2)
    
    # Add vertical lines at major divisions (1, 10, 100, etc.)
    for order in range(max_order + 1):
        plt.axvline(10**order, color='r', linestyle='--', alpha=0.3)
    
    # Customize x-axis ticks to show leading digits clearly
    major_ticks = [10**i for i in range(max_order + 1)]
    plt.xticks(major_ticks, [f'10^{i}' for i in range(max_order + 1)])
    
    # Print bin counts
    counts, _ = np.histogram(sizes, bins=bins)
    print("\nBin Counts:")
    print("Range (KiB) | Count")
    print("-" * 25)
    for i in range(len(counts)):
        if counts[i] > 0:  # Only print non-empty bins
            print(f"{bins[i]:<6.1f}-{bins[i+1]:<6.1f} | {counts[i]}")
    
    plt.tight_layout()
    plt.show()

# Read the file and analyze
with open('/home/jan/Desktop/benford/pacman.log', 'r') as file:
    content = file.read()
    sizes = extract_sizes(content)
    analyze_leading_digits(sizes)
