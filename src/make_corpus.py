#! python

import billion
 
## https://docs.python.org/3/library/argparse.html
import argparse

parser = argparse.ArgumentParser(description='Parses training file in to Corpus (for GloVe)')
parser.add_argument('-i','--input', help='Input file name', required=True)
parser.add_argument('-o','--output', help='Output file name', required=True)
parser.add_argument('-l','--lines', help='# of lines', type=int, default=50e6)

args = parser.parse_args()

## show values ##
#print ("Input file: %s" % args.input )
#print ("Output file: %s" % args.output )
#print ("Lines : %d" % args.lines )

inputfile = open(args.input)
outputfile = open(args.output, 'w')

regularize=billion.util.regularize

for l, line in enumerate(inputfile):  
  if 0 == l % 10000:
    print "Line : ", l
  
  words = regularize(line)
  
  outputfile.write(' '.join(words))
  outputfile.write("\n")
  
  #print(' '.join(words))
  if l>args.lines: break
  
inputfile.close()
outputfile.close()
