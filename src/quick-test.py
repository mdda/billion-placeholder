import sys

import billion

print(billion.util.filename_matching('data/3-gaps/1MM_model.pickle'))
print(billion.util.filename_matching('data/3-gaps/1MM_model-2layer_%02d.pickle'))
print(billion.util.filename_matching('data/3-gaps/1MM_model-2layer_%02d.pickle', offset=1))

