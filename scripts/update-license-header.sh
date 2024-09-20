#!/bin/bash

HEADER='/*
 * Copyright (C) 2023 The GDLang Team.
 *
 * This file is part of GDLang.
 *
 * GDLang is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GDLang is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GDLang.  If not, see <http://www.gnu.org/licenses/>.
 */
 '

if [ -z "$1" ]; then
    echo "Usage: $0 <path>"
    exit 1
fi

TARGET_PATH="$1"

update_header() {
    local file="$1"
    TEMP_FILE=$(mktemp)

    # Check if the file starts with a comment block and contains "Copyright (C)"
    if head -n 20 "$file" | grep -q '^\s*/\*' && grep -q 'Copyright (C)' "$file"; then
        # Remove the existing header
        sed '/\/\*/,/\*\//d' "$file" >"$TEMP_FILE"
    else
        cat "$file" >"$TEMP_FILE"
    fi

    {
        echo -e "$HEADER"
        cat "$TEMP_FILE"
    } >"$file"

    rm -f "$TEMP_FILE"
}

for file in $(find "$TARGET_PATH" -name "*.go"); do
    update_header "$file"
done

go fmt $TARGET_PATH
