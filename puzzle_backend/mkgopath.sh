set -eu

unset -v srcdir dstdir src dst mod ver
srcdir=$(go env GOPATH)
dstdir="${GOPATH_DST-".gopath"}"
while [ ! -f go.mod ]
do
	case "$(pwd)" in
	/)
		echo "can't find go.mod" >&2
		return 1
		;;
	esac
	cd ..
done
echo "go.mod found in $(pwd)"
go list -m -json all | jq -r '.Path + " " + .Dir' | \
while IFS=" " read -r mod src
do
	echo "===> ${mod}"
	case "${src}" in
	"") continue;;
	esac
	dst="${dstdir}/src/${mod}"
	case "${dst}" in
	*/*) mkdir -p "${dst%/*}";;
	esac
	ln -shf "${src}" "${dst}"
done
