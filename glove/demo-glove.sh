#!/bin/bash

make

DIR=demo-data/
CORPUS=${DIR}text8

if [ ! -e ${CORPUS} ]; then
  mkdir -p ${DIR}
  wget http://mattmahoney.net/dc/text8.zip --directory-prefix=${DIR}
  unzip ${CORPUS}.zip -d ${DIR}
  rm ${CORPUS}.zip
fi

VOCAB_FILE=${DIR}vocab.txt
COOCCURRENCE_FILE=${DIR}cooccurrence.bin
COOCCURRENCE_SHUF_FILE=${DIR}cooccurrence.shuf.bin
SAVE_FILE=${DIR}vectors
VERBOSE=2
MEMORY=4.0
VOCAB_MIN_COUNT=5
VECTOR_SIZE=50
MAX_ITER=15
WINDOW_SIZE=15
BINARY=2
NUM_THREADS=8
X_MAX=10

./vocab_count -min-count $VOCAB_MIN_COUNT -verbose $VERBOSE < $CORPUS > $VOCAB_FILE
if [[ $? -eq 0 ]]
  then
  ./cooccur -memory $MEMORY -vocab-file $VOCAB_FILE -verbose $VERBOSE -window-size $WINDOW_SIZE < $CORPUS > $COOCCURRENCE_FILE
  if [[ $? -eq 0 ]]
  then
    ./shuffle -memory $MEMORY -verbose $VERBOSE < $COOCCURRENCE_FILE > $COOCCURRENCE_SHUF_FILE
    if [[ $? -eq 0 ]]
    then
       ./glove -save-file $SAVE_FILE -threads $NUM_THREADS -input-file $COOCCURRENCE_SHUF_FILE -x-max $X_MAX -iter $MAX_ITER -vector-size $VECTOR_SIZE -binary $BINARY -vocab-file $VOCAB_FILE -verbose $VERBOSE
       if [[ $? -eq 0 ]]
       then
	   matlab -nodisplay -nodesktop -nojvm -nosplash < ./eval/read_and_evaluate.m 1>&2 
       fi
    fi
  fi
fi


