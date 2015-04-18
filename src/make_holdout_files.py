#! python
import billion
import sys

import argparse

parser = argparse.ArgumentParser(description='Creates Holdout files from original training file')
parser.add_argument('-i','--input',  help='Input original training file name', required=True)
parser.add_argument('-t','--train',  help='Output original training file name (no heldout lines)', required=True)

parser.add_argument(     '--heldout', help='Data file name heldout lines', required=True)

parser.add_argument(     '--truth',   help='Generated Test file name for "truth"', required=True)
parser.add_argument(     '--valid',   help='Generated Test file name for "validition"', required=True)

args = parser.parse_args()

input_file = open(args.input)

train_file = open(args.train, 'w')
heldout_file = open(args.heldout, 'w')
heldout_csv = open(args.heldout+".csv", 'w')

truth_file = open(args.truth, 'w')
validation_file = open(args.valid, 'w')

heldout_csv.write('"id","sentence"\n')
truth_file.write('"id","sentence"\n')
validation_file.write('"id","sentence"\n')

for l, line in enumerate(input_file):
  if not 0 == l % 10000:
	# Write into train_file - this is Ok for training
    train_file.write(line)
    continue
  
  heldout_file.write(line)  # Same format as training file
  #heldout_csv.write(',"+line+"\n')
    
  billion.util.print_thousands("Line # ", l)
  # Use this iter :: it's going into our 'holdout set'

  words = billion.util.stringize_test(line).split()
  heldout_csv.write('%d,"%s"\n' % (l, ' '.join(words), ))

  for i in range(1, len(words)-1):
    # i is the word we're going to drop, not first or last...
    
    word_dropped = list(words)
    del word_dropped[i]
    
    truth_file.write(   '%f,"%s"\n' % (l+(i/100.0), ' '.join(words), ))
    validation_file.write('%f,"%s"\n' % (l+(i/100.0), ' '.join(word_dropped), ))
  
input_file.close()

train_file.close()
heldout_file.close()
heldout_csv.close()

truth_file.close()
validation_file.close()
