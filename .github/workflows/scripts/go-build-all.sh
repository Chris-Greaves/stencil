#!/usr/bin/env bash

# Based on https://gist.github.com/eduncan911/68775dba9d3c028181e4
# but improved to use the `go` command so it never goes out of date.

type setopt >/dev/null 2>&1

contains() {
    # Source: https://stackoverflow.com/a/8063398/7361270
    [[ $1 =~ (^|[[:space:]])$2($|[[:space:]]) ]]
}

sudo apt install zip tar

SOURCE_FILE=$(echo "$@" | sed 's/\.go//')
CURRENT_DIRECTORY="${PWD##*/}"
OUTPUT=${SOURCE_FILE:-$CURRENT_DIRECTORY} # if no src file given, use current dir name
FAILURES=""
VERSION=$(echo "$GITHUB_REF" | sed -r 's/.{10}//')

# You can set your own flags on the command line
FLAGS=${FLAGS:-"-ldflags=\"-s -w\""}

# A list of OSes to not build for, space-separated
# It can be set from the command line when the script is called.
NOT_ALLOWED_OS=${NOT_ALLOWED_OS:-"js android ios solaris illumos aix plan9 darwin"} # Skipping darwin as well because it keeps failing and I don't know why

# Get all targets
while IFS= read -r target; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    BIN_FILEPATH="bin/${OUTPUT}_${VERSION}_${GOOS}-${GOARCH}"
    BIN_FILENAME="${BIN_FILEPATH}/${OUTPUT}"
    
    if contains "$NOT_ALLOWED_OS" "$GOOS" ; then
        continue
    fi
    
    # Check for arm and set arm version
    if [[ $GOARCH == "arm" ]]; then
        # Set what arm versions each platform supports
        if [[ $GOOS == "darwin" ]]; then
            arms="7"
        elif [[ $GOOS == "windows" ]]; then
             # This is a guess, it's not clear what Windows supports from the docs
             # But I was able to build all these on my machine
            arms="5 6 7" 
        elif [[ $GOOS == *"bsd"  ]]; then
            arms="6 7"
        else
            # Linux goes here
            arms="5 6 7"
        fi

        # Now do the arm build
        for GOARM in $arms; do
            BIN_FILEPATH="bin/${OUTPUT}_${VERSION}_${GOOS}-${GOARCH}${GOARM}"
            BIN_FILENAME="${BIN_FILEPATH}/${OUTPUT}"
            if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
            CMD="GOARM=${GOARM} GOOS=${GOOS} GOARCH=${GOARCH} go get && go build $FLAGS -o ${BIN_FILENAME} $@"
            echo "${CMD}"
            eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}" 
            if [[ "${GOOS}" == "windows" ]]; then
                zip -r ${BIN_FILEPATH}.zip ${BIN_FILEPATH}
            else
                tar -czvf ${BIN_FILEPATH}.tar.gz ${BIN_FILEPATH}
            fi
        done
    else
        # Build non-arm here
        if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
        CMD="GOOS=${GOOS} GOARCH=${GOARCH} go get && go build $FLAGS -o ${BIN_FILENAME} $@"
        echo "${CMD}"
        eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}"
        if [[ "${GOOS}" == "windows" ]]; then
            zip -r ${BIN_FILEPATH}.zip ${BIN_FILEPATH}
        else
            tar -czvf ${BIN_FILEPATH}.tar.gz ${BIN_FILEPATH}
        fi
    fi
done <<< "$(go tool dist list)"

if [[ "${FAILURES}" != "" ]]; then
    echo ""
    echo "${SCRIPT_NAME} failed on: ${FAILURES}"
    exit 1
fi
