#! python
import billion
import sys

import argparse

parser = argparse.ArgumentParser(description='Creates Holdout files from original training file')
parser.add_argument('-i','--input',  help='Input original training file name', required=True)
parser.add_argument('-t','--train',  help='Output original training file name (no heldout lines)', required=True)

parser.add_argument(     '--heldout', help='Holdout data file name for "truth"', required=True)
parser.add_argument(     '--valid',   help='Holdout data file name for "validition"', required=True)

args = parser.parse_args()

input_file = open(args.input)
train_file = open(args.train, 'w')

heldout_file = open(args.heldout, 'w')
validation_file = open(args.valid, 'w')


heldout_file.write('"id","sentence"\n')
validation_file.write('"id","sentence"\n')

for l, line in enumerate(input_file):
  if not 0 == l % 10000:
	# Write into train_file - this is Ok for training
    train_file.write(line)
    continue
    
  billion.util.print_thousands("Line # ", l)
  # Use this iter :: it's going into our 'holdout set'

  words = billion.util.stringize_test(line).split()

  for i in range(1, len(words)-1):
    # i is the word we're going to drop, not first or last...
    
    word_dropped = list(words)
    del word_dropped[i]
    
    heldout_file.write(   '%f,"%s"\n' % (l+(i/100.0), ' '.join(words), ))
    validation_file.write('%f,"%s"\n' % (l+(i/100.0), ' '.join(word_dropped), ))
  
input_file.close()
train_file.close()
heldout_file.close()
validation_file.close()
