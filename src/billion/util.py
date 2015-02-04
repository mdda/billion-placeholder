import re

import sys # for Flush
from collections import defaultdict

import os
#os.path import dirname, realpath, isfile, join
#from os import listdir

# import Levenshtein  # GPL = no thanks, today
import editdistance

# sample text string, just for demonstration to let you know how the data looks like
my_train = """
Fish , ranked 98th in the world , fired 22 aces en route to a 6-3 , 6-7 ( 5 / 7 ) , 7-6 ( 7 / 4 ) win over seventh-seeded Argentinian David Nalbandian .
Why does everything have to become such a big issue ?
AMMAN ( Reuters ) - King Abdullah of Jordan will meet U.S. President Barack Obama in Washington on April 21 to lobby on behalf of Arab states for a stronger U.S. role in Middle East peacemaking , palace officials said on Sunday .
"""

my_test = """
4,"The 's bloody body was discovered on a bed ."
5,\"Her adds that most Americans "" want to be seen in their big house with a big car . ""\"
6,"Michael Jackson could be forced to fly to the High Court in London to testify in a case being brought against him the King of Bahrain 's son ."
7,"The Wizards recovered from a 4-9 start season , and several of the team 's key players have been around long enough to know that a bad start does not necessarily lead to a bad finish ."
"""

_digits = re.compile(r'\d')
_digitsub = re.compile(r'[\d\,]+')

def regularize(line):
  words=line.split()
  
  # For consistency (most likely case), first word should be lowercased
  #  unless 'I' or a proper name (too soon to tell)
  words[0]=words[0].lower()  
  
  for i,w in enumerate(words):
    if bool(_digits.search(w)): 
      # Replace all consecutive digits with NUMBER
      words[i] = _digitsub.sub('{N}', w)
      #print "NUMBER : ", w, '', words[i]
      continue
      
  return words
  
def load_vocab(filename):
  v = defaultdict(lambda:None) # Mapping from word to number
  
  f = open(filename)
  for i,line in enumerate(f):
    w = line.split()[0] 
    v[w] = i
    
    #if i>=82905: break
    if i>82900: break
    
  f.close()
  
  return v

def comma_000(l):
  c = "    {:,}".format(l)
  return c[-(4*(len(c) // 4)):]  # Make it align in 000s too
  
def print_thousands(s_before, l, s_after="   ", overwrite=True):
  if overwrite:  # http://en.wikipedia.org/wiki/ANSI_escape_code
    print '\x1b[0G',
  print s_before, comma_000(l), s_after,   # Nice over-writing (no newline)
  sys.stdout.flush()

_test_line = re.compile(r'([\d\.]+),\"(.*)\"')  # NB no \" after the line (comments allowed)

def stringize_test(s):
  return s.replace('"', '""') # Our format for the test file should have \"\" not \"

def parse_test(line):
  m = _test_line.match(line)
  if m.groups:
    return float(m.group(1)), m.group(2).replace('""', '"')
  return 0,""


## http://hetland.org/coding/python/levenshtein.py
# Only used briefly : Prefer the C PyPi module 'python-Levenshtein'
def levenshtein_python(a,b):
    "Calculates the Levenshtein distance between a and b."
    n, m = len(a), len(b)
    if n > m:
        # Make sure n <= m, to use O(min(n,m)) space
        a,b = b,a
        n,m = m,n
        
    current = range(n+1)
    for i in range(1,m+1):
        previous, current = current, [i]+[0]*n
        for j in range(1,n+1):
            add, delete = previous[j]+1, current[j-1]+1
            change = previous[j-1]
            if a[j-1] != b[i-1]:
                change = change + 1
            current[j] = min(add, delete, change)
            
    return current[n]

## levenshtein distance measure
#levenshtein=Levenshtein.distance
levenshtein=editdistance.eval

def filename_matching(filename, offset=0):  #Offsets of 0 or 1 allowed
  if not '%' in filename:
    # This is just a regular file - no funky maths here
    if not 0==offset:
      print("Warning : Filename template must contain '%%' if you want to do offsets...")
    return filename

  if not ( offset==0 or offset==1 ):
      print("Offsets for filename incrementing can only be 0 or 1")
      exit(0)
    
  # Ok, so filename contains a number template '%d'
  # Grab all files in directory, and then do a (quicker) comparison
  base_d = os.path.dirname(os.path.realpath(filename))
  just_file = os.path.basename(filename)
  base_f = just_file[0 : just_file.index('%')]
  #print("Searching : ('%s','%s')" % (base_d, base_f))

  potentials = [ f for f in os.listdir(base_d) if os.path.isfile(os.path.join(base_d,f)) and f.startswith(base_f)]
  #print(potentials)
  
  ## Yes, this is terribly inefficient, but it allows great flexibility...
  count=1
  while (filename % count) in potentials:
    count = count+1
    
  # Now count points to an un-used file
  if 0==offset: # But (if offset==0) we want a used one (or zero)...
    count = count-1
   
  return (filename % (count+offset))
  
