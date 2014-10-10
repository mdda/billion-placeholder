#!/bin/bash

PREFIX=${1:-1MM}
#PREFIX=1MM
#PREFIX=ALL

X_MAX=25    # For 1MM
MAX_ITER=15 # For 1MM
MEMORY=4.0  # Standard

if [ ${PREFIX} == "ALL" ]; then
  echo "Doing ALL"
  X_MAX=100
  MAX_ITER=25
  MEMORY=24.0  # Larger machine...
fi

DIR=./data/glove/
CORPUS=${DIR}${PREFIX}_0-corpus.txt

if [ ! -e ${CORPUS} ]; then
  echo "Need to prepare corpus ${CORPUS} !"
  exit
#  mkdir -p ${DIR}
#  wget http://mattmahoney.net/dc/text8.zip --directory-prefix=${DIR}
#  unzip ${CORPUS}.zip -d ${DIR}
#  rm ${CORPUS}.zip
fi

VOCAB_FILE=${DIR}${PREFIX}_1-vocab.txt
COOCCURRENCE_FILE=${DIR}${PREFIX}_2-cooccurrence.bin
COOCCURRENCE_SHUF_FILE=${DIR}${PREFIX}_2-cooccurrence.shuf.bin
SAVE_FILE=${DIR}${PREFIX}_3-vectors
VERBOSE=2
VOCAB_MIN_COUNT=5
VECTOR_SIZE=240
#MAX_ITER=15
WINDOW_SIZE=15
BINARY=2
NUM_THREADS=8
GLOVE=./glove/

date
if [ ! -e ${VOCAB_FILE} ]; then
  ${GLOVE}/vocab_count -min-count $VOCAB_MIN_COUNT -verbose $VERBOSE < $CORPUS > $VOCAB_FILE
fi
if [[ $? -eq 0 ]]
then
  date
  if [ ! -e ${COOCCURRENCE_FILE} ]; then
    ${GLOVE}/cooccur -memory $MEMORY -vocab-file $VOCAB_FILE -verbose $VERBOSE -window-size $WINDOW_SIZE < $CORPUS > $COOCCURRENCE_FILE
  fi
  if [[ $? -eq 0 ]]
  then
    date
    if [ ! -e ${COOCCURRENCE_SHUF_FILE} ]; then
      ${GLOVE}/shuffle -memory $MEMORY -verbose $VERBOSE < $COOCCURRENCE_FILE > $COOCCURRENCE_SHUF_FILE
    fi
    if [[ $? -eq 0 ]]
    then
      date
      ${GLOVE}/glove -save-file $SAVE_FILE -threads $NUM_THREADS -input-file $COOCCURRENCE_SHUF_FILE -x-max $X_MAX -iter $MAX_ITER -vector-size $VECTOR_SIZE -binary $BINARY -vocab-file $VOCAB_FILE -verbose $VERBOSE
      if [[ $? -eq 0 ]]
      then
       echo "SUCCESS!"
        #matlab -nodisplay -nodesktop -nojvm -nosplash < ./eval/read_and_evaluate.m 1>&2 
      fi
    fi
  fi
fi

date
