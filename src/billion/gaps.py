import sys

import billion

regularize=billion.util.regularize

class Gaps(object):
  def __init__(self, vocab, small_limit):
    self.vocab_index = billion.util.load_vocab(vocab)
    self.small_limit = small_limit

    if False:
      for w in ['the', 'computer', 'investor', 'xNONEXISTENTx', ]:
        print w, ' -> ', self.vocab_index[w]

    if False:
      self.missing_ones = ['the', 'and', 'of', 'to', 'for', 'a', 'an', 'on', 'in', 'at', 'by', 'from', ]
      if len(self.missing_ones)>30:
        print "missing_ones list too long to pack into integer"
        
      print(["%d=%s" % (self.vocab_index[w], w) for w in self.missing_ones ], "\n")

  def generate_training(self, line):
    words = regularize(line)
    vocab_indices = [ self.vocab_index[w] for w in words ]
    
    #print "Words in Line : %d" % (len(words),)
    
    for i in range(len(words)-2):
      #if vocab_indices[i] is None or vocab_indices[i+1] is None or vocab_indices[i+2] is None:
      #  continue
        
      # Middle word missing :: So ans>0
      x = [ vocab_indices[i], vocab_indices[i+2] ]
      
      # This may loop around to end/beginning of sentence...
      #x = [ vocal_indices[i-1], vocab_indices[i], vocab_indices[i+2], vocab_indices[i+3] ]
        
      if None in x: 
        continue
      
      # Pick out small words for 'easy' identification
      missing = vocab_indices[i+1]
      if missing is None: 
        #print "Missing : ", words[i+1] # Check that these are 'very infrequent'
        continue
        
      ans = (missing+2) if missing<self.small_limit else 1
      
      # So, ans==1 if this is a 'complex' word
      # small_limit+2>ans>1 if this is a 'simple' word
      # i.e. ans>0  => there is some word missing
      #      ans==0 => no word missing
      
      yield (x, ans)
    
    for i in range(len(words)-1):
      #if vocab_indices[i] is None or vocab_indices[i+1] is None:
      #  continue
        
      # Word not missing :: So ans==0
      x = [ vocab_indices[i], vocab_indices[i+1] ]
      
      # This may loop around to end/beginning of sentence...
      #x = [ vocab_indices[i-1], vocab_indices[i], vocab_indices[i+1], vocab_indices[i+2] ]
      
      if None in x: 
        continue
        
      yield (x, 0)

"""
for l, line in enumerate(inputfile):  
  if 0 == l % 10000:
    billion.util.print_thousands("Line # ", l)
    # Skip this iter - since it's going into our 'holdout set'
    continue
"""
