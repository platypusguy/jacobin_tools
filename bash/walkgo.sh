
MYNAME=`basename $0`

showusage() {
	echo
	echo 'Usage: '$MYNAME'   [ {grep options} ]   {query-string}'
	echo
	exit 1
}

run() {
    grep --include='*.go' "$1" "$2"
    exit $?
}

if [ $# -eq 1 ]; then 
	run -Hnr "$1"
fi
if [ $# -eq 2 ]; then 
    arg="$1"r
	run $arg "$2"
fi

showusage
exit 1

# Documentation
# -------------

# Position in jacobin/src.
# `walkgo.sh STRING` will show every occurence of STRING in the tree.
# `walkgo.sh -i STRING` will ignore the character case.
# For all grep's options, see `man grep`.
# https://man7.org/linux/man-pages/man1/grep.1.html
# The subdirectory name, file name, and line number are included. 
