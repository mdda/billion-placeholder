#! python
import billion
import sys

import argparse

parser = argparse.ArgumentParser(description='Creates Holdout files from Corpus file')
parser.add_argument('-i','--input', help='Input Corpus file name', required=True)
parser.add_argument('-o','--orig', help='Holdout data file name for "truth"', required=True)
parser.add_argument('-t','--test', help='Holdout data file name for "test"', required=True)

args = parser.parse_args()

inputfile = open(args.input)
origfile = open(args.orig, 'w')
testfile = open(args.test, 'w')

for l, line in enumerate(inputfile):  
  if not 0 == l % 10000:
    continue
    
  billion.util.print_thousands("Line # ", l)
  # Use this iter :: it's going into our 'holdout set'

  words = line.split()
    
  for i in range(1, len(words)-1):
    # i is the word we're going to drop, not first or last...
    
    word_dropped = list(words)
    del word_dropped[i]
    
    origfile.write('%f,"%s"\n' % (l+(i/100.0), ' '.join(words), ))
    testfile.write('%f,"%s"\n' % (l+(i/100.0), ' '.join(word_dropped), ))
  
inputfile.close()
origfile.close()
testfile.close()
