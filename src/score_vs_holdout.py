#! python
import billion
import sys
import math

import argparse

parser = argparse.ArgumentParser(description='Tests files vs ideal - for Holdout data, for instance')
parser.add_argument('-o','--orig', help='Data file name for "truth"', required=True)
parser.add_argument('-s','--submission', help='Data file name for submission', required=True)

args = parser.parse_args()

origfile = open(args.orig)
subsfile = open(args.submission)

header_orig = origfile.readline()
header_subs = subsfile.readline()

def pretty_print(l, total, total2, total_cnt):
  mean = float(total)/total_cnt
  sd   = float(total2)/total_cnt - mean*mean
  conf = math.sqrt(sd/total_cnt)
  billion.util.print_thousands("Lines = ", l, " -> %7.5f  1SD=(%7.5f, %7.5f)" % (mean, mean-conf, mean+conf))

total, total2, total_cnt = 0, 0, 0
for l, line in enumerate(subsfile):  
  line_orig = origfile.readline()
  
  i_subs, text = billion.util.parse_test(line)
  i_orig, text_orig = billion.util.parse_test(line_orig)
  
  if i_subs != i_orig:
    print "Indices don't match up - aborting"
    exit()
  
  ## Now compare text to text_orig, and total up differences
  dist = billion.util.levenshtein(text, text_orig)  # Shorter one first...
  total  += dist
  total2 += dist*dist
  total_cnt+=1
  
  if 0 == l % 1000:
    pretty_print(l, total, total2, total_cnt)

pretty_print(l, total, total2, total_cnt)

origfile.close()
subsfile.close()

"""
Benchmark 
PurePython
71620 5.54362547298
71621 5.54360392058
71622 5.54358236879
71623 5.5435608176
71624 5.54352530541
71625 5.54353167844
71626 5.54349616765
71627 5.54351650193

EditDistance
71620 5.54362547298
71621 5.54360392058
71622 5.54358236879
71623 5.5435608176
71624 5.54352530541
71625 5.54353167844
71626 5.54349616765
71627 5.54351650193

"""
