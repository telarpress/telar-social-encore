#!/bin/bash

# ensure LOCAL_MODULES_BASE_PATH is set
if [[ -z "$LOCAL_MODULES_BASE_PATH" ]]; then
    echo "Error: LOCAL_MODULES_BASE_PATH environment variable is not set."
    exit 1
fi

# expand the LOCAL_MODULES_BASE_PATH
LOCAL_MODULES_BASE_PATH=$(eval echo "$LOCAL_MODULES_BASE_PATH")

echo "LOCAL_MODULES_BASE_PATH: $LOCAL_MODULES_BASE_PATH"

# define the local module paths (relative to the base path)
local_modules=(
    /telar-web
    /telar-web/micros/actions
    /telar-web/micros/admin
    /telar-web/micros/auth
    /telar-web/micros/notifications
    /telar-web/micros/profile
    /telar-web/micros/setting
    /telar-web/micros/storage
    /telar-social-go
    /telar-social-go/micros/circles
    /telar-social-go/micros/comments
    /telar-social-go/micros/gallery
    /telar-social-go/micros/posts
    /telar-social-go/micros/user-rels
    /telar-social-go/micros/vang
    /telar-social-go/micros/votes
)

# create go.work file
create_go_work() {
    echo "Creating go.work file..."
    {
        echo "go 1.18"
        echo ""
        echo "use ("
        echo "    ." # Add the current module directory
        for module in "${local_modules[@]}"; do
            echo "    ${LOCAL_MODULES_BASE_PATH}${module}"
        done
        echo ")"
    } > go.work
    echo "go.work file created with local modules."
}

# remove go.work file
remove_go_work() {
    if [[ -f go.work ]]; then
        rm go.work
        echo "Removed go.work file."
    else
        echo "go.work file does not exist."
    fi
    if [[ -f go.work.sum ]]; then
        rm go.work.sum
        echo "Removed go.work.sum file."
    else
        echo "go.work.sum file does not exist."
    fi
}

# print usage
usage() {
    echo "Usage: $0 {local|remote}"
    exit 1
}

# main script logic
if [[ $# -ne 1 ]]; then
    usage
fi

case $1 in
    local)
        create_go_work
        ;;
    remote)
        remove_go_work
        ;;
    *)
        usage
        ;;
esac
