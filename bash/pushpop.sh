nl -ba | grep -E 'PUSH|POP' | sort -n
exit $?

# Documentation
# -------------

# Reads from stdin
# Extracts only lines containing PUSH or POP
# Prefixes each with its line number
# Sorts by line number (helpful if input might be unordered)

# nl -ba
# Numbers all lines from stdin (-ba = number all lines, even blank ones).

# grep -E 'PUSH|POP'
# Keeps only the lines with PUSH or POP.

# sort -n
# Sorts numerically by the line number at the start of each line.
